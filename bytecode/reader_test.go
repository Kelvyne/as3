package bytecode

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNewReader(t *testing.T) {
	r := NewReader(nil)
	if r == nil {
		t.Fatalf("expected non-nil, got %v", r)
	}
}

func Test_reader_ReadU8(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint8
		wantErr bool
	}{
		{
			"valid",
			fields{bytes.NewReader([]byte{0x5a})},
			0x5a,
			false,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadU8()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadU8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadU8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadU16(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint16
		wantErr bool
	}{
		{
			"valid",
			fields{bytes.NewReader([]byte{0x5a, 0x8c})},
			0x8c5a,
			false,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadU16()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadU16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadU16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadS24(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    int32
		wantErr bool
	}{
		{
			"positive",
			fields{bytes.NewReader([]byte{0x5a, 0x8c, 0x07})},
			494682,
			false,
		},
		{
			"negative",
			fields{bytes.NewReader([]byte{0xa6, 0x73, 0xf8})},
			-494682,
			false,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadS24()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadS24() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadS24() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_reader_ReadU32(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint32
		wantErr bool
	}{
		{
			"valid",
			fields{bytes.NewReader([]byte{0x5a})},
			0x5a,
			false,
		},
		{
			"multi bytes",
			fields{bytes.NewReader([]byte{0xff, 0x03})},
			0x1ff,
			false,
		},
		{
			"malformed variable length",
			fields{bytes.NewReader([]byte{0x81, 0x82, 0x83, 0x84, 0x85, 0x86})},
			0,
			true,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadU32()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadU32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadU32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadU30(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint32
		wantErr bool
	}{
		{
			"valid",
			fields{bytes.NewReader([]byte{0x5a})},
			0x5a,
			false,
		},
		{
			"multi bytes",
			fields{bytes.NewReader([]byte{0xff, 0x03})},
			0x1ff,
			false,
		},
		{
			"malformed variable length",
			fields{bytes.NewReader([]byte{0x81, 0x82, 0x83, 0x84, 0x85, 0x86})},
			0,
			true,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadU30()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadU30() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadU30() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadS32(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    int32
		wantErr bool
	}{
		{
			"positive",
			fields{bytes.NewReader([]byte{0x3a})},
			0x3a,
			false,
		},
		{
			"multi bytes",
			fields{bytes.NewReader([]byte{0xff, 0x03})},
			0x1ff,
			false,
		},
		{
			"negative",
			fields{bytes.NewReader([]byte{0x81, 0x80, 0x80, 0x80, 0x08})},
			-1,
			false,
		},
		{
			"five bytes",
			fields{bytes.NewReader([]byte{0x90, 0xaf, 0xee, 0xdf, 0x04})},
			1274779536,
			false,
		},
		{
			"malformed variable length",
			fields{bytes.NewReader([]byte{0x81, 0x82, 0x83, 0x84, 0x85, 0x86})},
			0,
			true,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadS32()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadS32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadS32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadD64(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{
			"valid",
			fields{bytes.NewReader([]byte{0x5a, 0x8c, 0x5a, 0x8c, 0x5a, 0x8c, 0x5a, 0x8c})},
			-3.7079989838049655e-249,
			false,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadD64()
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadD64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.ReadD64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadBytes(t *testing.T) {
	type fields struct {
		Reader io.Reader
	}
	type args struct {
		n uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"complete",
			fields{bytes.NewReader([]byte{0x01, 0x02, 0x03})},
			args{3},
			[]byte{0x01, 0x02, 0x03},
			false,
		},
		{
			"incomplete",
			fields{bytes.NewReader([]byte{0x01, 0x02})},
			args{3},
			nil,
			true,
		},
		{
			"EOF",
			fields{&bytes.Buffer{}},
			args{3},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				Reader: tt.fields.Reader,
			}
			got, err := r.ReadBytes(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reader.ReadBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
