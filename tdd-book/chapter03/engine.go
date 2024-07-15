package calculator

type Engine struct{}

type Operation struct {
	Expression string
	Operator   string
	Operands   []float64
}

//go:noinline
func (e *Engine) Add(x, y float64) float64 {
	return x + y
}
