package period

import (
	"github.com/govalues/decimal"
	"github.com/rickb777/expect"
	"testing"
)

func TestSet(t *testing.T) {
	p := Period{}
	err := p.Set("P2D")

	expect.Error(err).Not().ToHaveOccurred(t)
	expect.Any(p).ToBe(t, Period{days: decimal.Two})

	expect.Any(p.Get()).ToBe(t, "P2D")

	err = p.Set("Foo")

	expect.Error(err).ToHaveOccurred(t)

	typ := p.Type()

	expect.String(typ).ToBe(t, "period")
}
