package httpclient

import (
	"net/url"

	"github.com/gorilla/schema"
)

func toQueryString(v interface{}, tagAlias string) (string, error) {
	if v == nil {
		return "", nil
	}

	if m, ok := v.(url.Values); ok {
		return m.Encode(), nil
	}

	if m, ok := v.(map[string][]string); ok {
		return url.Values(m).Encode(), nil
	}

	if m, ok := v.(map[string]string); ok {
		values := url.Values{}
		for key, value := range m {
			values.Set(key, value)
		}
		return values.Encode(), nil
	}

	// Note: gorilla/schema only supports structs

	e := schema.NewEncoder()
	e.SetAliasTag(tagAlias)

	values := url.Values{}
	if err := e.Encode(v, values); err != nil {
		return "", err
	}

	return values.Encode(), nil
}
