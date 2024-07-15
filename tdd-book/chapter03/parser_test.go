package calculator_test

import (
	calculator "github.com/tdd-book/chapter03"
	"github.com/tdd-book/chapter03/mocks"
	"testing"
)

func TestParser_ProcessExpr(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {

		expr := "2 + 3"
		operator := "+"
		operands := []float64{2, 3}
		want := "2 + 3 = 5.5"

		engine := mocks.NewOperationProcessor(t)
		validator := mocks.NewValidator(t)

		validator.On("validate", operator, operands).Return(nil).Once()

		engine.On("ProcessOp", &calculator.Operation{
			Expression: expr,
			Operator:   operator,
			Operands:   operands,
		}).Return(want).Once()
	})
}
