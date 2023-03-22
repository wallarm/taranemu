package tarantella

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// StartServer starts the tarantool emulator
func StartServer(ctx context.Context, listenOn, dataDir string) error {
	log.Debug().Msgf("Launching server on %s...", listenOn)

	lc := &net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", listenOn)
	if err != nil {
		return errors.Wrapf(err, "unable to listen on %s", listenOn)
	}

	log.Info().Msgf("Server started on %s...", ln.Addr().String())

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing listen socket...")
		ln.Close() //nolint: errcheck
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return errors.Wrapf(err, "unable to accept on %s", ln.Addr().String())
		}
		go processClient(ctx, conn, dataDir) //nolint: errcheck
	}
}
