package dino

// Middleware is something
type Middleware func(Handler) Handler

func applyMiddlewares(h Handler, middlewares ...Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
