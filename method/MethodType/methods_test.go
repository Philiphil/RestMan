package method_type_test

import (
	"testing"

	. "github.com/philiphil/apiman/method/MethodType"
)

func TestApiMethod_String(t *testing.T) {
	tests := []struct {
		name string
		e    ApiMethod
		want string
	}{
		{
			name: "Patch",
			e:    Patch,
			want: "PATCH",
		},
		{
			name: "Post",
			e:    Post,
			want: "POST",
		},
		{
			name: "Put",
			e:    Put,
			want: "PUT",
		},
		{
			name: "Get",
			e:    Get,
			want: "GET",
		},
		{
			name: "Head",
			e:    Head,
			want: "HEAD",
		},
		{
			name: "Delete",
			e:    Delete,
			want: "DELETE",
		},
		{
			name: "Options",
			e:    Options,
			want: "OPTIONS",
		},
		{
			name: "Trace",
			e:    Trace,
			want: "TRACE",
		},
		{
			name: "Connect",
			e:    Connect,
			want: "CONNECT",
		},
		{
			name: "Undefined",
			e:    Undefined,
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ApiMethod.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
