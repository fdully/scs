package support

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

type Sender interface {
	Send(ctx context.Context, msg string) error
}

func Handle(ctx context.Context, sender Sender) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg string

		switch r.Method {
		case "GET":
			msg = r.URL.Query().Get("message")
		case "POST":
			msg = r.PostFormValue("message")
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if msg == "" {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if err := sender.Send(r.Context(), msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func NewSupport() *Support {
	return &Support{
		slackURL:  os.Getenv("SLACK_SUPPORT_URL"),
		maxMsgLen: 2000,
	}
}

type Support struct {
	slackURL  string
	maxMsgLen int
}

func (s Support) Send(ctx context.Context, msg string) error {
	if len(msg) > s.maxMsgLen {
		msg = msg[0:s.maxMsgLen]
	}

	message := struct {
		Text string `json:"text"`
	}{
		Text: msg,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(s.slackURL, "application/json", bytes.NewReader(b))

	return err
}
