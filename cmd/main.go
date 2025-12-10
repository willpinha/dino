package main

import (
	"fmt"
	"net/http"

	"github.com/willpinha/dino"
)

func HelloHandler() dino.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		name := dino.NewPathParam(r, "name").String()

		if name == "bob" {
			return dino.NewError(
				http.StatusBadRequest,
				"Bob is banned",
				dino.WithDetails("Please contact the HR"),
			)
		}

		msg := fmt.Sprintf("Hello, %s!", name)

		return dino.WriteJSON(w, 200, msg)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /hello/{name}", HelloHandler())

	h := dino.AdaptHandler(mux).WithMiddlewares(
		dino.AccessLogMiddleware(),
	)

	http.ListenAndServe(":8080", h)
}
