package vm

import (
	"fmt"
	"mokey-type/code"
	"mokey-type/compiler"
	"mokey-type/object"
)

const StackSize = 2048
const GlobalsSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.NullValue{}

type VM struct {
	constant     []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int //Always points to the next value. Top of the stack is stack[sp - 1]

	globals []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constant:     bytecode.Constanst,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	return &VM{
		constant:     bytecode.Constanst,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: globals,
	}
}

func (vm *VM) LastPopedStackElement() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constant[constIndex])
			if err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return nil
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return nil
			}

		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			condition := vm.pop()

			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) push(ob object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("Stack Overflow")
	}

	vm.stack[vm.sp] = ob
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	ob := vm.stack[vm.sp-1]
	vm.sp--

	return ob
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	if rightType == object.INTEGER_OBJ && leftType == object.INTEGER_OBJ {
		return vm.executeIntegerBinaryOperation(op, left, right)
	}

	if rightType == object.STRING_OBJ && leftType == object.STRING_OBJ {
		return vm.executeStringBinaryOperation(op, left, right)
	}
	return fmt.Errorf("unsoported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeIntegerBinaryOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result int64
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue

	case code.OpSub:
		result = leftValue - rightValue

	case code.OpMul:
		result = leftValue * rightValue

	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknow integer operation: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeStringBinaryOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknow integer operation: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	left := vm.pop()
	right := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	var result bool
	switch op {
	case code.OpEqual:
		result = left == right

	case code.OpNotEqual:
		result = left != right
	default:
		return fmt.Errorf("unknow operator: %d (%s %s)", op, left.Type(), right.Type())
	}

	return vm.push(nativeBooleanObject(result))
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result bool
	switch op {
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue

	case code.OpGreaterThan:
		result = rightValue > leftValue
	default:
		return fmt.Errorf("unknown operator %d", op)
	}

	return vm.push(nativeBooleanObject(result))
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False, Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsuported type for negation: %s", operand.Type())
	}
	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func nativeBooleanObject(value bool) *object.Boolean {
	if value {
		return True
	} else {
		return False
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value

	case *object.NullValue:
		return false

	default:
		return true
	}
}
