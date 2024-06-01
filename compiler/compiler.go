package compiler

import (
	"fmt"
	"mokey-type/ast"
	"mokey-type/code"
	"mokey-type/object"
)

type Compiler struct {
	instructions code.Instructions
	constanst    []object.Object
}

type Bytecode struct {
	Instructions code.Instructions
	Constanst    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constanst:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)

		case "*":
			c.emit(code.OpMul)

		case "/":
			c.emit(code.OpDiv)

		case ">":
			c.emit(code.OpGreaterThan)

		case "==":
			c.emit(code.OpEqual)

		case "!=":
			c.emit(code.OpNotEqual)

		default:
			return fmt.Errorf("unknow operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)

		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operato %s", node.Operator)
		}
	}
	return nil
}

func (c *Compiler) addConstant(ob object.Object) int {
	c.constanst = append(c.constanst, ob)
	return len(c.constanst) - 1
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constanst:    c.constanst,
	}
}