package bytecode

import "bytes"
import "io"
import "errors"
import "fmt"

// InstrOperand is the type of an operand
type InstrOperand uint8

// These are possible operand types
const (
	InstrOperandU30 = InstrOperand(iota)
	InstrOperandU8
	InstrOperandS24
	InstrOperandCaseCount
)

// InstrModel represents the model for an avm2 instruction
type InstrModel struct {
	Code     uint8
	Name     string
	Operands []InstrOperand
}

// Instr represents a disassembled avm2 instruction.
// Since all the operands are at most 30 bits, operands are stored as uint32.
type Instr struct {
	Model    InstrModel
	Operands []uint32
}

// ErrUnknownInstruction means that an invalid instruction code was read
// when disassembling a method.
var ErrUnknownInstruction = errors.New("unknown instruction")

// ErrUnknownInstructionOperand means that an invalid instruction operand
// was found in the Instructions map
var ErrUnknownInstructionOperand = errors.New("unknown instruction operand")

// Instructions is a map of Instruction. Each key represents the instruction code
var Instructions = map[uint8]InstrModel{
	0xa0: {0xa0, "add", nil},
	0xc5: {0xc5, "add_i", nil},
	0x86: {0x86, "astype", []InstrOperand{InstrOperandU30}},
	0x87: {0x87, "astypelate", nil},
	0xa8: {0xa8, "bitand", nil},
	0x97: {0x97, "bitnot", nil},
	0xa9: {0xa9, "bitor", nil},
	0xaa: {0xaa, "bitxor", nil},
	0x41: {0x41, "call", []InstrOperand{InstrOperandU30}},
	0x43: {0x43, "callmethod", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x46: {0x46, "callproperty", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x4c: {0x4c, "callproplex", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x4f: {0x4f, "callpropvoid", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x44: {0x44, "callstatic", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x45: {0x45, "callsuper", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x4e: {0x4e, "callsupervoid", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x78: {0x78, "checkfilter", nil},
	0x80: {0x80, "coerce", []InstrOperand{InstrOperandU30}},
	0x82: {0x82, "coerce_a", nil},
	0x85: {0x85, "coerce_s", nil},
	0x42: {0x42, "construct", []InstrOperand{InstrOperandU30}},
	0x4a: {0x4a, "contructprop", []InstrOperand{InstrOperandU30, InstrOperandU30}},
	0x49: {0x49, "constructsuper", []InstrOperand{InstrOperandU30}},
	0x76: {0x76, "convert_b", nil},
	0x73: {0x73, "convert_i", nil},
	0x75: {0x75, "convert_d", nil},
	0x77: {0x77, "convert_o", nil},
	0x74: {0x74, "convert_u", nil},
	0x70: {0x70, "convert_s", nil},
	0xef: {0xef, "debug", []InstrOperand{InstrOperandU8, InstrOperandU30, InstrOperandU8, InstrOperandU30}},
	0xf1: {0xf1, "debugfile", []InstrOperand{InstrOperandU30}},
	0xf0: {0xf0, "debugline", []InstrOperand{InstrOperandU30}},
	0x94: {0x94, "declocal", []InstrOperand{InstrOperandU30}},
	0xc3: {0xc3, "declocal_i", []InstrOperand{InstrOperandU30}},
	0x93: {0x93, "decrement", nil},
	0xc1: {0xc1, "decrement_i", nil},
	0x6a: {0x6a, "deleteproperty", []InstrOperand{InstrOperandU30}},
	0xa3: {0xa3, "divide", nil},
	0x2a: {0x2a, "dup", nil},
	0x06: {0x06, "dxns", []InstrOperand{InstrOperandU30}},
	0x07: {0x07, "dxnslate", nil},
	0xab: {0xab, "equals", nil},
	0x72: {0x72, "esc_xattr", nil},
	0x71: {0x71, "esc_xelem", nil},
	0x5e: {0x5e, "findproperty", []InstrOperand{InstrOperandU30}},
	0x5d: {0x5d, "findpropstrict", []InstrOperand{InstrOperandU30}},
	0x59: {0x59, "getdescendants", []InstrOperand{InstrOperandU30}},
	0x64: {0x64, "getglobalscope", nil},
	0x6e: {0x6e, "getglobalslot", []InstrOperand{InstrOperandU30}},
	0x60: {0x60, "getlex", []InstrOperand{InstrOperandU30}},
	0x62: {0x62, "getlocal", []InstrOperand{InstrOperandU30}},
	0xd0: {0xd0, "getlocal_0", nil},
	0xd1: {0xd1, "getlocal_1", nil},
	0xd2: {0xd2, "getlocal_2", nil},
	0xd3: {0xd3, "getlocal_3", nil},
	0x66: {0x66, "getproperty", []InstrOperand{InstrOperandU30}},
	0x65: {0x65, "getscopeobject", []InstrOperand{InstrOperandU30}},
	0x6c: {0x6c, "getslot", []InstrOperand{InstrOperandU30}},
	0x04: {0x04, "getsuper", []InstrOperand{InstrOperandU30}},
	0xb0: {0xb0, "greaterequals", nil},
	0xaf: {0xaf, "greaterthan", nil},
	0x1f: {0x1f, "hasnext", nil},
	0x32: {0x32, "hasnext2", []InstrOperand{InstrOperandU30}},
	0x13: {0x13, "ifeq", []InstrOperand{InstrOperandS24}},
	0x12: {0x12, "iffalse", []InstrOperand{InstrOperandS24}},
	0x18: {0x18, "ifge", []InstrOperand{InstrOperandS24}},
	0x17: {0x17, "ifgt", []InstrOperand{InstrOperandS24}},
	0x16: {0x16, "ifle", []InstrOperand{InstrOperandS24}},
	0x15: {0x15, "iflt", []InstrOperand{InstrOperandS24}},
	0x0f: {0x0f, "ifnge", []InstrOperand{InstrOperandS24}},
	0x0e: {0x0e, "ifngt", []InstrOperand{InstrOperandS24}},
	0x0d: {0x0d, "ifnle", []InstrOperand{InstrOperandS24}},
	0x0c: {0x0c, "ifnlt", []InstrOperand{InstrOperandS24}},
	0x14: {0x14, "ifne", []InstrOperand{InstrOperandS24}},
	0x19: {0x19, "ifstricteq", []InstrOperand{InstrOperandS24}},
	0x1a: {0x1a, "ifstrictne", []InstrOperand{InstrOperandS24}},
	0x11: {0x11, "iftrue", []InstrOperand{InstrOperandS24}},
	0xb4: {0xb4, "in", nil},
	0x92: {0x92, "inclocal", []InstrOperand{InstrOperandU30}},
	0xc2: {0xc2, "inclocal_i", []InstrOperand{InstrOperandU30}},
	0x91: {0x91, "increment", nil},
	0xc0: {0xc0, "increment_i", nil},
	0x68: {0x68, "initproperty", []InstrOperand{InstrOperandU30}},
	0xb1: {0xb1, "instanceof", nil},
	0xb2: {0xb2, "istype", []InstrOperand{InstrOperandU30}},
	0xb3: {0xb3, "istypelate", nil},
	0x10: {0x10, "jump", []InstrOperand{InstrOperandS24}},
	0x08: {0x08, "kill", []InstrOperand{InstrOperandU30}},
	0x09: {0x09, "label", nil},
	0xae: {0xae, "lessequals", nil},
	0xad: {0xad, "lessthan", nil},
	0x1b: {0x1b, "lookupswitch", []InstrOperand{InstrOperandS24, InstrOperandCaseCount}},
	0xa5: {0xa5, "lshift", nil},
	0xa4: {0xa4, "modulo", nil},
	0xa2: {0xa2, "multiply", nil},
	0xc7: {0xc7, "multiply_i", nil},
	0x90: {0x90, "negate", nil},
	0xc4: {0xc4, "negate_i", nil},
	0x57: {0x57, "newactivation", nil},
	0x56: {0x56, "newarray", []InstrOperand{InstrOperandU30}},
	0x5a: {0x5a, "newcatch", []InstrOperand{InstrOperandU30}},
	0x58: {0x58, "newclass", []InstrOperand{InstrOperandU30}},
	0x40: {0x40, "newfunction", []InstrOperand{InstrOperandU30}},
	0x55: {0x55, "newobject", []InstrOperand{InstrOperandU30}},
	0x1e: {0x1e, "nextname", nil},
	0x23: {0x23, "nextvalue", nil},
	0x02: {0x02, "nop", nil},
	0x96: {0x96, "not", nil},
	0x29: {0x29, "pop", nil},
	0x1d: {0x1d, "popscope", nil},
	0x24: {0x24, "pushbyte", []InstrOperand{InstrOperandU8}},
	0x2f: {0x2f, "pushdouble", []InstrOperand{InstrOperandU30}},
	0x27: {0x27, "pushfalse", nil},
	0x2d: {0x2d, "pushint", []InstrOperand{InstrOperandU30}},
	0x31: {0x31, "pushnamespace", []InstrOperand{InstrOperandU30}},
	0x28: {0x28, "pushnan", nil},
	0x20: {0x20, "pushnull", nil},
	0x30: {0x30, "pushscope", nil},
	0x25: {0x25, "pushshort", []InstrOperand{InstrOperandU30}},
	0x2c: {0x2c, "pushstring", []InstrOperand{InstrOperandU30}},
	0x26: {0x26, "pushtrue", nil},
	0x2e: {0x2e, "pushuint", []InstrOperand{InstrOperandU30}},
	0x21: {0x21, "undefined", nil},
	0x1c: {0x1c, "pushwith", nil},
	0x48: {0x48, "returnvalue", nil},
	0x47: {0x47, "returnvoid", nil},
	0xa6: {0xa6, "rshift", nil},
	0x63: {0x63, "setlocal", []InstrOperand{InstrOperandU30}},
	0xd4: {0xd4, "setlocal_0", nil},
	0xd5: {0xd5, "setlocal_1", nil},
	0xd6: {0xd6, "setlocal_2", nil},
	0xd7: {0xd7, "setlocal_3", nil},
	0x6f: {0x6f, "setglobalslot", []InstrOperand{InstrOperandU30}},
	0x61: {0x61, "setproperty", []InstrOperand{InstrOperandU30}},
	0x6d: {0x6d, "setslot", []InstrOperand{InstrOperandU30}},
	0x05: {0x05, "setsuper", []InstrOperand{InstrOperandU30}},
	0xac: {0xac, "strictequals", nil},
	0xa1: {0xa1, "subtract", nil},
	0xc6: {0xc6, "subtract_i", nil},
	0x2b: {0x2b, "swap", nil},
	0x03: {0x03, "throw", nil},
	0x95: {0x95, "typeof", nil},
	0xa7: {0xa7, "urshift", nil},
}

func disassembleInstrOperand(r Reader, t InstrOperand) (uint32, error) {
	switch t {
	case InstrOperandU30:
		v, err := r.ReadU30()
		if err != nil {
			return 0, err
		}
		return v, nil
	case InstrOperandS24:
		v, err := r.ReadS24()
		if err != nil {
			return 0, err
		}
		return uint32(v), nil
	case InstrOperandU8:
		v, err := r.ReadU8()
		if err != nil {
			return 0, err
		}
		return uint32(v), nil
	}
	return 0, ErrUnknownInstructionOperand
}

func dissassembleInstr(r Reader, code uint8) (Instr, error) {
	model, ok := Instructions[code]
	if !ok {
		return Instr{}, fmt.Errorf("unknown instruction %v", code)
	}
	var operands []uint32
	for _, t := range model.Operands {
		if t == InstrOperandCaseCount {
			count, err := r.ReadU30()
			if err != nil {
				return Instr{}, err
			}
			for i := uint32(0); i < count; i++ {
				v, err := disassembleInstrOperand(r, InstrOperandS24)
				if err != nil {
					return Instr{}, err
				}
				operands = append(operands, v)
			}
		} else {
			v, err := disassembleInstrOperand(r, t)
			if err != nil {
				return Instr{}, err
			}
			operands = append(operands, v)
		}
	}
	return Instr{model, operands}, nil
}

// Disassemble parses the instructions of the method body
func (m *MethodBodyInfo) Disassemble() (err error) {
	base := bytes.NewReader(m.Code)
	r := NewReader(base)
	var instructions []Instr
	for {
		var code uint8
		var instr Instr
		code, err = r.ReadU8()
		if err != nil {
			break
		}
		instr, err = dissassembleInstr(r, code)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			break
		}
		instructions = append(instructions, instr)

	}
	if err != nil && err != io.EOF {
		return
	}
	m.Instructions = instructions
	err = nil
	return
}
