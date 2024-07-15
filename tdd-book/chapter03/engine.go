package calculator

type Engine struct{}

//go:noinline
func (e *Engine) Add(x, y float64) float64 {
	return x + y
}

type Op interface {
}
