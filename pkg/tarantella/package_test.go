package tarantella

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"path"
	"testing"

	"github.com/ryboe/q"
	"github.com/stretchr/testify/require"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func TestParse(t *testing.T) {
	parsePackageData := func(t *testing.T, name string) {
		packData, err := os.Open("../testdata/" + name + ".hex")
		require.NoError(t, err)
		defer packData.Close()

		hexdec := hex.NewDecoder(packData)
		bb, err := io.ReadAll(hexdec)
		require.NoError(t, err)

		msgdec := msgpack.NewDecoder(bytes.NewBuffer(bb))

		var length1 uint32
		err = msgdec.Decode(&length1)
		require.NoError(t, err)

		q.Q(length1)

		mh, err := msgdec.DecodeMap()
		require.NoError(t, err)

		q.Q("header", mh, Cast(mh.(map[any]any)))

		mb, err := msgdec.DecodeMap()
		require.NoError(t, err)

		q.Q("body", mb, Cast(mb.(map[any]any)))

		os.WriteFile("../testdata/"+name+".header.yaml", MustYaml(mh.(map[any]any), Cast), 0o644)
		os.WriteFile("../testdata/"+name+".body.yaml", MustYaml(mb.(map[any]any), Cast), 0o644)
	}

	t.Run("select.281", func(t *testing.T) {
		parsePackageData(t, path.Base(t.Name()))
	})
	t.Run("select.289", func(t *testing.T) {
		parsePackageData(t, path.Base(t.Name()))
	})
	t.Run("insert.514", func(t *testing.T) {
		parsePackageData(t, path.Base(t.Name()))
	})
	t.Run("ping", func(t *testing.T) {
		parsePackageData(t, path.Base(t.Name()))
	})
	t.Run("execute.1", func(t *testing.T) {
		parsePackageData(t, path.Base(t.Name()))
	})
}
