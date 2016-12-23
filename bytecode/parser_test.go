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
	reader := NewReader(openFixture(t, "frame1"))
	a, err := Parse(reader)
	if err != nil {
		t.Fatalf("expected non-nil, got %+v", err)
	}
	fmt.Printf("Major %v Minor %v\n", a.MajorVersion, a.MinorVersion)

	cpool := &a.ConstantPool

	if cpool.IntCount != 2926 {
		t.Errorf("expected 2926, got %v", cpool.IntCount)
	}
	if cpool.UIntCount != 15 {
		t.Errorf("expected 15, got %v", cpool.UIntCount)
	}
	if cpool.DoubleCount != 857 {
		t.Errorf("expected 857, got %v", cpool.DoubleCount)
	}
	if cpool.StringCount != 51453 {
		t.Errorf("expected 51453, got %v", cpool.StringCount)
	}
	if cpool.NamespaceCount != 11524 {
		t.Errorf("expected 11524, got %v", cpool.NamespaceCount)
	}
	if cpool.NsSetCount != 2226 {
		t.Errorf("expected 2226, got %v", cpool.NsSetCount)
	}
	if cpool.MultinameCount != 48187 {
		t.Errorf("expected 48187, got %v", cpool.MultinameCount)
	}
}
