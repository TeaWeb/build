package requests

import "net/http"

type Response struct {
	*http.Response

	BodyData []byte
}

func NewResponse(resp *http.Response) *Response {
	return &Response{
		Response: resp,
	}
}
