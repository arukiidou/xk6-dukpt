package example

import (
	"encoding/base32"
	"errors"
	"fmt"
	"reflect"

	"github.com/grafana/sobek"
)

func (m *module) b32encode(input sobek.Value, encoding string) (any, error) {
	data, err := stringOrArrayBuffer(input, m.vu.Runtime())
	if err != nil {
		return nil, err
	}

	enc, err := encodingFor(encoding)
	if err != nil {
		return nil, err
	}

	return m.vu.Runtime().ToValue(enc.EncodeToString(data)), nil
}

func (m *module) b32decode(input string, encoding string, format string) (any, error) {
	rt := m.vu.Runtime()

	enc, err := encodingFor(encoding)
	if err != nil {
		return nil, err
	}

	data, err := enc.DecodeString(input)
	if err != nil {
		return nil, err
	}

	if format == stringFormat {
		return rt.ToValue(string(data)).ToObject(rt), nil
	}

	return rt.NewArrayBuffer(data), nil
}

var (
	errInvalidEncoding = errors.New("invalid encoding")
	errInvalidType     = errors.New("invalid type")
)

const (
	stdEncoding    = "std"
	stdrawEncoding = "stdraw"
	hexEncoding    = "hex"
	hexrawEncoding = "hexraw"

	stringFormat = "s"
)

func stringOrArrayBuffer(input sobek.Value, runtime *sobek.Runtime) ([]byte, error) {
	var data []byte

	switch input.ExportType() {
	case reflect.TypeFor[string]():
		var str string

		if err := runtime.ExportTo(input, &str); err != nil {
			return nil, err
		}

		data = []byte(str)

	case reflect.TypeFor[[]byte]():
		if err := runtime.ExportTo(input, &data); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("%w: String or ArrayBuffer expected", errInvalidType)
	}

	return data, nil
}

func encodingFor(encoding string) (*base32.Encoding, error) {
	switch encoding {
	case stdEncoding, "":
		return base32.StdEncoding, nil
	case stdrawEncoding:
		return base32.StdEncoding.WithPadding(base32.NoPadding), nil
	case hexEncoding:
		return base32.HexEncoding, nil
	case hexrawEncoding:
		return base32.HexEncoding.WithPadding(base32.NoPadding), nil
	default:
		return nil, fmt.Errorf("%w: %s", errInvalidEncoding, encoding)
	}
}
