package bytecode

import "testing"

func TestMethodBodyInfo_Disassemble(t *testing.T) {
	t.Skip("skipping test because some instructions are not implemented")

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
	if len(a.MethodBodies) != 45243 {
		t.Errorf("expected 45243, got %v", len(a.MethodBodies))
	}

	for i, body := range a.MethodBodies {
		if err = body.Disassemble(); err != nil {
			t.Errorf("method_body %v: %v", i, err)
		}
	}
}
