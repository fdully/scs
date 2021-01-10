package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/fdully/scs/internal/logging"
	"github.com/fdully/scs/internal/scs/support"
	"github.com/fdully/scs/internal/server"
	"github.com/gorilla/handlers"
	"github.com/sethvargo/go-signalcontext"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		return errors.New("PORT env is not set")
	}

	mux := http.NewServeMux()
	mux.Handle("/health", server.HandleHealthz(ctx))
	mux.Handle("/api/v1/send/support", support.Handle(ctx, support.NewSupport()))

	logger.Infof("serving on port: %s", port)
	srv, err := server.New(port)
	if err != nil {
		return err
	}

	return srv.ServeHTTPHandler(ctx, handlers.CombinedLoggingHandler(os.Stdout, mux))
}
