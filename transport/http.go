package transport

import (
	"context"
	"encoding/json"
	"errors"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"ther.cool/micro-go-course-user/endpoint"
)

var ErrorBadRequest = errors.New("invalid request parameter")

// make http handler use mux
func NewHttpHandler(_ context.Context, endpoints *endpoint.UserEndpoints) http.Handler {
	router := mux.NewRouter()

	kitLog := kitlog.NewLogfmtLogger(os.Stderr)
	kitLog = kitlog.With(kitLog, "ts", kitlog.DefaultTimestampUTC)
	kitLog = kitlog.With(kitLog, "caller", kitlog.DefaultCaller)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kithttp.ServerErrorEncoder(encodeErr),
	}

	router.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		options...,
	))

	router.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		options...,
	))

	return router
}

func decodeRegisterRequest(_ context.Context, request *http.Request) (interface{}, error) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	email := request.FormValue("email")
	if username == "" || password == "" || email == "" {
		return nil, ErrorBadRequest
	}

	return &endpoint.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}, nil
}

func decodeLoginRequest(_ context.Context, request *http.Request) (interface{}, error) {
	email := request.FormValue("email")
	password := request.FormValue("password")
	if email == "" || password == "" {
		return nil, ErrorBadRequest
	}

	return &endpoint.LoginRequest{
		Email:    email,
		Password: password,
	}, nil
}

func encodeJSONResponse(_ context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

func encodeErr(_ context.Context, err error, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch err {
	default:
		writer.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
