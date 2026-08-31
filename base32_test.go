package example

import (
	"testing"

	"github.com/grafana/sobek"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

const (
	message = "The quick brown fox jumps over the lazy dog."
	std     = "KRUGKIDROVUWG2ZAMJZG653OEBTG66BANJ2W24DTEBXXMZLSEB2GQZJANRQXU6JAMRXWOLQ="
	stdraw  = "KRUGKIDROVUWG2ZAMJZG653OEBTG66BANJ2W24DTEBXXMZLSEB2GQZJANRQXU6JAMRXWOLQ"
	hex     = "AHK6A83HELKM6QP0C9P6UTRE41J6UU10D9QMQS3J41NNCPBI41Q6GP90DHGNKU90CHNMEBG="
	hexraw  = "AHK6A83HELKM6QP0C9P6UTRE41J6UU10D9QMQS3J41NNCPBI41Q6GP90DHGNKU90CHNMEBG"
)

func Test_module_b32encode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		encoding string
		want     string
		wantErr  bool
	}{
		{name: "std", input: message, encoding: "std", want: std},
		{name: "stdraw", input: message, encoding: "stdraw", want: stdraw},
		{name: "hex", input: message, encoding: "hex", want: hex},
		{name: "hexraw", input: message, encoding: "hexraw", want: hexraw},
		{name: "invalid encoding", input: message, encoding: "invalid", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod := &module{vu: modulestest.NewRuntime(t).VU}
			toValue := mod.vu.Runtime().ToValue

			got, err := mod.b32encode(toValue(tt.input), tt.encoding)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			value, ok := got.(sobek.Value)
			require.True(t, ok)

			require.Equal(t, tt.want, value.String())
		})
	}
}

func Test_module_b32encode_invalid_type(t *testing.T) {
	t.Parallel()

	mod := &module{vu: modulestest.NewRuntime(t).VU}

	_, err := mod.b32encode(mod.vu.Runtime().NewObject(), "std")

	require.ErrorIs(t, err, errInvalidType)
}

func Test_module_b32decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		encoding string
		want     string
		wantErr  bool
	}{
		{name: "std", input: std, encoding: "std", want: message},
		{name: "stdraw", input: stdraw, encoding: "stdraw", want: message},
		{name: "hex", input: hex, encoding: "hex", want: message},
		{name: "hexraw", input: hexraw, encoding: "hexraw", want: message},
		{name: "invalid encoding", input: message, encoding: "invalid", wantErr: true},
		{name: "invalid input", input: "invalid!", encoding: "std", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod := &module{vu: modulestest.NewRuntime(t).VU}

			got, err := mod.b32decode(tt.input, tt.encoding, "s")

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			value, ok := got.(sobek.Value)
			require.True(t, ok)
			require.NotNil(t, value)

			require.Equal(t, tt.want, value.String())

			want := value.String()

			got, err = mod.b32decode(tt.input, tt.encoding, "")
			require.NoError(t, err)

			buff, ok := got.(sobek.ArrayBuffer)
			require.True(t, ok)

			require.Equal(t, want, string(buff.Bytes()))
		})
	}
}

func Test_encodingFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		encoding string
		want     string
		wantErr  bool
	}{
		{name: "std", encoding: "std", want: std},
		{name: "stdraw", encoding: "stdraw", want: stdraw},
		{name: "hex", encoding: "hex", want: hex},
		{name: "hexraw", encoding: "hexraw", want: hexraw},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := encodingFor(tt.encoding)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got.EncodeToString([]byte(message)))
		})
	}
}

func Test_stringOrArrayBuffer(t *testing.T) {
	t.Parallel()

	rt := sobek.New()

	tests := []struct {
		name    string
		input   sobek.Value
		want    string
		wantErr bool
	}{
		{name: "string", input: rt.ToValue(message), want: message},
		{name: "ArrayBuffer", input: rt.ToValue([]byte(message)), want: message},
		{name: "invalid type", input: rt.NewObject(), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := stringOrArrayBuffer(tt.input, sobek.New())

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, []byte(tt.want), got)
		})
	}
}
