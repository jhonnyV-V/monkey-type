package vm

import (
	"fmt"
	"mokey-type/code"
	"mokey-type/compiler"
	"mokey-type/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.NullValue{}

type VM struct {
	constant []object.Object

	stack []object.Object
	sp    int //Always points to the next value. Top of the stack is stack[sp - 1]

	globals []object.Object

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	return &VM{
		constant: bytecode.Constanst,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	return &VM{
		constant: bytecode.Constanst,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: globals,

		frames:      frames,
		framesIndex: 1,
	}
}

func (vm *VM) LastPopedStackElement() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	var ip int
	var op code.Opcode
	var ins code.Instructions

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions()) {
		vm.currentFrame().ip++
		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		//TODO: check how I messed up this indexing
		op = code.Opcode(ins[ip-1])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip:])
			vm.currentFrame().ip += 2
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
			pos := int(code.ReadUint16(ins[ip:]))
			vm.currentFrame().ip = pos

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()

			if !isTruthy(condition) {
				vm.currentFrame().ip = pos
			}

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip:])
			vm.currentFrame().ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpArray:
			numOfElements := code.ReadUint16(ins[ip:])
			vm.currentFrame().ip += 2

			newSp := vm.sp - int(numOfElements)
			array := vm.buildArray(newSp, vm.sp)
			vm.sp = newSp

			err := vm.push(array)
			if err != nil {
				return err
			}

		case code.OpHash:
			numOfElements := code.ReadUint16(ins[ip:])
			vm.currentFrame().ip += 2

			newSp := vm.sp - int(numOfElements)

			hash, err := vm.buildHash(newSp, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = newSp

			err = vm.push(hash)
			if err != nil {
				return err
			}

		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case code.OpCall:
			numArgs := int(code.ReadUint8(ins[ip:]))
			vm.currentFrame().ip += 1

			err := vm.executeCall(numArgs)
			if err != nil {
				return err
			}

		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetLocal:
			localIndex := int(code.ReadUint8(ins[ip:]))
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			vm.stack[frame.basePointer+localIndex] = vm.pop()

		case code.OpGetLocal:
			localIndex := int(code.ReadUint8(ins[ip:]))
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			err := vm.push(vm.stack[frame.basePointer+localIndex])
			if err != nil {
				return err
			}

		case code.OpGetBuiltin:
			localIndex := int(code.ReadUint8(ins[ip:]))
			vm.currentFrame().ip += 1

			definition := object.Builtins[localIndex]
			err := vm.push(definition.Builtin)

			if err != nil {
				return err
			}

		case code.OpClosure:
			constIndex := int(code.ReadUint16(ins[ip:]))
			numFree := int(code.ReadUint8(ins[ip+2:]))
			vm.currentFrame().ip += 3

			err := vm.pushClosure(constIndex, numFree)
			if err != nil {
				return err
			}

		case code.OpGetFree:
			freeIndex := int(code.ReadUint8(ins[ip:]))
			vm.currentFrame().ip += 1
			currentCl := vm.currentFrame().cl

			err := vm.push(currentCl.Free[freeIndex])
			if err != nil {
				return err
			}

		case code.OpCurrentClosure:
			currentCl := vm.currentFrame().cl

			err := vm.push(currentCl)
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
	case False:
		return vm.push(True)
	case Null:
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

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Value: value, Key: key}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return nil, fmt.Errorf("unusable as hashkey: %s", key.Type())
		}

		pairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: pairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)

	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)

	case left.Type() == object.HASH_OBJ && index.Type() != object.INTEGER_OBJ:
		return fmt.Errorf("unsuported index of array %s", index.Type())

	default:
		return fmt.Errorf("index operation not suported %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
	array := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)
	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(array.Elements[i])
}

func (vm *VM) executeHashIndex(left, index object.Object) error {
	hash := left.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return fmt.Errorf("unable to hash key %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]

	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(frame *Frame) {
	vm.frames[vm.framesIndex] = frame
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--

	return vm.frames[vm.framesIndex]
}

func (vm *VM) callClosure(cl *object.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d", cl.Fn.NumParameters, numArgs)
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)
	vm.sp = frame.basePointer + cl.Fn.NumLocals
	return nil
}

func (vm *VM) executeCall(numArg int) error {
	callee := vm.stack[vm.sp-1-numArg]
	switch callee := callee.(type) {
	case *object.Closure:
		return vm.callClosure(callee, numArg)

	case *object.Builtin:
		return vm.callBuiltin(callee, numArg)

	default:
		return fmt.Errorf("calling non-function")
	}
}

func (vm *VM) callBuiltin(fn *object.Builtin, numArg int) error {
	args := vm.stack[vm.sp-numArg : vm.sp]
	result := fn.Fn(args...)

	if result != nil {
		switch result := result.(type) {

		case *object.Boolean:
			if result.Value {
				return vm.push(True)
			}
			return vm.push(False)

		default:
			return vm.push(result)

		}
	} else {
		return vm.push(Null)
	}
}

func (vm *VM) pushClosure(constIndex, numFree int) error {
	constant := vm.constant[constIndex]

	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	free := make([]object.Object, numFree)

	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}

	vm.sp -= numFree

	closure := &object.Closure{Fn: function, Free: free}
	return vm.push(closure)
}
