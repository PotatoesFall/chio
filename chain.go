package chi

import "net/http"

// Chain returns a Middlewares type from a slice of middleware handlers.
func Chain(middlewares ...func(http.HandlerFunc) http.HandlerFunc) Middlewares {
	return Middlewares(middlewares)
}

// Handler builds and returns a http.Handler from the chain of middlewares,
// with `h http.Handler` as the final handler.
func (mws Middlewares) Handler(h http.Handler) http.Handler {
	return &ChainHandler{h.ServeHTTP, chain(mws, h.ServeHTTP), mws}
}

// HandlerFunc builds and returns a http.Handler from the chain of middlewares,
// with `h http.Handler` as the final handler.
func (mws Middlewares) HandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

// ChainHandler is a http.Handler with support for handler composition and
// execution.
type ChainHandler struct {
	Endpoint    http.HandlerFunc
	chain       http.HandlerFunc
	Middlewares Middlewares
}

func (c *ChainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.chain.ServeHTTP(w, r)
}

// chain builds a http.Handler composed of an inline middleware stack and endpoint
// handler in the order they are passed.
func chain(middlewares []func(http.HandlerFunc) http.HandlerFunc, endpoint http.HandlerFunc) http.HandlerFunc {
	// Return ahead of time if there aren't any middlewares for the chain
	if len(middlewares) == 0 {
		return endpoint
	}

	// Wrap the end handler with the middleware chain
	h := middlewares[len(middlewares)-1](endpoint)
	for i := len(middlewares) - 2; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
