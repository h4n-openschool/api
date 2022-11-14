package server

import "net/http"

func validRequest(req *http.Request) *http.Response {
	res := NewResponse()

	if ok := validMethod(req.Method); !ok {
		res.StatusCode = 405
		return &res
	}

	if ok := req.ProtoAtLeast(1, 1); !ok {
		res.StatusCode = 400
		return &res
	}

	return nil
}

func validMethod(method string) bool {
	isValid := (method == http.MethodGet ||
		method == http.MethodPost ||
		method == http.MethodPatch ||
		method == http.MethodPut ||
		method == http.MethodDelete ||
		method == http.MethodOptions ||
		method == http.MethodConnect ||
		method == http.MethodTrace ||
		method == http.MethodHead)
	return isValid
}
