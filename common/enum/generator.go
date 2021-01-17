package enum

var gen *generator

func init() {
	gen = newGenerator()
	gen.prepare()
}

type generator struct {
	c chan Enum
}

func (g *generator) Next() Enum {
	return <-g.c
}

func (g *generator) prepare() {
	var count Enum
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
