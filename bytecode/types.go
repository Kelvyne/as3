package bytecode

// These constants are possible value of the kind field of a MultinameInfo struct
const (
	MultinameKindQName       = 0x07
	MultinameKindQNameA      = 0x0d
	MultinameKindRTQName     = 0x0f
	MultinameKindRTQNameA    = 0x10
	MultinameKindRTQNameL    = 0x11
	MultinameKindRTQNameLA   = 0x12
	MultinameKindMultiname   = 0x09
	MultinameKindMultinameA  = 0x0e
	MultinameKindMultinameL  = 0x1b
	MultinameKindMultinameLA = 0x1c
	MultinameKindTypename    = 0x1d
)

// These constants are possible flags of the flags field of a MethodInfo struct
const (
	MethodNeedArguments = 1 << iota
	MethodNeedActivation
	MethodNeedRest
	MethodHasOptional
	MethodSetDxns       = 0x40
	MethodHasParamNames = 0x80
)

// These constants are possible flags of the flags field of a InstanceInfo
// struct
const (
	InstanceInfoClassSealed = 1 << iota
	InstanceInfoClassFinal
	InstanceInfoClassInterface
	InstanceInfoClassProtectedNs
)

// These constants are possible values of the kind field of a TraitsInfo struct
const (
	TraitsInfoSlot = iota
	TraitsInfoMethod
	TraitsInfoGetter
	TraitsInfoSetter
	TraitsInfoClass
	TraitsInfoFunction
	TraitsInfoConst
)

// These constants are possible flags of the kind field of a TraitsInfo struct
const (
	TraitsInfoAttributeFinal = 1 << (4 + iota)
	TraitsInfoAttributeOverride
	TraitsInfoAttributeMetadata
)

// AbcFile is the root structure of an AS3 file
type AbcFile struct {
	MinorVersion uint16
	MajorVersion uint16
	ConstantPool CpoolInfo
	Methods      []MethodInfo
	Metadatas    []MetadataInfo
	Instances    []InstanceInfo
	Classes      []ClassInfo
	Scripts      []ScriptInfo
	MethodBodies []MethodBodyInfo
}

// CpoolInfo represents the constant pool informations of an AbcFile
type CpoolInfo struct {
	Integers   []int32
	UIntegers  []uint32
	Doubles    []float64
	Strings    []string
	Namespaces []NamespaceInfo
	NsSets     []NsSetInfo
	Multinames []MultinameInfo
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
type MultinameInfo struct {
	Kind      uint8
	Name      uint32   // reserved for QName, RTQName, Multiname and Typename
	Namespace uint32   // reserved for Namespace
	NsSet     uint32   // reserved for Multiname and MultinameL
	Params    []uint32 // reserved for Typename
}

// MethodInfo represents a method_info data structure
type MethodInfo struct {
	ParamCount uint32
	ReturnType uint32
	ParamTypes []uint32
	Name       uint32
	Flags      uint8
	OptionInfo OptionInfo
	ParamInfo  ParamInfo
}

// OptionInfo represents a option_info data structure
type OptionInfo struct {
	Options []OptionDetail
}

// OptionDetail represents a option_detail data structure
type OptionDetail struct {
	Value uint32
	Kind  uint8
}

// ParamInfo represents a param_info data structure
type ParamInfo struct {
	ParamNames []uint32
}

// MetadataInfo represents a metadata_info data structure
type MetadataInfo struct {
	Names uint32
	Items []ItemInfo
}

// ItemInfo represents a item_info data structure
type ItemInfo struct {
	Key   uint32
	Value uint32
}

// TraitsInfo represents a traits_info data structure
type TraitsInfo struct {
	Name      uint32
	Kind      uint8
	SlotID    uint32 // Reserved for slot, constant, class and function
	Typename  uint32 // Reserved for slot and constant
	VIndex    uint32 // Reserved for slot and constant
	VKind     uint8  // Reserved for slot and constant
	ClassI    uint32 // Reserved for class
	Function  uint32 // Reserved for function
	DispID    uint32 // Reserved for method, getter and setter
	Method    uint32 // Reserved for method, getter and setter
	Metadatas []uint32
}

// GetType isolates the type (4 lower bits) of a TraitsInfo.
func (t TraitsInfo) GetType() uint8 { return t.Kind & 0x0f }

// InstanceInfo represents a instance_info data structure
type InstanceInfo struct {
	Name        uint32
	SuperName   uint32
	Flags       uint8
	ProtectedNs uint32
	Interfaces  []uint32
	IInit       uint32
	Traits      []TraitsInfo
}

// ClassInfo represents a class_info data structure
type ClassInfo struct {
	CInit  uint32
	Traits []TraitsInfo
}

// ScriptInfo represents a script_info data structure
type ScriptInfo struct {
	Init   uint32
	Traits []TraitsInfo
}

// MethodBodyInfo represents a method_body_info data structure
type MethodBodyInfo struct {
	Method          uint32
	MaxStack        uint32
	LocalCount      uint32
	InitScopeLength uint32
	MaxScopeLength  uint32
	Code            []byte
	Exceptions      []ExceptionInfo
	Traits          []TraitsInfo
	Instructions    []Instr
}

// ExceptionInfo represents a exception_info data structure
type ExceptionInfo struct {
	From    uint32
	To      uint32
	Target  uint32
	ExcType uint32
	VarName uint32
}
