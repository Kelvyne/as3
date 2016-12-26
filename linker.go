package as3

import "github.com/kelvyne/as3/bytecode"
import "errors"

// ErrLinkerUnknownTrait means that an unknown trait_info kind was found when
// build a TraitsObject
var ErrLinkerUnknownTrait = errors.New("linker unknown trait")

type linker struct {
	abc *bytecode.AbcFile
}

// Link produces a linked version of the AbcFile provided. It :
// - Links instance_info and class_info and resolve informations about each class
// - Resolve method names, parameters and return types
// and link method_body_info when the method has a body
func Link(abcFile *bytecode.AbcFile) (AbcFile, error) {
	l := linker{abcFile}
	return l.Link()
}

func (l *linker) Link() (AbcFile, error) {
	classes, err := l.LinkClasses()
	if err != nil {
		return AbcFile{}, err
	}
	methods, err := l.LinkMethods()
	if err != nil {
		return AbcFile{}, err
	}
	return AbcFile{l.abc, classes, methods}, nil
}

func (l *linker) LinkClasses() ([]Class, error) {
	classes := make([]Class, len(l.abc.Classes))
	for i := range classes {
		class, err := l.LinkClass(i)
		if err != nil {
			return nil, err
		}
		classes[i] = class
	}
	return classes, nil
}

func (l *linker) LinkClass(index int) (c Class, err error) {
	c.InstanceInfo = l.abc.Instances[index]
	c.ClassInfo = l.abc.Classes[index]
	c.Name = l.abc.ConstantPool.MultinameString(c.InstanceInfo.Name)
	c.SuperName = l.abc.ConstantPool.MultinameString(c.InstanceInfo.SuperName)
	c.Interfaces = make([]string, len(c.InstanceInfo.Interfaces))
	for i := range c.Interfaces {
		c.Interfaces[i] = l.abc.ConstantPool.MultinameString(c.InstanceInfo.Interfaces[i])
	}
	instanceTraits, err := l.BuildTraits(c.InstanceInfo.Traits)
	if err != nil {
		c = Class{}
		return
	}
	c.InstanceTraits = instanceTraits
	classTraits, err := l.BuildTraits(c.ClassInfo.Traits)
	if err != nil {
		c = Class{}
		return
	}
	c.ClassTraits = classTraits
	return
}

func (l *linker) BuildTraits(info []bytecode.TraitsInfo) (TraitsObject, error) {
	o := TraitsObject{}
	mapping := map[uint8]*[]Trait{
		bytecode.TraitsInfoSlot:     &o.Slots,
		bytecode.TraitsInfoConst:    &o.Slots,
		bytecode.TraitsInfoClass:    &o.Classes,
		bytecode.TraitsInfoFunction: &o.Functions,
		bytecode.TraitsInfoMethod:   &o.Methods,
		bytecode.TraitsInfoSetter:   &o.Methods,
		bytecode.TraitsInfoGetter:   &o.Methods,
	}
	for i := range info {
		t := info[i].GetType()
		arrayPtr, ok := mapping[t]
		if !ok {
			return TraitsObject{}, ErrLinkerUnknownTrait
		}
		name := l.abc.ConstantPool.MultinameString(info[i].Name)
		var typename string
		if t == bytecode.TraitsInfoSlot || t == bytecode.TraitsInfoConst {
			typename = l.abc.ConstantPool.MultinameString(info[i].Typename)
		}
		*arrayPtr = append(*arrayPtr, Trait{info[i], name, typename})
	}
	return o, nil
}

func (l *linker) LinkMethods() ([]Method, error) {
	methods := make([]Method, len(l.abc.Methods))
	cpool := &l.abc.ConstantPool
	for i := range methods {
		info := l.abc.Methods[i]
		name := cpool.MultinameString(info.Name)
		returnType := cpool.MultinameString(info.ReturnType)
		paramTypes := make([]string, len(info.ParamTypes))
		for iParam := range paramTypes {
			paramTypes[iParam] = cpool.MultinameString(info.ParamTypes[iParam])
		}
		methods[i] = Method{
			info, bytecode.MethodBodyInfo{}, false,
			name, returnType, paramTypes,
		}
	}

	for i := range l.abc.MethodBodies {
		info := l.abc.MethodBodies[i]
		methods[info.Method].HasBody = true
		methods[info.Method].BodyInfo = info
	}
	return methods, nil
}
