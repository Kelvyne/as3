package bytecode

import "testing"
import "io/ioutil"
import "bytes"
import "reflect"

func TestExtractor(t *testing.T) {
	file := openFixture(t, "frame1")
	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("file.close: %v", err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll failed: %v\n", err)
	}

	reader := NewReader(bytes.NewReader(b))
	a, err := Parse(reader)
	if err != nil {
		t.Fatalf("expected non-nil, got %v", err)
	}

	buffer := bytes.Buffer{}

	err = Extract(&buffer, a)
	if err != nil {
		t.Fatalf("Extract failed: %v\n", err)
	}

	if !reflect.DeepEqual(buffer.Bytes(), b) {
		t.Errorf("buffer are not equal (required len: %v, got %v)", len(b), buffer.Len())
	}
}
