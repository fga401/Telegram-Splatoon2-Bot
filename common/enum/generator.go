package enum

const (
	// EmptyValue is the value not used by assigner.
	EmptyValue = 0
)

var gen *generator

func init() {
	gen = newGenerator()
	gen.prepare()
}

type generator struct {
	c chan Enum
}

func (g *generator) next() Enum {
	return <-g.c
}

func (g *generator) prepare() {
	var count Enum = EmptyValue + 1
	go func() {
		for {
			g.c <- count
			count++
		}
	}()
}
func newGenerator() *generator {
	return &generator{c: make(chan Enum)}
}
