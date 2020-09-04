package transport

import (
	"context"
	"errors"
)

var ErrorBadRequest = errors.New("invalid request parameter")

// make http handler use mux
func MakeHttpHandler(ctx context.Context) {

}
