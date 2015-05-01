package main

import (
	"net/http"
	"os"
	"time"

	"github.com/paddycarey/gack"
)

// ClockHandler returns the current time in a given IANA timezone.
type ClockHandler struct{}

func (h *ClockHandler) CanHandle(sc *gack.SlashCommand) bool {
	return true
}

// Handle returns the current time for the provided timezone.
func (h *ClockHandler) Handle(sc *gack.SlashCommand) (string, error) {
	zone, err := time.LoadLocation(sc.Text)
	if err != nil {
		return "", err
	}
	t := time.Now().In(zone)
	return t.Format(time.UnixDate), nil
}

func main() {
	mux := http.NewServeMux()

	srv := gack.NewServer(
		[]string{os.Getenv("SLACK_API_TOKEN")},
		[]gack.Handler{&ClockHandler{}},
	)

	mux.Handle("/", srv)
	http.ListenAndServe(":3000", mux)
}
