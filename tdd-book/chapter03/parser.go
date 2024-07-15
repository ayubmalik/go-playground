package calculator

import "fmt"

type OperationProcessor interface {
	ProcessOp(op any) (*string, error)
}

type Validator interface {
	Validator(operator string, operands []float64) error
}

type Parser struct {
	engine    OperationProcessor
	validator Validator
}

func (p *Parser) ProcessExpr(expr string) (*string, error) {
	op, err := p.getOp(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid expression '%s': %v", expr, err)
	}
	return p.engine.ProcessOp(op)
}

func (p *Parser) getOp(expr string) (*Op, interface{}) {
	return nil, nil
}
