# Handlers & Routing

## Handlers

Handlers (a.k.a. Controllers) represents the core mechanism for processing HTTP requests and
returning HTTP responses to the clients

Similar to the [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) type of `net/http`,
a handler in Dino is a function, with the difference that this function returns an error. See below
for an example of a handler:

```go
func MyHandler() dino.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Handler logic here
		return nil
	}
}
```

Defining a custom handler that returns an error is a very common pattern because it simplifies error
handling, as we will see in the following sections. This pattern is used, for example, on
[Fiber](https://gofiber.io) and [Echo](https://echo.labstack.com) handlers

## Routing

Dino's handlers implement the [`http.Handler`](https://pkg.go.dev/net/http#Handler) interface, and
therefore can be used anywhere that interface is used. An example of this is in
[`http.ServeMux`](https://pkg.go.dev/net/http#ServeMux):

```go
mux := http.NewServeMux()

mux.Handle("GET /my/path", MyHandler())
mux.Handle("POST /another/path", AnotherHandler())
```

In fact, `http.ServeMux` is the ideal multiplexer (router) to use with Dino. This is because many
of Dino's features are built upon the functionalities that `http.ServeMux` provides. Besides, it's
part of the standard library and therefore a very stable router
