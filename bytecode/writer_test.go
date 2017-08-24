package bytecode

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_writer_WriteU8(t *testing.T) {
	type args struct {
		x uint8
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"valid",
			args{0x58},
			[]byte{0x58},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteU8(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteU8() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteU8() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}

func Test_writer_WriteU16(t *testing.T) {
	type args struct {
		x uint16
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"single",
			args{0x58},
			[]byte{0x58, 0x00},
			false,
		},
		{
			"double bytes",
			args{0x1234},
			[]byte{0x34, 0x12},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteU16(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteU16() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteU16() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}

func Test_writer_WriteS24(t *testing.T) {
	type args struct {
		x int32
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"positive",
			args{494682},
			[]byte{0x5a, 0x8c, 0x07},
			false,
		},
		{
			"negative",
			args{-494682},
			[]byte{0xa6, 0x73, 0xf8},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteS24(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteS24() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteS24() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}

func Test_writer_WriteU30(t *testing.T) {
	type args struct {
		x uint32
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"single byte",
			args{0x5a},
			[]byte{0x5a},
			false,
		},
		{
			"multi bytes",
			args{0x1ff},
			[]byte{0xff, 0x03},
			false,
		},
		{
			"five bytes",
			args{1274779536},
			[]byte{0x90, 0xaf, 0xee, 0xdf, 0x04},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteU30(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteU30() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteU30() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}

func Test_writer_WriteU32(t *testing.T) {
	type args struct {
		x uint32
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"single byte",
			args{0x5a},
			[]byte{0x5a},
			false,
		},
		{
			"multi bytes",
			args{0x1ff},
			[]byte{0xff, 0x03},
			false,
		},
		{
			"five bytes",
			args{1274779536},
			[]byte{0x90, 0xaf, 0xee, 0xdf, 0x04},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteU32(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteU32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteU32() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}

func Test_writer_WriteS32(t *testing.T) {
	type args struct {
		x int32
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"positive",
			args{0x3a},
			[]byte{0x3a},
			false,
		},
		{
			"multi bytes",
			args{0x1ff},
			[]byte{0xff, 0x03},
			false,
		},
		{
			"five bytes",
			args{1274779536},
			[]byte{0x90, 0xaf, 0xee, 0xdf, 0x04},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteS32(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteS32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteS32() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}

func Test_writer_WriteD64(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name        string
		args        args
		wantedBytes []byte
		wantErr     bool
	}{
		{
			"valid",
			args{-3.7079989838049655e-249},
			[]byte{0x5a, 0x8c, 0x5a, 0x8c, 0x5a, 0x8c, 0x5a, 0x8c},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &writer{
				Writer: buf,
			}
			if err := w.WriteD64(tt.args.x); (err != nil) != tt.wantErr {
				t.Errorf("writer.WriteD64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.wantedBytes) {
				t.Errorf("writer.WriteD64() wantedBytes %v, actualBytes %v", tt.wantedBytes, buf.Bytes())
			}
		})
	}
}
