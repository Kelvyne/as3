package as3

import (
	"fmt"
	"os"
	"testing"

	"github.com/kelvyne/as3/bytecode"
)

func openFixture(t *testing.T, name string) *os.File {
	file, err := os.Open(fmt.Sprintf("./bytecode/fixtures/%v.abc", name))
	if err != nil {
		t.Fatalf("openFixture: %v", err)
	}
	return file
}

func TestLink(t *testing.T) {
	file := openFixture(t, "frame1")
	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("file.close: %v", err)
		}
	}()
	reader := bytecode.NewReader(file)
	a, err := bytecode.Parse(reader)
	if err != nil {
		t.Fatalf("expected non-nil, got %v", err)
	}

	l, err := Link(&a)
	if err != nil {
		t.Fatalf("expected non-nil, got %v", err)
	}

	if len(l.Classes) != len(a.Classes) {
		t.Errorf("expected %v, got %v", len(a.Classes), len(l.Classes))
	}

	if len(l.Methods) != len(a.Methods) {
		t.Errorf("expected %v, got %v", len(a.Methods), len(l.Methods))
	}
}
