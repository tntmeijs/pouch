package pouch

import (
	"context"
	"net/http"
)

type stubContextKey string
type StubContextStorage struct {
	Data map[string]any
}

const contextStubsKey stubContextKey = "stubs"

// Custom interceptor object, used as an http.Transport object in requests.
type StubInterceptorRoundTrip struct {
	originalTransport http.RoundTripper
}

// ConfigureTransportForStubbing wraps an http.RoundTripper in a custom transport object that supports stubbing.
// Use this transport in requests to ensure that your requests are intercepted and, potentially, stubbed.
// This object effectively acts as "middleware", if you will. For those familiar with popular web frameworks.
func ConfigureTransportForStubbing(transport http.RoundTripper) StubInterceptorRoundTrip {
	return StubInterceptorRoundTrip{originalTransport: transport}
}

// NewStubbedContext creates a new, empty, context with a special key that marks the context as eligible for stubbing.
func NewStubbedContext() context.Context {
	return EnableStubsForExistingContext(context.Background())
}

// EnableStubsForExistingContext appends a special key to the existing context that marks it as eligible for stubbing.
func EnableStubsForExistingContext(c context.Context) context.Context {
	return EnableStubsForExistingContextWithData(c, &StubContextStorage{})
}

// EnableStubsForExistingContextWithData appends a special key to the existing context that marks it as eligible for stubbing.
// It will also store an object in the context that can later be referenced in stub templates.
// This is particularly useful when you want to insert dynamic application data into a stub, such as a value from a database.
func EnableStubsForExistingContextWithData(c context.Context, storage *StubContextStorage) context.Context {
	return context.WithValue(c, contextStubsKey, storage)
}

// RoundTrip implements the http.RoundTripper interface and extends it with request interception capabilities.
func (s *StubInterceptorRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) {
	storage, ok := request.Context().Value(contextStubsKey).(*StubContextStorage)

	if ok {
		// Intercept the request and return a stubbed result instead.
		// The original request will never hit the server it was intended for.
		return handleStubRequest(request, storage)
	}

	// Allow the request to hit the original target instead of the stubs.
	return s.originalTransport.RoundTrip(request)
}

// The handleStubRequest method generates a stubbed response for the given request.
func handleStubRequest(request *http.Request, storage *StubContextStorage) (*http.Response, error) {
	// TODO: implement stubbing logic
	return &http.Response{
		Status: "I'm a teapot",
		StatusCode: 418,
	}, nil
}
