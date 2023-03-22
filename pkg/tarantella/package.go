package tarantella

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/yaml.v3"
)

// see https://www.tarantool.io/en/doc/latest/dev_guide/internals/iproto/keys/

type (
	// Package represents a request or response PROTO package
	Package struct {
		rawLen  []byte // 5 bytes length
		rawData []byte // header + body
		header  map[any]any
		body    map[any]any
	}

	// RequestInfo describes package while marshall or unmarshal package into or from a file
	RequestInfo struct {
		RT string      `yaml:"rt,omitempty"` // request type
		H  map[any]any `yaml:"h,omitempty,flow"`
		B  map[any]any `yaml:"b,omitempty,flow"`
	}
)

// Yaml mutates map with mutators (if any), an return yaml representation for it
func Yaml(src map[any]any, mutators ...func(src map[any]any) map[any]any) ([]byte, error) {
	for _, m := range mutators {
		src = m(src)
	}
	return yaml.Marshal(src)
}

// MustYaml is like Yaml but doesn't make mistakes
func MustYaml(src map[any]any, mutators ...func(src map[any]any) map[any]any) []byte {
	bb, err := Yaml(src, mutators...)
	if err != nil {
		panic(err)
	}
	return bb
}

// Cast are trying to cast keys to IPROTO constants and some keys values too
func Cast(src map[any]any) map[any]any {
	dst := map[any]any{}

	for k, v := range src {
		var (
			castk = k
			castv = v
		)
		if kname, ok := iproto_key[k]; ok {
			castk = fmt.Sprintf("⚡%s(%#v)", kname, k)
		}

		switch k {
		case IPROTO_REQUEST_TYPE:
			if vname, ok := iproto_type[v]; ok {
				castv = fmt.Sprintf("⚡%s(%#v)", vname, v)
			}
		}

		dst[castk] = castv
	}
	return dst
}

// Info concatenate header + body and returns it like map with 2 values
func (pack *Package) Info() any {
	return &RequestInfo{
		RT: RequestTypeDescr(pack.HeaderRequestType()),
		H:  Cast(pack.header),
		B:  Cast(pack.body),
	}
}

// RequestTypeDescr dress request type into beautiful cloths
func RequestTypeDescr(requestType uint64) string {
	if v, ok := iproto_type[requestType]; ok {
		return fmt.Sprintf("%s(%#v)", v, requestType)
	}
	return fmt.Sprintf("%#v", requestType)
}

// CastInfo applies returns info of the package
func (pack *Package) CastInfo() any {
	info := &struct {
		Length int
		Header map[any]any
		Body   map[any]any
	}{}
	info.Length = len(pack.rawData)
	info.Header = Cast(pack.header)
	info.Body = Cast(pack.body)

	return info
}

// HeaderRequestType returns IPROTO_REQUEST_TYPE
func (pack *Package) HeaderRequestType() uint64 {
	return header[uint64](pack, IPROTO_REQUEST_TYPE)
}

// HeaderSync returns IPROTO_SYNC
func (pack *Package) HeaderSync() uint64 {
	return header[uint64](pack, IPROTO_SYNC)
}

// BodyVersion returns IPROTO_VERSION
func (pack *Package) BodyVersion() uint64 {
	return body[uint64](pack, IPROTO_VERSION)
}

// BodySpaceID returns IPROTO_SPACE_ID
func (pack *Package) BodySpaceID() uint64 {
	return body[uint64](pack, IPROTO_SPACE_ID)
}

// BodySQLText returns IPROTO_SQL_TEXT
func (pack *Package) BodySQLText() string {
	return body[string](pack, IPROTO_SQL_TEXT)
}

// BodyFeatures returns IPROTO_FEATURES
func (pack *Package) BodyFeatures() []any {
	return body[[]any](pack, IPROTO_FEATURES)
}

// BodyUsername returns IPROTO_USER_NAME
func (pack *Package) BodyUsername() string {
	return body[string](pack, IPROTO_USER_NAME)
}

// header returns a key value from the header
func header[T any](pack *Package, key uint64) T {
	if v, ok := pack.header[key].(T); !ok {
		panic(errors.Errorf("unable to cast header value '%#v' with key '%v' for type '%T'", pack.header[key], key, *new(T)))
	} else {
		return v
	}
}

// body returns a key value from the body
func body[T any](pack *Package, key uint64) T {
	if v, ok := pack.body[key].(T); !ok {
		panic(errors.Errorf("unable to cast body value '%#v' with key '%v' for type '%T'", pack.body[key], key, *new(T)))
	} else {
		return v
	}
}

// SetHeader sets one field of the package header
func (pack *Package) SetHeader(k, v any) {
	if pack.header == nil {
		pack.header = make(map[any]any)
	}
	pack.header[k] = v
}

// SetBody sets one field of the package body
func (pack *Package) SetBody(k, v any) {
	if pack.body == nil {
		pack.body = make(map[any]any)
	}
	pack.body[k] = v
}

// Encode encodes Header, Body and Len into byte arrays
func (pack *Package) Encode() error {
	if pack.header == nil {
		pack.header = make(map[any]any)
	}
	if pack.body == nil {
		pack.body = make(map[any]any)
	}

	buff := &bytes.Buffer{}
	e := msgpack.NewEncoder(buff)

	if e := e.Encode(pack.header); e != nil {
		return errors.Wrap(e, "unable to encode header")
	}
	if e := e.Encode(pack.body); e != nil {
		return errors.Wrap(e, "unable to encode body")
	}
	pack.rawData = buff.Bytes()
	pack.rawLen = pack.len()
	return nil
}

// Decode gets byte stream and read from it header and body of the package as a map[any]any
func (pack *Package) Decode(rawData []byte) error {
	pack.rawData = make([]byte, len(rawData))
	copy(pack.rawData, rawData)

	d := msgpack.NewDecoder(bytes.NewBuffer(pack.rawData))

	var (
		err error
		ok  bool
	)

	m, err := d.DecodeMap()
	if err != nil {
		return errors.Wrap(err, "unable to decode package header")
	}

	if pack.header, ok = m.(map[any]any); !ok {
		return errors.New("unable to decode package header")
	}

	m, err = d.DecodeMap()
	if err != nil {
		return errors.Wrap(err, "unable to decode package body")
	}

	if pack.body, ok = m.(map[any]any); !ok {
		return errors.New("unable to decode package body")
	}

	return nil
}

// len returns length rawData package size preamble
func (pack *Package) len() []byte {
	l := [5]byte{}
	l[0] = 0xce
	n := len(pack.rawData)
	l[1] = byte(n >> 24)
	l[2] = byte(n >> 16)
	l[3] = byte(n >> 8)
	l[4] = byte(n)

	return l[:]
}

// SetLen sets len of the package
func (pack *Package) SetLen(rawLen [5]byte) (uint32, error) {
	var length uint32
	err := msgpack.NewDecoder(bytes.NewBuffer(rawLen[:])).Decode(&length)
	if err != nil {
		return 0, errors.Wrap(err, "unable to decode package length")
	}

	if length == 0 {
		return 0, errors.New("Response should not be 0 length")
	}
	pack.rawLen = rawLen[:]
	return length, nil
}

// ToBytes emit Len+Head+Body as a whole
func (pack *Package) ToBytes() []byte {
	var buff []byte
	buff = append(buff, pack.rawLen...)
	buff = append(buff, pack.rawData...)
	return buff
}

func createGreeting() []byte {
	greetingBuf := &bytes.Buffer{}

	h := IPROTO_GREETING_SIZE / 2

	fmt.Fprintf(greetingBuf, "Tarantool %d.%d.%d (Binary) ", versionMajor, versionMinor, versionPatch)
	greetingBuf.WriteString(dummyInstanceID)

	r := len(greetingBuf.Bytes())
	greetingBuf.WriteString(strings.Repeat(" ", h-r-1))
	greetingBuf.WriteString("\n")

	salt := [IPROTO_SALT_SIZE]byte{}
	rand.Reader.Read(salt[:]) //nolint: errcheck
	greetingBuf.WriteString(base64.StdEncoding.EncodeToString(salt[:]))

	rest := IPROTO_GREETING_SIZE - len(greetingBuf.Bytes())
	// salt
	greetingBuf.WriteString(strings.Repeat(" ", rest-1))
	greetingBuf.WriteString("\n")

	if len(greetingBuf.Bytes()) != IPROTO_GREETING_SIZE {
		panic("Illegal size of greeting. Fix the implementation.")
	}

	log.Debug().Str("greeting", greetingBuf.String()).Msg("Greeting prepared")

	return greetingBuf.Bytes()
}
