package example

import (
	"math/rand"

	"github.com/grafana/sobek"
)

const maxSafeInteger int64 = 9007199254740991

type random struct {
	seed int64
	rnd  *rand.Rand
}

func newRandom(seed int64) *random {
	r := new(random)

	r.seed = seed
	r.rnd = rand.New(rand.NewSource(seed)) // #nosec G115,G404

	return r
}

func (r *random) seedGetter() int64 {
	return r.seed
}

func (r *random) seedSetter(seed int64) {
	r.rnd.Seed(seed)
	r.seed = seed
}

func (r *random) intMethod(n sobek.Value) int64 {
	var m int64

	if n == nil || sobek.IsNaN(n) || sobek.IsUndefined(n) {
		m = maxSafeInteger
	} else {
		m = min(maxSafeInteger, n.ToInteger())
	}

	return r.rnd.Int63n(m)
}

func (r *random) floatMethod(n sobek.Value) float64 {
	var m float64

	if n == nil || sobek.IsNaN(n) || sobek.IsUndefined(n) {
		m = 1
	} else {
		m = n.ToFloat()
	}

	return r.rnd.Float64() * m
}

func (m *module) random(call sobek.ConstructorCall) *sobek.Object {
	var seed int64

	if val := call.Argument(0); !sobek.IsNaN(val) && !sobek.IsUndefined(val) {
		seed = val.ToInteger()
	} else {
		seed = rand.Int63() // #nosec G404
	}

	r := newRandom(seed)

	this := call.This
	toValue := m.vu.Runtime().ToValue
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	must(this.DefineAccessorProperty(
		"seed",
		toValue(r.seedGetter),
		toValue(r.seedSetter),
		sobek.FLAG_FALSE,
		sobek.FLAG_TRUE,
	))

	must(this.Set("int", toValue(r.intMethod)))
	must(this.Set("float", toValue(r.floatMethod)))

	return nil
}
