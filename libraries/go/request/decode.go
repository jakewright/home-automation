package request

import (
	"encoding/json"
	"home-automation/libraries/go/errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// Decode unmarshals URL parameters and the JSON body of the given request into the output interface
func Decode(r *http.Request, v interface{}) error {
	// This does a load of reflection to unmarshal URL params into body
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true,
		Result:           v,
	})
	if err != nil {
		return errors.Wrap(err, nil)
	}
	if err := decoder.Decode(mux.Vars(r)); err != nil {
		return errors.Wrap(err, nil)
	}

	// Read the body of the request
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, nil)
	}

	// Assume the data is JSON and unmarshal into body
	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, nil)
	}

	return nil
}
