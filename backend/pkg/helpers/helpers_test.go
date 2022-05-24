package helpers

import (
	"testing"

	"github.com/gofrs/uuid"
)

func Test_IsURL_IsHTTPURL(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		want     bool // result for IsURL
		wantHTTP bool // result for IsHTTPURL
	}{
		{
			name:     "empty string",
			arg:      "",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "invalid",
			arg:      "://",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "empty host x",
			arg:      "x://",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "empty host http",
			arg:      "http://",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "empty host https",
			arg:      "https://",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "non-http 1",
			arg:      "x://xxx",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "non-http 2",
			arg:      "aaa://xxx",
			want:     false,
			wantHTTP: false,
		},
		{
			name:     "one letter host 1",
			arg:      "http://x",
			want:     true,
			wantHTTP: true,
		},
		{
			name:     "one letter host 2",
			arg:      "https://x",
			want:     true,
			wantHTTP: false,
		},
		{
			name:     "full url 1",
			arg:      "https://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x1234",
			want:     true,
			wantHTTP: false,
		},
		{
			name:     "full url 2",
			arg:      "http://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x1234",
			want:     true,
			wantHTTP: true,
		},
		{
			name:     "full url non-http",
			arg:      "aaa://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x1234",
			want:     false,
			wantHTTP: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.arg); got != tt.want {
				t.Errorf("IsURL() = %v, want %v", got, tt.want)
			}
			if got := IsHTTPURL(tt.arg); got != tt.wantHTTP {
				t.Errorf("IsHTTPURL() = %v, wantHTTP %v", got, tt.wantHTTP)
			}
		})
	}
}

func TestRewriteToHTTPS(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "empty string",
			arg:  "",
			want: "",
		},
		{
			name: "invalid 1",
			arg:  ":",
			want: ":",
		},
		{
			name: "invalid 2",
			arg:  ":/",
			want: ":/",
		},
		{
			name: "invalid 3",
			arg:  "://",
			want: "://",
		},
		{
			name: "invalid 4",
			arg:  "x:",
			want: "x:",
		},
		{
			name: "invalid 5",
			arg:  "x://",
			want: "x://",
		},
		{
			name: "valid 1",
			arg:  "x://y",
			want: "https://y",
		},
		{
			name: "valid 2",
			arg:  "http://y",
			want: "https://y",
		},
		{
			name: "valid 3",
			arg:  "https://y",
			want: "https://y",
		},
		{
			name: "full url 1",
			arg:  "https://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x1234",
			want: "https://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x1234",
		},
		{
			name: "full url 2",
			arg:  "http://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x123",
			want: "https://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x123",
		},
		{
			name: "full url non-http",
			arg:  "aaa://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x12345",
			want: "https://foo.bar.example.com/baz/boo?d=12&a=b&c[f]=0x12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RewriteToHTTPS(tt.arg); got != tt.want {
				t.Errorf("RewriteToHTTPS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameUUID(t *testing.T) {
	type args struct {
		a uuid.UUID
		b uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same null uuid",
			args: args{
				a: uuid.NullUUID{}.UUID,
				b: uuid.NullUUID{}.UUID,
			},
			want: true,
		},
		{
			name: "same uuid",
			args: args{
				a: uuid.FromStringOrNil("9919a392-1548-4e13-b8d9-36f0159a10ad"),
				b: uuid.FromStringOrNil("9919a392-1548-4e13-b8d9-36f0159a10ad"),
			},
			want: true,
		},
		{
			name: "different uuid",
			args: args{
				a: uuid.FromStringOrNil("9919a392-1548-4e13-b8d9-36f0159a10ab"),
				b: uuid.FromStringOrNil("9919a392-1548-4e13-b8d9-36f0159a10ac"),
			},
			want: false,
		},
		{
			name: "different than null uuid 1",
			args: args{
				a: uuid.NullUUID{}.UUID,
				b: uuid.FromStringOrNil("9919a392-1548-4e13-b8d9-36f0159a10aa"),
			},
			want: false,
		},
		{
			name: "different than null uuid 2",
			args: args{
				a: uuid.FromStringOrNil("9919a392-1548-4e13-b8d9-36f0159a10ab"),
				b: uuid.NullUUID{}.UUID,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameUUID(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("IsSameUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
