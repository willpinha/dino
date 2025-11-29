package main

import (
	"fmt"
	"net/http"

	"github.com/willpinha/httpbox"
)

func HelloHandler() httpbox.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		name := httpbox.NewPathParam(r, "name").String()

		if name == "bob" {
			return httpbox.NewError(
				http.StatusBadRequest,
				"Bob is banned",
				httpbox.WithDetails("Please contact the HR"),
			)
		}

		msg := fmt.Sprintf("Hello, %s!", name)

		return httpbox.WriteJSON(w, 200, msg)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /hello/{name}", HelloHandler())

	h := httpbox.AdaptHandler(mux).WithMiddlewares(
		httpbox.AccessLogMiddleware(),
	)

	http.ListenAndServe(":8080", h)
}
