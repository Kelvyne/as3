package as3

import (
	"os"
	"reflect"
	"testing"

	"github.com/kelvyne/as3/bytecode"
)

func TestAbcFile_GetClassByName(t *testing.T) {
	f, err := os.Open("./bytecode/fixtures/frame1.abc")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	abc, err := bytecode.Parse(bytecode.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	l, err := Link(&abc)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  Class
		want1 bool
	}{
		{"first", args{l.Classes[0].Name}, l.Classes[0], true},
		{"not found", args{"unknoooooown class"}, Class{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := l
			got, got1 := f.GetClassByName(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AbcFile.GetClassByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AbcFile.GetClassByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
