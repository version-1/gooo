package response

import "net/http"

func JSON[O any]() *Response[O] {
	return &Response[O]{
		adapter: JSONAdapter{},
		status:  http.StatusOK,
	}
}

func HTML[O any]() *Response[O] {
	return &Response[O]{
		adapter: HTMLAdapter{},
		status:  http.StatusOK,
	}
}
