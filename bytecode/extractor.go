package bytecode

import "io"
import "fmt"

type extractor struct {
	w   Writer
	ex  *extractWriter
	abc AbcFile
}

type extractWriter struct {
	w   io.Writer
	err error
	n   int
}

func (w *extractWriter) Write(b []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.w.Write(b)
	w.err = err
	w.n += n
	return n, err
}

// Extract is used to serialize an AbcFile
func Extract(w io.Writer, abc AbcFile) error {
	wrappedWriter := &extractWriter{w: w}
	ex := extractor{NewWriter(wrappedWriter), wrappedWriter, abc}
	ex.Extract()
	return wrappedWriter.err
}

func (e *extractor) Extract() {
	e.w.WriteU16(e.abc.MinorVersion)
	e.w.WriteU16(e.abc.MajorVersion)
	e.extractCpool()
	fmt.Printf("cpool end: %x\n", e.ex.n)
	e.extractMethods()
	fmt.Printf("methods end: %x\n", e.ex.n)
	e.extractMetadatas()
	fmt.Printf("metadatas end: %x\n", e.ex.n)
	e.extractInstancesClasses()
	fmt.Printf("instances/classes end: %x\n", e.ex.n)
	e.extractScripts()
	fmt.Printf("scripts end: %x\n", e.ex.n)
	e.extractMethodBodies()
}

func (e *extractor) extractCpoolInt() {
	ints := e.abc.ConstantPool.Integers
	e.w.WriteU30(uint32(len(ints)))
	for i := 1; i < len(ints); i++ {
		e.w.WriteS32(ints[i])
	}
}

func (e *extractor) extractCpoolUInt() {
	uints := e.abc.ConstantPool.UIntegers
	e.w.WriteU30(uint32(len(uints)))
	for i := 1; i < len(uints); i++ {
		e.w.WriteU32(uints[i])
	}
}

func (e *extractor) extractCpoolDouble() {
	doubles := e.abc.ConstantPool.Doubles
	e.w.WriteU30(uint32(len(doubles)))
	for i := 1; i < len(doubles); i++ {
		e.w.WriteD64(doubles[i])
	}
}

func (e *extractor) extractCpoolString() {
	strs := e.abc.ConstantPool.Strings
	e.w.WriteU30(uint32(len(strs)))
	for i := 1; i < len(strs); i++ {
		e.extractStringInfo(strs[i])
	}
}

func (e *extractor) extractStringInfo(v string) {
	e.w.WriteU30(uint32(len(v)))
	e.w.Write([]byte(v))
}

func (e *extractor) extractCpoolNamespace() {
	namespaces := e.abc.ConstantPool.Namespaces
	e.w.WriteU30(uint32(len(namespaces)))
	for i := 1; i < len(namespaces); i++ {
		e.extractNamespaceInfo(namespaces[i])
	}
}

func (e *extractor) extractNamespaceInfo(v NamespaceInfo) {
	e.w.WriteU8(v.Kind)
	e.w.WriteU30(v.Name)
}

func (e *extractor) extractCpoolNsSet() {
	nsSets := e.abc.ConstantPool.NsSets
	e.w.WriteU30(uint32(len(nsSets)))
	for i := 1; i < len(nsSets); i++ {
		e.extractNsSetInfo(nsSets[i])
	}
}

func (e *extractor) extractNsSetInfo(v NsSetInfo) {
	count := uint32(len(v.Namespaces))
	e.w.WriteU30(count)
	for _, ns := range v.Namespaces {
		e.w.WriteU30(ns)
	}
}

func (e *extractor) extractCpoolMultiname() {
	multinames := e.abc.ConstantPool.Multinames
	e.w.WriteU30(uint32(len(multinames)))
	for i := 1; i < len(multinames); i++ {
		e.extractMultinameInfo(multinames[i])
	}
}

func (e *extractor) extractMultinameInfo(v MultinameInfo) {
	e.w.WriteU8(v.Kind)

	extractors := map[uint8]func(MultinameInfo){
		MultinameKindQName: e.extractQName, MultinameKindQNameA: e.extractQName,
		MultinameKindRTQName: e.extractRTQName, MultinameKindRTQNameA: e.extractRTQName,
		MultinameKindRTQNameL: e.extractRTQNameL, MultinameKindRTQNameLA: e.extractRTQNameL,
		MultinameKindMultiname: e.extractMultiname, MultinameKindMultinameA: e.extractMultiname,
		MultinameKindMultinameL: e.extractMultinameL, MultinameKindMultinameLA: e.extractMultinameL,
		MultinameKindTypename: e.extractTypename,
	}

	extract, ok := extractors[v.Kind]
	if !ok {
		panic("unknown multiname kind")
	}
	extract(v)
}

func (e *extractor) extractQName(v MultinameInfo) {
	e.w.WriteU30(v.Namespace)
	e.w.WriteU30(v.Name)
}

func (e *extractor) extractRTQName(v MultinameInfo) {
	e.w.WriteU30(v.Name)
}

func (e *extractor) extractRTQNameL(v MultinameInfo) {
}

func (e *extractor) extractMultiname(v MultinameInfo) {
	e.w.WriteU30(v.Name)
	e.w.WriteU30(v.NsSet)
}

func (e *extractor) extractMultinameL(v MultinameInfo) {
	e.w.WriteU30(v.NsSet)
}

func (e *extractor) extractTypename(v MultinameInfo) {
	e.w.WriteU30(v.Name)
	e.w.WriteU30(uint32(len(v.Params)))
	for _, param := range v.Params {
		e.w.WriteU30(param)
	}
}

func (e *extractor) extractCpool() {
	e.extractCpoolInt()
	e.extractCpoolUInt()
	e.extractCpoolDouble()
	e.extractCpoolString()
	e.extractCpoolNamespace()
	e.extractCpoolNsSet()
	e.extractCpoolMultiname()
}

func (e *extractor) extractOptionInfo(v OptionInfo) {
	e.w.WriteU30(uint32(len(v.Options)))
	for _, detail := range v.Options {
		e.w.WriteU30(detail.Value)
		e.w.WriteU8(detail.Kind)
	}
}

func (e *extractor) extractParamInfo(v ParamInfo) {
	for _, name := range v.ParamNames {
		e.w.WriteU30(name)
	}
}

func (e *extractor) extractMethod(v MethodInfo) {
	e.w.WriteU30(uint32(len(v.ParamTypes)))
	e.w.WriteU30(v.ReturnType)
	for _, paramType := range v.ParamTypes {
		e.w.WriteU30(paramType)
	}
	e.w.WriteU30(v.Name)
	e.w.WriteU8(v.Flags)
	if v.Flags&MethodHasOptional != 0 {
		e.extractOptionInfo(v.OptionInfo)
	}
	if v.Flags&MethodHasParamNames != 0 {
		e.extractParamInfo(v.ParamInfo)
	}
}

func (e *extractor) extractMethods() {
	methods := e.abc.Methods
	e.w.WriteU30(uint32(len(methods)))
	for _, method := range methods {
		e.extractMethod(method)
	}
}

func (e *extractor) extractMetadata(v MetadataInfo) {
	e.w.WriteU30(v.Names)
	items := v.Items
	e.w.WriteU30(uint32(len(items)))
	for _, item := range items {
		e.w.WriteU30(item.Key)
		e.w.WriteU30(item.Value)
	}
}

func (e *extractor) extractMetadatas() {
	metadatas := e.abc.Metadatas
	e.w.WriteU30(uint32(len(metadatas)))
	for _, metadata := range metadatas {
		e.extractMetadata(metadata)
	}
}

func (e *extractor) extractTrait(v TraitsInfo) {
	e.w.WriteU30(v.Name)
	e.w.WriteU8(v.Kind)

	switch v.Kind & 0xf {
	default:
		panic("unknown traits info")
	case TraitsInfoSlot, TraitsInfoConst:
		e.w.WriteU30(v.SlotID)
		e.w.WriteU30(v.Typename)
		e.w.WriteU30(v.VIndex)
		if v.VIndex != 0 {
			e.w.WriteU8(v.VKind)
		}
	case TraitsInfoClass:
		e.w.WriteU30(v.SlotID)
		e.w.WriteU30(v.ClassI)
	case TraitsInfoFunction:
		e.w.WriteU30(v.SlotID)
		e.w.WriteU30(v.Function)
	case TraitsInfoMethod, TraitsInfoGetter, TraitsInfoSetter:
		e.w.WriteU30(v.DispID)
		e.w.WriteU30(v.Method)
	}
	if v.Kind&TraitsInfoAttributeMetadata != 0 {
		metadatas := v.Metadatas
		e.w.WriteU30(uint32(len(metadatas)))
		for _, metadata := range metadatas {
			e.w.WriteU30(metadata)
		}
	}
}

func (e *extractor) extractTraits(v []TraitsInfo) {
	e.w.WriteU30(uint32(len(v)))
	for _, trait := range v {
		e.extractTrait(trait)
	}
}

func (e *extractor) extractInstance(v InstanceInfo) {
	e.w.WriteU30(v.Name)
	e.w.WriteU30(v.SuperName)
	e.w.WriteU8(v.Flags)

	if v.Flags&InstanceInfoClassProtectedNs != 0 {
		e.w.WriteU30(v.ProtectedNs)
	}

	e.w.WriteU30(uint32(len(v.Interfaces)))
	for _, intf := range v.Interfaces {
		e.w.WriteU30(intf)
	}

	e.w.WriteU30(v.IInit)
	e.extractTraits(v.Traits)
}

func (e *extractor) extractClass(v ClassInfo) {
	e.w.WriteU30(v.CInit)
	e.extractTraits(v.Traits)
}

func (e *extractor) extractInstancesClasses() {
	instances := e.abc.Instances
	classes := e.abc.Classes
	e.w.WriteU30(uint32(len(instances)))
	for _, instance := range instances {
		e.extractInstance(instance)
	}
	for _, class := range classes {
		e.extractClass(class)
	}
}

func (e *extractor) extractScript(v ScriptInfo) {
	e.w.WriteU30(v.Init)
	e.extractTraits(v.Traits)
}

func (e *extractor) extractScripts() {
	scripts := e.abc.Scripts
	e.w.WriteU30(uint32(len(scripts)))
	for _, script := range scripts {
		e.extractScript(script)
	}
}

func (e *extractor) extractExceptionInfo(v ExceptionInfo) {
	e.w.WriteU30(v.From)
	e.w.WriteU30(v.To)
	e.w.WriteU30(v.Target)
	e.w.WriteU30(v.ExcType)
	e.w.WriteU30(v.VarName)
}

func (e *extractor) extractMethodBody(v MethodBodyInfo) {
	e.w.WriteU30(v.Method)
	e.w.WriteU30(v.MaxStack)
	e.w.WriteU30(v.LocalCount)
	e.w.WriteU30(v.InitScopeLength)
	e.w.WriteU30(v.MaxScopeLength)
	e.w.WriteU30(uint32(len(v.Code)))
	e.w.Write(v.Code)

	e.w.WriteU30(uint32(len(v.Exceptions)))
	for _, exception := range v.Exceptions {
		e.extractExceptionInfo(exception)
	}
	e.extractTraits(v.Traits)
}

func (e *extractor) extractMethodBodies() {
	methodBodies := e.abc.MethodBodies
	e.w.WriteU30(uint32(len(methodBodies)))
	for _, methodBody := range methodBodies {
		e.extractMethodBody(methodBody)
	}
}
