package server

import (
	"context"
	"fmt"
	"net/http"
)

func HandleHealthz(hctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"status": "ok"}`)
	})
}
