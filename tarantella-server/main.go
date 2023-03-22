package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/wallarm/tarantella/pkg/tarantella"
)

var (
	cfgLevel   = os.Getenv("LEVEL")
	cfgListen  = os.Getenv("LISTEN")
	cfgDataDir = os.Getenv("DATA_DIR")
)

func main() {
	if cfgLevel == "" {
		cfgLevel = "info"
	}
	if cfgDataDir == "" {
		cfgDataDir = os.TempDir() + "/tarantella"
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.With().Stack().Logger()

	if level, e := zerolog.ParseLevel(cfgLevel); e == nil {
		zerolog.SetGlobalLevel(level)
	} else {
		fmt.Fprintf(os.Stderr, "unable to parse level %s", cfgLevel)
	}

	doMain(func(ctx context.Context, cancel context.CancelFunc) error {
		defer cancel()
		return tarantella.StartServer(ctx, cfgListen, cfgDataDir)
	})
}

// doMain starts function runFunc with specified context. The context will be canceled
// by SIGTERM or SIGINT signal (Ctrl+C for example)
// beforeExit function must be executed immediately before exit
func doMain(runFunc func(ctx context.Context, cancel context.CancelFunc) error) {
	// context should be canceled while Int signal will be caught
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// main processing loop
	retChan := make(chan error, 1)
	go func() {
		err2 := runFunc(ctx, cancel)
		if err2 != nil {
			retChan <- err2
		}
		close(retChan)
	}()

	// Waiting signals from OS
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		log.Warn().Msgf("Signal '%s' was caught. Exiting", <-quit)
		cancel()
	}()

	// Listening for the main loop response
	for e := range retChan {
		log.Info().Err(e).Msg("Exiting.")
	}
}
