package tlvrw

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	type args struct {
		tlv *TLV
	}
	tests := []struct {
		name     string
		args     args
		wantW    string
		wantErr  bool
		prehook  func()
		posthook func()
	}{
		{
			name: "Valid TLV write",
			args: args{
				tlv: NewTLV(0, []byte("hello")),
			},
			wantW:   "\x00\x05\x00\x00\x00hello",
			wantErr: false,
		},
		{
			name: "Valid TLV write - 2",
			args: args{
				tlv: NewTLV(1, []byte("hello")),
			},
			wantW:   "\x01\x05\x00\x00\x00hello",
			wantErr: false,
		},
		{
			name: "Invalid TLV write - len too long",
			args: args{
				tlv: NewTLV(0, make([]byte, 6)),
			},
			wantW:   "",
			wantErr: true,
			prehook: func() {
				MaxLen = 5
			},
			posthook: func() {
				MaxLen = math.MaxUint32
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prehook != nil {
				tt.prehook()
			}
			if tt.posthook != nil {
				defer tt.posthook()
			}

			w := &bytes.Buffer{}
			if err := Write(w, tt.args.tlv); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Write() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *TLV
		wantErr bool
	}{
		{
			name: "Valid TLV read",
			args: args{
				r: bytes.NewBuffer([]byte("\x00\x05\x00\x00\x00hello")),
			},
			want:    NewTLV(0, []byte("hello")),
			wantErr: false,
		},
		{
			name: "Valid TLV read - 2",
			args: args{
				r: bytes.NewBuffer([]byte("\x02\x05\x00\x00\x00hello")),
			},
			want:    NewTLV(2, []byte("hello")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
