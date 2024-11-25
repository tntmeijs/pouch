package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tntmeijs/pouch"
)

func main() {
	// You can use any transport configuration you would like, for example if you use TLS, you would use your own transport here.
	stubTransport := pouch.ConfigureTransportForStubbing(http.DefaultTransport)

	client := http.Client{
		Transport: &stubTransport,
	}

	// To use stubs, the context must be configured.
	// You can use any of the other exported context-related functions if you would like to use an existing context instead.
	ctx := pouch.NewStubbedContext()

	// Make the request and ensure that the context is passed along with it - this ensures that the interceptor triggers.
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost/pouch/00/basic", nil)

	// Any other logic from this point onwards does not need to know about the stubs.
	if err != nil {
		log.Fatalf("Unable to create request: %s", err.Error())
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatalf("Unable to send request: %s", err.Error())
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("Unable to read response: %s", err.Error())
	}

	fmt.Printf("HTTP %d: %s\n", res.StatusCode, res.Status)
	fmt.Println(string(bytes))
}
