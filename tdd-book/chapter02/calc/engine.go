package calc

type Engine struct{}

//go:noinline
func (e *Engine) Add(x, y float64) float64 {
	return x + y
}

//go:noinline
func (e *Engine) Subtract(x, y float64) float64 {
	return x - y
}

//go:noinline
func (e *Engine) Double(x float64) float64 {
	return x + x
}

//go:noinline
func (e *Engine) Double2(x float64) float64 {
	return 2 * x
}
