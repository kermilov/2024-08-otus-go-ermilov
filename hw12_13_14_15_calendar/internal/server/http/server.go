package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Server struct {
	http.Server
	logger Logger
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface { // TODO
}

type myHandler struct{}

// реализуем интерфейс `http.Handler`.
func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/hello" {
		args := r.URL.Query()
		name := args.Get("name")
		results, err := h.hello(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	}
}

func (h *myHandler) hello(name string) (string, error) {
	if name != "" {
		if name == "world" {
			return "", errors.New("error")
		}
		return `{"hello":"` + name + `"}`, nil
	}
	return `{"hello":"world"}`, nil
}

func NewServer(logger Logger, _ Application) *Server {
	handler := &myHandler{}
	return &Server{
		Server: http.Server{
			Addr:         ":8080",
			Handler:      loggingMiddleware(handler),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

// TODO
