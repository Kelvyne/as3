package bytecode

func (p *parser) parseCpoolInt() (n uint32, slice []int32, err error) {
	n, err = p.r.ReadU30()
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

func (p *parser) parseCpoolUInt() (n uint32, slice []uint32, err error) {
	n, err = p.r.ReadU30()
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

func (p *parser) parseCpoolDouble() (n uint32, slice []float64, err error) {
	n, err = p.r.ReadU30()
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

func (p *parser) parseCpoolString() (n uint32, slice []string, err error) {
	n, err = p.r.ReadU30()
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

func (p *parser) parseCpoolNamespace() (n uint32, slice []NamespaceInfo, err error) {
	n, err = p.r.ReadU30()
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

func (p *parser) parseCpoolNsSet() (n uint32, slice []NsSetInfo, err error) {
	n, err = p.r.ReadU30()
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

func (p *parser) parseCpoolMultiname() (n uint32, slice []MultinameInfo, err error) {
	n, err = p.r.ReadU30()
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
	name, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	ns, err := p.r.ReadU30()
	if err != nil {
		return nil, err
	}
	return QName{b, name, ns}, nil
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

func (p *parser) ParseCpool() (CpoolInfo, error) {
	nInt, ints, err := p.parseCpoolInt()
	if err != nil {
		return CpoolInfo{}, err
	}
	nUInt, uints, err := p.parseCpoolUInt()
	if err != nil {
		return CpoolInfo{}, err
	}
	nDouble, doubles, err := p.parseCpoolDouble()
	if err != nil {
		return CpoolInfo{}, err
	}
	nStrings, strings, err := p.parseCpoolString()

	nNamespaces, namespaces, err := p.parseCpoolNamespace()
	if err != nil {
		return CpoolInfo{}, err
	}

	nNsSets, nsSets, err := p.parseCpoolNsSet()
	if err != nil {
		return CpoolInfo{}, err
	}

	nMultinames, multinames, err := p.parseCpoolMultiname()

	return CpoolInfo{
		nInt, ints,
		nUInt, uints,
		nDouble, doubles,
		nStrings, strings,
		nNamespaces, namespaces,
		nNsSets, nsSets,
		nMultinames, multinames,
	}, nil
}
