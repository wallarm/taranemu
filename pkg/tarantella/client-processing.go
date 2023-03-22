package tarantella

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
	"gopkg.in/yaml.v3"
)

const (
	versionMajor = 2
	versionMinor = 10
	versionPatch = 4

	IPROTO_SALT_SIZE     = 32  //nolint
	IPROTO_GREETING_SIZE = 128 //nolint
)

type (
	clientConnection struct {
		ctx      context.Context
		c        net.Conn
		username string // from IPROTO_AUTH
		baseDir  string
	}
)

var errUnanswerable = errors.New("unanswerable")

func processClient(ctx context.Context, conn net.Conn, baseDir string) error {
	// uncomment this block, if you want to stop propositioning panic
	// 	defer func() {
	// 	if r := recover(); r != nil {
	// 		if e, ok := r.(error); ok {
	// 			st := pkgerrors.MarshalStack(e)
	// 			log.Error().Err(e).Any("stack", st).Msg("Recovered")
	// 		} else {
	// 			log.Error().Msgf("Recovered %#v", r)
	// 		}
	// 	}
	// }()

	clc := &clientConnection{
		ctx:     ctx,
		c:       conn,
		baseDir: baseDir,
	}
	return clc.loop()
}

func (clc *clientConnection) loop() error {
	log.Info().Any("local", clc.c.LocalAddr()).
		Any("remote", clc.c.RemoteAddr()).
		Str("base-dir", clc.baseDir).
		Msg("Processing connection")

	go func() {
		<-clc.ctx.Done()
		log.Info().Msg("Closing client socket")
		clc.c.Close() //nolint: errcheck
	}()

	// dc := &tio.DeadlineIO{To: timeout, Conn: clc.c}
	// r := bufio.NewReaderSize(dc, 128*1024)
	// w := bufio.NewWriterSize(dc, 128*1024)

	_, err := clc.c.Write(createGreeting())
	if err != nil {
		return errors.Wrap(err, "unable to send greeting")
	}

	for {
		req, err := clc.readRequest(clc.c)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse incoming request")
			return errors.Wrap(err, "failed to parse incoming request")
		}
		res, err := clc.prepareResponse(req)
		if errors.Is(err, errUnanswerable) {
			continue
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to prepare response")
			return errors.Wrap(err, "failed to prepare response")
		}

		err = clc.writeResponse(res, clc.c)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send response")
		}

	}
}

// sinkDir returns the directory for data files in form baseDir + <username>
func (clc *clientConnection) sinkDir() string {
	d := filepath.Join(clc.baseDir, clc.username)
	os.MkdirAll(d, 0o755) //nolint: errcheck
	return d
}

// sinkDir returns the directory for data files in form baseDir + <username>
func (clc *clientConnection) spaceFile(spaceID uint64) string {
	return filepath.Join(clc.sinkDir(), fmt.Sprintf("%d.yaml", spaceID))
}

// prepareResponse can returns nil, errUnanswerable if request
// doesn't require a response like IPROTO_WATCH request
func (clc *clientConnection) prepareResponse(req *Package) (*Package, error) {
	res := &Package{}

	requestType := req.HeaderRequestType()

	requestTypeDescription := RequestTypeDescr(requestType)

	log.Debug().Str("request-type", requestTypeDescription).Msg("Incoming request")

	// each response has to have this field
	res.SetHeader(IPROTO_SYNC, req.HeaderSync())
	// indeed, as a stub we will be pretend to be good boy
	res.SetHeader(IPROTO_REQUEST_TYPE, IPROTO_OK)

	switch requestType {
	case IPROTO_ID:
		res.SetBody(IPROTO_VERSION, 4)
		res.SetBody(IPROTO_FEATURES, req.BodyFeatures())
	case IPROTO_AUTH:
		clc.username = req.BodyUsername()
		if clc.username == "" {
			clc.username = "_incognito_"
		}
	case IPROTO_PING:
		res.SetHeader(IPROTO_SCHEMA_VERSION, schemaVersion)
	case IPROTO_EXECUTE:
		res.SetHeader(IPROTO_SCHEMA_VERSION, schemaVersion)
		return clc.processExecute(req, res)
	case IPROTO_WATCH:
		return nil, errUnanswerable
	case IPROTO_SELECT:
		return clc.processSelect(req, res)
	case IPROTO_INSERT:
		return clc.processInsert(req, res)
	default:
		log.Warn().Str("request-type", requestTypeDescription).Msg("Unimplemented or unknown request type")
		res.SetHeader(IPROTO_REQUEST_TYPE, IPROTO_TYPE_ERROR|tarantool.ErrUnknownRequestType)
	}
	return res, nil
}

func (clc *clientConnection) processExecute(req, res *Package) (*Package, error) {
	sqlText := req.BodySQLText()
	log.Info().Str("sql-text", sqlText).Msg("SQL execute")

	res.SetBody(IPROTO_DATA, dummyExecute1["data"])
	res.SetBody(IPROTO_METADATA, dummyExecute1["metadata"])

	return res, nil
}

func (clc *clientConnection) processInsert(req, res *Package) (*Package, error) {
	res.SetHeader(IPROTO_SCHEMA_VERSION, schemaVersion)

	spaceFile := clc.spaceFile(req.BodySpaceID())

	log.Debug().Str("tgt-file", spaceFile).Msg("Writing INSERT event")

	// append to the space file
	f, err := os.OpenFile(spaceFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error().Err(err).
			Str("user", clc.username).Str("file", spaceFile).
			Msg("Unable to open file for write")
		res.SetBody(IPROTO_ERROR_24, fmt.Sprintf("TARANTELLA: unable to open file %s while trying to save insert data", spaceFile))
		return res, nil
	}
	defer f.Close() //nolint: errcheck
	if _, e := f.WriteString("---\n"); e != nil {
		log.Error().Err(e).Msg("Unable to save request into file")
	}
	enc := yaml.NewEncoder(f)
	if e := enc.Encode(req.Info()); e != nil {
		log.Error().Err(e).Msg("Unable to save request into file")
	}
	enc.Close() //nolint: errcheck

	return res, nil
}

func (clc *clientConnection) processSelect(req, res *Package) (*Package, error) {
	res.SetHeader(IPROTO_SCHEMA_VERSION, schemaVersion)

	spaceID := req.BodySpaceID()

	log.Debug().Uint64("space-id", spaceID).Msg("IPROTO_SELECT(0x1) on space")

	spaceFile := clc.spaceFile(req.BodySpaceID())
	if f, e := os.Open(spaceFile); e == nil {
		defer f.Close() //nolint: errcheck
		dec := yaml.NewDecoder(f)

		body := []any{}
		for {
			ri := new(RequestInfo)
			err := dec.Decode(ri)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Warn().Err(err).Str("space-file", spaceFile).Msg("Unable to decode spaceFile")
				res.SetBody(IPROTO_ERROR_24, fmt.Sprintf("TARANTELLA: unable to decode file for space %d", spaceID))
				return res, nil
			}
			body = append(body, ri.B["⚡IPROTO_TUPLE(0x21)"])
		}
		res.SetBody(IPROTO_DATA, body)
		return res, nil
	}

	switch spaceID {
	case BOX_VSPACE_ID:
		// here we ignore all other passed flags!
		res.SetBody(IPROTO_DATA, dummySpaces)
	case BOX_VINDEX_ID:
		// here we ignore all other passed flags!
		res.SetBody(IPROTO_DATA, dummyIndexes)

	default:
		log.Warn().Uint64("space-id", spaceID).Msg("IPROTO_SELECT(0x1) on space unsupported")
		res.SetBody(IPROTO_ERROR_24, fmt.Sprintf("TARANTELLA: space with id %d is not exist or not served", spaceID))
	}
	return res, nil
}

func (clc *clientConnection) writeResponse(res *Package, w io.Writer) error {
	if e := res.Encode(); e != nil {
		return errors.Wrap(e, "unable to encode response")
	}
	bb := res.ToBytes()

	if _, e := w.Write(bb); e != nil {
		return errors.Wrap(e, "unable to write response packet")
	}

	return nil
}

func (clc *clientConnection) readRequest(r io.Reader) (*Package, error) {
	req := &Package{}

	var rawLen [5]byte

	if _, e := io.ReadFull(r, rawLen[:]); e != nil {
		return nil, e
	}

	length, err := req.SetLen(rawLen)
	if err != nil {
		return nil, err
	}

	rawData := make([]byte, length)
	_, err = io.ReadFull(r, rawData)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read request package")
	}

	err = req.Decode(rawData)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode request package")
	}

	return req, nil
}
