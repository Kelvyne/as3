package bytecode

import "errors"

type parser struct {
	r Reader
}

// ErrUnknownTraitsInfoKind means that an unknown traits_info was found in classes
var ErrUnknownTraitsInfoKind = errors.New("unknown traits_info kind")

// ErrUnknownMultinameKind means that an unknown multiname was found in the constant pool
var ErrUnknownMultinameKind = errors.New("unknown multiname kind")

// Parse parses an AS3 file
func Parse(r Reader) (AbcFile, error) {
	p := parser{r}
	return p.Parse()
}

func (p *parser) Parse() (AbcFile, error) {
	minor, err := p.r.ReadU16()
	if err != nil {
		return AbcFile{}, err
	}

	major, err := p.r.ReadU16()
	if err != nil {
		return AbcFile{}, err
	}

	cpoolInfo, err := p.ParseCpool()
	if err != nil {
		return AbcFile{}, err
	}

	methods, err := p.ParseMethods()
	if err != nil {
		return AbcFile{}, err
	}

	metadatas, err := p.ParseMetadatas()
	if err != nil {
		return AbcFile{}, err
	}

	instances, classes, err := p.ParseInstancesClasses()
	if err != nil {
		return AbcFile{}, err
	}

	scripts, err := p.ParseScripts()
	if err != nil {
		return AbcFile{}, err
	}

	methodBodies, err := p.ParseMethodBodies()
	if err != nil {
		return AbcFile{}, err
	}

	return AbcFile{minor, major, cpoolInfo, methods, metadatas, instances, classes, scripts, methodBodies}, nil
}

func (p *parser) parseCpoolInt() (slice []int32, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]int32, n)
	for i := uint32(1); i < n; i++ {
		v, vErr := p.r.ReadS32()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = v
	}
	return
}

func (p *parser) parseCpoolUInt() (slice []uint32, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]uint32, n)
	for i := uint32(1); i < n; i++ {
		v, vErr := p.r.ReadU32()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = v
	}
	return
}

func (p *parser) parseCpoolDouble() (slice []float64, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]float64, n)
	for i := uint32(1); i < n; i++ {
		v, vErr := p.r.ReadD64()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = v
	}
	return
}

func (p *parser) parseCpoolString() (slice []string, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]string, n)
	for i := uint32(1); i < n; i++ {
		str, vErr := p.parseStringInfo()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = str
	}
	return
}

func (p *parser) parseStringInfo() (s string, err error) {
	length, err := p.r.ReadU30()
	if err != nil {
		return
	}
	bytes, err := p.r.ReadBytes(length)
	if err != nil {
		return
	}
	s = string(bytes)
	return
}

func (p *parser) parseCpoolNamespace() (slice []NamespaceInfo, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]NamespaceInfo, n)
	for i := uint32(1); i < n; i++ {
		ns, vErr := p.parseNamespaceInfo()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = ns
	}
	return
}

func (p *parser) parseNamespaceInfo() (NamespaceInfo, error) {
	kind, err := p.r.ReadU8()
	if err != nil {
		return NamespaceInfo{}, err
	}
	name, err := p.r.ReadU30()
	if err != nil {
		return NamespaceInfo{}, err
	}
	return NamespaceInfo{kind, name}, nil
}

func (p *parser) parseCpoolNsSet() (slice []NsSetInfo, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]NsSetInfo, n)
	for i := uint32(1); i < n; i++ {
		nsSetInfo, vErr := p.parseNsSetInfo()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = nsSetInfo
	}
	return
}

func (p *parser) parseNsSetInfo() (NsSetInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return NsSetInfo{}, err
	}
	namespaces := make([]uint32, count)
	for i := range namespaces {
		ns, err := p.r.ReadU30()
		if err != nil {
			return NsSetInfo{}, err
		}
		namespaces[i] = ns
	}
	return NsSetInfo{count, namespaces}, nil
}

func (p *parser) parseCpoolMultiname() (slice []MultinameInfo, err error) {
	n, err := p.r.ReadU30()
	if err != nil {
		return
	}

	slice = make([]MultinameInfo, n)
	for i := uint32(1); i < n; i++ {
		multinameInfo, vErr := p.parseMultinameInfo()
		if vErr != nil {
			err = vErr
			return
		}
		slice[i] = multinameInfo
	}
	return
}

func (p *parser) parseMultinameInfo() (MultinameInfo, error) {
	kind, err := p.r.ReadU8()
	if err != nil {
		return nil, ErrUnknownMultinameKind
	}

	parsers := map[uint8]func(multinameInfo) (MultinameInfo, error){
		KindQName: p.parseQName, KindQNameA: p.parseQName,
		KindRTQName: p.parseRTQName, KindRTQNameA: p.parseRTQName,
		KindRTQNameL: p.parseRTQNameL, KindRTQNameLA: p.parseRTQNameL,
		KindMultiname: p.parseMultiname, KindMultinameA: p.parseMultiname,
		KindMultinameL: p.parseMultinameL, KindMultinameLA: p.parseMultinameL,
		KindTypename: p.parseTypename,
	}
	parser, ok := parsers[kind]
	if !ok {
		return nil, ErrUnknownMultinameKind
	}
	b := multinameInfo{kind}
	mInfo, err := parser(b)
	if err != nil {
		return nil, err
	}
	return mInfo, nil
}

func (p *parser) parseQName(b multinameInfo) (MultinameInfo, error) {
	ns, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	name, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	return QName{b, ns, name}, nil
}

func (p *parser) parseRTQName(b multinameInfo) (MultinameInfo, error) {
	name, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	return RTQName{b, name}, nil
}

func (p *parser) parseRTQNameL(b multinameInfo) (MultinameInfo, error) {
	return RTQNameL{b}, nil
}

func (p *parser) parseMultiname(b multinameInfo) (MultinameInfo, error) {
	name, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	nsSet, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	return Multiname{b, name, nsSet}, nil
}

func (p *parser) parseMultinameL(b multinameInfo) (MultinameInfo, error) {
	nsSet, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	return MultinameL{b, nsSet}, nil
}

func (p *parser) parseTypename(b multinameInfo) (MultinameInfo, error) {
	name, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	paramLength, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	params := make([]uint32, paramLength)
	for i := range params {
		param, err := p.r.ReadU30()
		if err != nil {
			return nil, err
		}
		params[i] = param
	}
	return Typename{b, name, params}, nil
}

func (p *parser) ParseCpool() (CpoolInfo, error) {
	ints, err := p.parseCpoolInt()
	if err != nil {
		return CpoolInfo{}, err
	}
	uints, err := p.parseCpoolUInt()
	if err != nil {
		return CpoolInfo{}, err
	}
	doubles, err := p.parseCpoolDouble()
	if err != nil {
		return CpoolInfo{}, err
	}
	strings, err := p.parseCpoolString()

	namespaces, err := p.parseCpoolNamespace()
	if err != nil {
		return CpoolInfo{}, err
	}

	nsSets, err := p.parseCpoolNsSet()
	if err != nil {
		return CpoolInfo{}, err
	}

	multinames, err := p.parseCpoolMultiname()
	if err != nil {
		return CpoolInfo{}, err
	}
	return CpoolInfo{
		ints, uints, doubles, strings, namespaces, nsSets, multinames,
	}, nil
}

func (p *parser) ParseOptionInfo() (OptionInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return OptionInfo{}, err
	}
	optionDetails := make([]OptionDetail, count)
	for i := range optionDetails {
		val, err := p.r.ReadU30()
		if err != nil {
			return OptionInfo{}, err
		}
		kind, err := p.r.ReadU8()
		if err != nil {
			return OptionInfo{}, err
		}
		optionDetails[i] = OptionDetail{val, kind}
	}
	return OptionInfo{optionDetails}, nil
}

func (p *parser) ParseParamInfo(paramCount uint32) (ParamInfo, error) {
	paramNames := make([]uint32, paramCount)
	for i := range paramNames {
		paramName, err := p.r.ReadU30()
		if err != nil {
			return ParamInfo{}, err
		}
		paramNames[i] = paramName
	}
	return ParamInfo{paramNames}, nil
}

func (p *parser) ParseMethod() (MethodInfo, error) {
	paramCount, err := p.r.ReadU30()
	if err != nil {
		return MethodInfo{}, err
	}
	returnType, err := p.r.ReadU30()
	if err != nil {
		return MethodInfo{}, err
	}
	paramTypes := make([]uint32, paramCount)
	for i := range paramTypes {
		paramType, pErr := p.r.ReadU30()
		if pErr != nil {
			return MethodInfo{}, pErr
		}
		paramTypes[i] = paramType
	}
	name, err := p.r.ReadU30()
	if err != nil {
		return MethodInfo{}, err
	}
	flags, err := p.r.ReadU8()
	if err != nil {
		return MethodInfo{}, err
	}
	var optionInfo OptionInfo
	if flags&MethodHasOptional != 0 {
		optionInfo, err = p.ParseOptionInfo()
		if err != nil {
			return MethodInfo{}, err
		}
	}
	var paramInfo ParamInfo
	if flags&MethodHasParamNames != 0 {
		paramInfo, err = p.ParseParamInfo(paramCount)
		if err != nil {
			return MethodInfo{}, err
		}
	}
	return MethodInfo{paramCount, returnType, paramTypes, name, flags, optionInfo, paramInfo}, nil
}

func (p *parser) ParseMethods() ([]MethodInfo, error) {
	nMethods, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}

	methods := make([]MethodInfo, nMethods)
	for i := range methods {
		method, err := p.ParseMethod()
		if err != nil {
			return nil, err
		}
		methods[i] = method
	}
	return methods, nil
}

func (p *parser) ParseMetadata() (MetadataInfo, error) {
	name, err := p.r.ReadU30()
	if err != nil {
		return MetadataInfo{}, err
	}
	count, err := p.r.ReadU30()
	if err != nil {
		return MetadataInfo{}, err
	}

	items := make([]ItemInfo, count)
	for i := range items {
		key, err := p.r.ReadU30()
		if err != nil {
			return MetadataInfo{}, err
		}
		value, err := p.r.ReadU30()
		if err != nil {
			return MetadataInfo{}, err
		}
		items[i] = ItemInfo{key, value}
	}
	return MetadataInfo{name, items}, nil
}

func (p *parser) ParseMetadatas() ([]MetadataInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	metadatas := make([]MetadataInfo, count)
	for i := range metadatas {
		metadata, err := p.ParseMetadata()
		if err != nil {
			return nil, err
		}
		metadatas[i] = metadata
	}
	return metadatas, nil
}

func (p *parser) ParseTrait() (TraitsInfo, error) {
	var t TraitsInfo
	name, err := p.r.ReadU30()
	if err != nil {
		return TraitsInfo{}, err
	}
	t.Name = name

	kind, err := p.r.ReadU8()
	if err != nil {
		return TraitsInfo{}, err
	}
	t.Kind = kind

	switch kind & 0xf {
	default:
		return TraitsInfo{}, ErrUnknownTraitsInfoKind
	case TraitsInfoSlot, TraitsInfoConst:
		slotID, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		typename, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		vIndex, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		var vKind uint8
		if vIndex != 0 {
			vKind, tErr = p.r.ReadU8()
			if tErr != nil {
				return TraitsInfo{}, tErr
			}
		}
		t.SlotID = slotID
		t.Typename = typename
		t.VIndex = vIndex
		t.VKind = vKind
	case TraitsInfoClass:
		slotID, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		classI, tErr := p.r.ReadU30()
		if err != nil {
			return TraitsInfo{}, tErr
		}
		t.SlotID = slotID
		t.ClassI = classI
	case TraitsInfoFunction:
		slotID, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		function, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		t.SlotID = slotID
		t.Function = function
	case TraitsInfoMethod, TraitsInfoGetter, TraitsInfoSetter:
		dispID, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		method, tErr := p.r.ReadU30()
		if tErr != nil {
			return TraitsInfo{}, tErr
		}
		t.DispID = dispID
		t.Method = method
	}
	if kind&TraitsInfoAttributeMetadata != 0 {
		metadataCount, err := p.r.ReadU30()
		if err != nil {
			return TraitsInfo{}, err
		}
		t.Metadatas = make([]uint32, metadataCount)
		for i := range t.Metadatas {
			metadata, err := p.r.ReadU30()
			if err != nil {
				return TraitsInfo{}, err
			}
			t.Metadatas[i] = metadata
		}
	}

	return t, nil
}

func (p *parser) ParseTraits() ([]TraitsInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	traits := make([]TraitsInfo, count)
	for i := range traits {
		trait, err := p.ParseTrait()
		if err != nil {
			return nil, err
		}
		traits[i] = trait
	}
	return traits, nil
}

func (p *parser) ParseInstance() (InstanceInfo, error) {
	name, err := p.r.ReadU30()
	if err != nil {
		return InstanceInfo{}, nil
	}
	superName, err := p.r.ReadU30()
	if err != nil {
		return InstanceInfo{}, nil
	}
	flags, err := p.r.ReadU8()
	if err != nil {
		return InstanceInfo{}, nil
	}

	var protectedNs uint32
	if flags&InstanceInfoClassProtectedNs != 0 {
		protectedNs, err = p.r.ReadU30()
		if err != nil {
			return InstanceInfo{}, nil
		}
	}

	interfaceCount, err := p.r.ReadU30()
	if err != nil {
		return InstanceInfo{}, err
	}
	interfaces := make([]uint32, interfaceCount)
	for i := range interfaces {
		intrf, iErr := p.r.ReadU30()
		if iErr != nil {
			return InstanceInfo{}, iErr
		}
		interfaces[i] = intrf
	}
	iInit, err := p.r.ReadU30()
	if err != nil {
		return InstanceInfo{}, err
	}
	traits, err := p.ParseTraits()
	if err != nil {
		return InstanceInfo{}, err
	}
	return InstanceInfo{name, superName, flags, protectedNs, interfaces, iInit, traits}, nil
}

func (p *parser) ParseClass() (ClassInfo, error) {
	cinit, err := p.r.ReadU30()
	if err != nil {
		return ClassInfo{}, err
	}
	traits, err := p.ParseTraits()
	if err != nil {
		return ClassInfo{}, err
	}
	return ClassInfo{cinit, traits}, nil
}

func (p *parser) ParseInstancesClasses() ([]InstanceInfo, []ClassInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return nil, nil, err
	}
	instances := make([]InstanceInfo, count)
	for i := range instances {
		instance, iErr := p.ParseInstance()
		if iErr != nil {
			return nil, nil, iErr
		}
		instances[i] = instance
	}
	classes := make([]ClassInfo, count)
	for i := range instances {
		class, cErr := p.ParseClass()
		if cErr != nil {
			return nil, nil, cErr
		}
		classes[i] = class
	}
	return instances, classes, nil
}

func (p *parser) ParseScript() (ScriptInfo, error) {
	init, err := p.r.ReadU30()
	if err != nil {
		return ScriptInfo{}, err
	}
	traits, err := p.ParseTraits()
	if err != nil {
		return ScriptInfo{}, err
	}
	return ScriptInfo{init, traits}, nil
}

func (p *parser) ParseScripts() ([]ScriptInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	scripts := make([]ScriptInfo, count)
	for i := range scripts {
		script, sErr := p.ParseScript()
		if sErr != nil {
			return nil, sErr
		}
		scripts[i] = script
	}
	return scripts, nil
}

func (p *parser) ParseExceptionInfo() (ExceptionInfo, error) {
	from, err := p.r.ReadU30()
	if err != nil {
		return ExceptionInfo{}, err
	}
	to, err := p.r.ReadU30()
	if err != nil {
		return ExceptionInfo{}, err
	}
	target, err := p.r.ReadU30()
	if err != nil {
		return ExceptionInfo{}, err
	}
	excType, err := p.r.ReadU30()
	if err != nil {
		return ExceptionInfo{}, err
	}
	varName, err := p.r.ReadU30()
	if err != nil {
		return ExceptionInfo{}, err
	}
	return ExceptionInfo{from, to, target, excType, varName}, nil
}

func (p *parser) ParseMethodBody() (MethodBodyInfo, error) {
	method, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	maxStack, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	localCount, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	initScopeLength, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	maxScopeLength, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	codeLength, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	code, err := p.r.ReadBytes(codeLength)
	if err != nil {
		return MethodBodyInfo{}, err
	}
	exceptionLength, err := p.r.ReadU30()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	exceptions := make([]ExceptionInfo, exceptionLength)
	for i := range exceptions {
		exception, eErr := p.ParseExceptionInfo()
		if eErr != nil {
			return MethodBodyInfo{}, eErr
		}
		exceptions[i] = exception
	}
	traits, err := p.ParseTraits()
	if err != nil {
		return MethodBodyInfo{}, err
	}
	return MethodBodyInfo{method, maxStack, localCount, initScopeLength, maxScopeLength, code, exceptions, traits}, nil
}

func (p *parser) ParseMethodBodies() ([]MethodBodyInfo, error) {
	count, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	methodBodies := make([]MethodBodyInfo, count)
	for i := range methodBodies {
		methodBody, mErr := p.ParseMethodBody()
		if mErr != nil {
			return nil, mErr
		}
		methodBodies[i] = methodBody
	}
	return methodBodies, nil
}
