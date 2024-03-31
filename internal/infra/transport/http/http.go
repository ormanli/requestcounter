package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ormanli/requestcounter/internal/app/requestcounter"
)

// Service is business logic implementation.
type Service interface {
	Increment(ctx context.Context) (requestcounter.IncrementResult, error)
}

// Transport is the HTTP component.
type Transport struct {
	service Service
	srv     *http.Server
	cfg     requestcounter.Config
}

// NewTransport initializes a new Transport.
func NewTransport(cfg requestcounter.Config, service Service, mw ...func(http.Handler) http.Handler) *Transport {
	a := &Transport{
		service: service,
		cfg:     cfg,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", a.count)

	a.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:           chainMiddlewares(mw...)(mux),
		ReadHeaderTimeout: cfg.ServerReadHeaderTimeout,
	}

	return a
}

const countResponseText = `You are talking to instance %s:%d.
This is request %d to this instance and request %d to the cluster.`

func (g *Transport) count(w http.ResponseWriter, r *http.Request) {
	result, err := g.service.Increment(r.Context())
	if err != nil {
		writeTextPlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeTextPlainResponse(w, http.StatusOK, fmt.Sprintf(countResponseText, g.cfg.GetServerHost(), g.cfg.ServerPort, result.LocalCount, result.ClusterCount))
}

func writeTextPlainResponse(w http.ResponseWriter, statusCode int, resp string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, resp)
}

// Run starts HTTP server.
func (g *Transport) Run() error {
	slog.Info("Server started", "port", g.cfg.ServerPort)

	if err := g.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Stop stops HTTP server.
func (g *Transport) Stop(ctx context.Context) error {
	slog.Info("Server stopped")

	return g.srv.Shutdown(ctx)
}

func chainMiddlewares(mw ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last.ServeHTTP(w, r)
		})
	}
}
