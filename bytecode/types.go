package bytecode

// AbcFile is the root structure of an AS3 file
type AbcFile struct {
	MinorVersion uint16
	MajorVersion uint16
	ConstantPool CpoolInfo
}

// These constants are possible value of the kind field of a MultinameInfo struct
const (
	KindQName       = 0x07
	KindQNameA      = 0x0d
	KindRTQName     = 0x0f
	KindRTQNameA    = 0x10
	KindRTQNameL    = 0x11
	KindRTQNameLA   = 0x12
	KindMultiname   = 0x09
	KindMultinameA  = 0x0e
	KindMultinameL  = 0x1b
	KindMultinameLA = 0x1c
)

// CpoolInfo represents the constant pool informations of an AbcFile
type CpoolInfo struct {
	IntCount       uint32
	Integers       []int32
	UIntCount      uint32
	UIntegers      []uint32
	DoubleCount    uint32
	Doubles        []float64
	StringCount    uint32
	Strings        []string
	NamespaceCount uint32
	Namespaces     []NamespaceInfo
	NsSetCount     uint32
	NsSets         []NsSetInfo
	MultinameCount uint32
	Multinames     []MultinameInfo
}

// NamespaceInfo represents a namespace info data structure
type NamespaceInfo struct {
	Kind uint8
	Name uint32
}

// NsSetInfo represents a namespace set info data structure
type NsSetInfo struct {
	Count      uint32
	Namespaces []uint32
}

// MultinameInfo represents a multiname info data structure
type MultinameInfo interface {
	Kind() uint8
}

type multinameInfo struct {
	kind uint8
}

func (m multinameInfo) Kind() uint8 {
	return m.kind
}

// QName represents a qualified name data structure
type QName struct {
	multinameInfo
	Namespace uint32
	Name      uint32
}

// RTQName represents a real time qualified name data structure
type RTQName struct {
	multinameInfo
	Name uint32
}

// RTQNameL represents a late real time qualified name data structure
type RTQNameL struct {
	multinameInfo
}

// Multiname represents a multiname data structure
type Multiname struct {
	multinameInfo
	Name  uint32
	NsSet uint32
}

// MultinameL represents a late multiname data structure
type MultinameL struct {
	multinameInfo
	NsSet uint32
}
