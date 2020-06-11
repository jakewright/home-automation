package httpclient

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToQueryString(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "nil",
			v:    nil,
			want: "",
		},
		{
			name: "url.Values{}",
			v:    url.Values{"foo": {"bar", "bat"}},
			want: "foo=bar&foo=bat",
		},
		{
			name: "map[string][]string",
			v: map[string][]string{
				"foo": {"bar", "bat"},
			},
			want: "foo=bar&foo=bat",
		},
		{
			name: "map[string]string",
			v: map[string]string{
				"foo": "bar",
				"bat": "baz",
			},
			want: "bat=baz&foo=bar",
		},
		{
			name: "empty struct",
			v:    struct{}{},
			want: "",
		},
		{
			name: "struct",
			v: struct {
				Foo string `form:"foo"`
				Bar int    `form:"bar"`
				Bat bool   `form:"bat"`
			}{"foo", 5, true},
			want: "bar=5&bat=true&foo=foo",
		},
		{
			name: "struct with slice",
			v: struct {
				Foo []string `form:"foo"`
			}{
				Foo: []string{"bar", "bat"},
			},
			want: "foo=bar&foo=bat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toQueryString(tt.v, "form")
			if (err != nil) != tt.wantErr {
				t.Errorf("ToURLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
