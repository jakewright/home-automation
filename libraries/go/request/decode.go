package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jakewright/home-automation/libraries/go/errors"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// Decode unmarshals URL parameters and the JSON body of the given request into the output interface.
// Parameters can be unmarshalled into primitive types or time.Time providing it conforms to time.RFC3339.
func Decode(r *http.Request, v interface{}) error {
	// This does a load of reflection to unmarshal a map into the type of v
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true,
		Result:           v,

		// Override the TagName to match the one used by the encoding/json package
		// so users of this function only have to define a single tag on struct fields
		TagName: "json",
	})
	if err != nil {
		return errors.Wrap(err, nil)
	}

	// Unmarshal route parameters
	if err := decoder.Decode(mux.Vars(r)); err != nil {
		return errors.Wrap(err, nil)
	}

	// Query parameters come out as a map[string][]string so we loop through them all
	// to remove the unnecessary slice if the parameter just has a single value
	paramSlices := r.URL.Query()
	params := map[string]interface{}{}
	for key, value := range paramSlices {
		switch len(value) {
		case 0:
			params[key] = nil
		case 1:
			params[key] = value[0]
		default:
			params[key] = value
		}
	}

	// Unmarshal query parameters
	if err := decoder.Decode(params); err != nil {
		return errors.Wrap(err, nil)
	}

	// If there's no body, return early
	if r.Body == nil {
		return nil
	}

	// Read the body of the request
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, nil)
	}

	// If the body is empty, return early
	if len(body) == 0 {
		return nil
	}

	// Assume the body is JSON and unmarshal into v
	if err := json.Unmarshal(body, v); err != nil {
		return errors.Wrap(err, nil)
	}

	return nil
}
