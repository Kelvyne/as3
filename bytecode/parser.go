package bytecode

import "errors"

type parser struct {
	r Reader
}

// ErrUnknownMultinameKind means that an unknown multiname was found in the constant pool
var ErrUnknownMultinameKind = errors.New("unknown multiname kind")

// Parse parses an AS3 file
func Parse(r Reader) (AbcFile, error) {
	p := parser{r}
	return p.Parse()
}

func (p *parser) Parse() (abcFile AbcFile, err error) {
	minor, err := p.r.ReadU16()
	if err != nil {
		return
	}
	abcFile.MinorVersion = minor

	major, err := p.r.ReadU16()
	if err != nil {
		return
	}
	abcFile.MajorVersion = major

	cpoolInfo, err := p.ParseCpool()
	if err != nil {
		return
	}
	abcFile.ConstantPool = cpoolInfo
	return
}
