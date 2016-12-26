package as3

import (
	"github.com/kelvyne/as3/bytecode"
)

// AbcFile represents a linked AbcFile
type AbcFile struct {
	src     *bytecode.AbcFile
	Classes []Class
	Methods []Method
}

// Class represents an actionscript Class
type Class struct {
	InstanceInfo   bytecode.InstanceInfo
	ClassInfo      bytecode.ClassInfo
	Name           string
	SuperName      string
	Interfaces     []string
	InstanceTraits TraitsObject
	ClassTraits    TraitsObject
}

// TraitsObject represents an object that has traits
type TraitsObject struct {
	Slots     []Trait
	Classes   []Trait
	Functions []Trait
	Methods   []Trait
}

// Trait represents a single trait from a TraitsObject
type Trait struct {
	Source   bytecode.TraitsInfo
	Name     string
	Typename string
}

// Method represents a linked method
type Method struct {
	Info       bytecode.MethodInfo
	BodyInfo   bytecode.MethodBodyInfo
	HasBody    bool
	Name       string
	ReturnType string
	ParamTypes []string
}
