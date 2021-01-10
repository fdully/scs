package support

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type okSupport struct{}

func (f okSupport) Send(ctx context.Context, msg string) error {
	return nil
}

type badSupport struct{}

func (f badSupport) Send(ctx context.Context, msg string) error {
	return errors.New("test error")
}

const testUrl = "http://example.com/api/v1/send/support"

func TestSupport(t *testing.T) {

	s := okSupport{}
	ctx := context.Background()

	h := Handle(ctx, s)

	t.Run("bad request", func(t *testing.T) {
		req := httptest.NewRequest("GET", testUrl, nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		resp := w.Result()

		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, "", string(body))
	})

	t.Run("ok", func(t *testing.T) {
		urlData := url.Values{}
		urlData.Set("message", "test")

		req := httptest.NewRequest("GET", testUrl+"?"+urlData.Encode(), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		resp := w.Result()

		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "", string(body))
	})

	t.Run("internal server error", func(t *testing.T) {
		s := badSupport{}
		h := Handle(ctx, s)

		urlData := url.Values{}
		urlData.Set("message", "test")

		req := httptest.NewRequest("GET", testUrl+"?"+urlData.Encode(), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		resp := w.Result()

		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		require.Equal(t, "", string(body))
	})

}
