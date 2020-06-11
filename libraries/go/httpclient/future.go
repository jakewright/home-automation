package httpclient

// Future represents an in-flight request
type Future struct {
	done     <-chan struct{}
	request  *Request
	response *Response
	err      error
}

// Response blocks until the response is available
func (f *Future) Response() (*Response, error) {
	<-f.done
	return f.response, f.err
}

// DecodeResponse blocks until the response is available. If there was an error
// during the request, it is returned immediately. Otherwise, the response is
// decoded into v. If there is an error while decoding, both the response and
// the error are returned. Decoding errors will be of type DecodeError.
func (f *Future) DecodeResponse(v interface{}) (*Response, error) {
	<-f.done

	if f.err != nil {
		return f.response, f.err
	}

	if v != nil {
		dec, err := decoder(f.request, f.response)
		if err != nil {
			return f.response, err
		}

		body, err := f.response.BodyBytes()
		if err != nil {
			return f.response, err
		}

		if err := dec.Decode(body, v); err != nil {
			return f.response, &DecodeError{
				Format: dec.Name(),
				Err:    err,
			}
		}
	}

	return f.response, nil
}

func decoder(req *Request, rsp *Response) (Decoder, error) {
	if req.ResponseDecoder != nil {
		return req.ResponseDecoder, nil
	}

	return inferDecoder(rsp.Header.Get("Content-Type"))
}
