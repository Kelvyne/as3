package bytecode

import (
	"fmt"
	"os"
	"testing"
)

func openFixture(t *testing.T, name string) *os.File {
	file, err := os.Open(fmt.Sprintf("./fixtures/%v.abc", name))
	if err != nil {
		t.Fatalf("openFixture: %v", err)
	}
	return file
}

func TestParse(t *testing.T) {
	file := openFixture(t, "frame1")
	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("file.close: %v", err)
		}
	}()
	reader := NewReader(file)
	a, err := Parse(reader)
	if err != nil {
		t.Fatalf("expected non-nil, got %v", err)
	}

	cpool := &a.ConstantPool

	if len(cpool.Integers) != 2926 {
		t.Errorf("expected 2926, got %v", len(cpool.Integers))
	}
	if len(cpool.UIntegers) != 15 {
		t.Errorf("expected 15, got %v", len(cpool.UIntegers))
	}
	if len(cpool.Doubles) != 857 {
		t.Errorf("expected 857, got %v", len(cpool.Doubles))
	}
	if len(cpool.Strings) != 51453 {
		t.Errorf("expected 51453, got %v", len(cpool.Strings))
	}
	if len(cpool.Namespaces) != 11524 {
		t.Errorf("expected 11524, got %v", len(cpool.Namespaces))
	}
	if len(cpool.NsSets) != 2226 {
		t.Errorf("expected 2226, got %v", len(cpool.NsSets))
	}
	if len(cpool.Multinames) != 48187 {
		t.Errorf("expected 48187, got %v", len(cpool.Multinames))
	}
	if len(a.Methods) != 46498 {
		t.Errorf("expected 46498, got %v", len(a.Methods))
	}
	if len(a.Metadatas) != 104 {
		t.Errorf("expected 104, got %v", len(a.Metadatas))
	}
	if len(a.Instances) != 5143 {
		t.Errorf("expected 5143, got %v", len(a.Instances))
	}
	if len(a.Scripts) != 4442 {
		t.Errorf("expected 4442, got %v", len(a.Scripts))
	}
	if len(a.MethodBodies) != 45243 {
		t.Errorf("expected 45243, got %v", len(a.MethodBodies))
	}
}
