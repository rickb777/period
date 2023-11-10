package period

import (
	"github.com/govalues/decimal"
	. "github.com/onsi/gomega"
	"testing"
)

func TestSet(t *testing.T) {
	g := NewGomegaWithT(t)

	p := Period{}
	err := p.Set("P2D")

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(p).To(Equal(Period{days: decimal.Two}))

	g.Expect(p.Get()).To(Equal("P2D"))

	err = p.Set("Foo")

	g.Expect(err).To(HaveOccurred())

	typ := p.Type()

	g.Expect(typ).To(Equal("period"))
}
