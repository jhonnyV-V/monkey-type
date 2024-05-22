package code

import (
	"fmt"
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]

	if !ok {
		return nil, fmt.Errorf("opcode %d undefined\n", op)
	}
	return def, nil
}

type Instructions []byte
type Opcode byte

const (
	OpConstant Opcode = iota
)