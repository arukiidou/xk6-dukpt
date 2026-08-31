package example

import (
	"testing"

	"github.com/grafana/sobek"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

func Test_newRandom(t *testing.T) {
	t.Parallel()

	r := newRandom(42)

	require.Equal(t, int64(42), r.seed)
	require.NotNil(t, r.rnd)
	require.Equal(t, int64(675), r.rnd.Int63n(2000))
}

func Test_random_seedGetter(t *testing.T) {
	t.Parallel()

	r := newRandom(42)

	require.Equal(t, int64(42), r.seedGetter())

	r = newRandom(1)

	require.Equal(t, int64(1), r.seedGetter())
}

func Test_random_seedSetter(t *testing.T) {
	t.Parallel()

	r := newRandom(1)

	r.seedSetter(42)

	require.Equal(t, int64(42), r.seedGetter())
	require.Equal(t, int64(675), r.rnd.Int63n(2000))
}

func Test_random_intMethod(t *testing.T) {
	t.Parallel()

	toValue := sobek.New().ToValue

	r := newRandom(42)

	require.Equal(t, int64(675), r.intMethod(toValue(2000)))

	r.seedSetter(42)
	require.Equal(t, int64(8836438174961104), r.intMethod(nil))

	r.seedSetter(42)
	require.Equal(t, int64(8836438174961104), r.intMethod(sobek.NaN()))

	r.seedSetter(42)
	require.Equal(t, int64(8836438174961104), r.intMethod(sobek.Undefined()))
}

func Test_random_floatMethod(t *testing.T) {
	t.Parallel()

	toValue := sobek.New().ToValue

	r := newRandom(42)

	require.InEpsilon(t, float64(746), r.floatMethod(toValue(2000)), float64(1))

	r.seedSetter(42)
	require.InEpsilon(t, float64(.37), r.floatMethod(nil), float64(.01))

	r.seedSetter(42)
	require.InEpsilon(t, float64(.37), r.floatMethod(sobek.NaN()), float64(.01))

	r.seedSetter(42)
	require.InEpsilon(t, float64(.37), r.floatMethod(sobek.Undefined()), float64(.01))
}

func Test_module_random(t *testing.T) {
	t.Parallel()

	mod := &module{vu: modulestest.NewRuntime(t).VU}
	rt := mod.vu.Runtime()

	require.NoError(t, rt.Set("Random", mod.random))

	val, err := rt.RunString("new Random(42).int(2000)")

	require.NoError(t, err)
	require.Equal(t, int64(675), val.ToInteger())
}
