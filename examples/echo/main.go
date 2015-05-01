package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/paddycarey/gack"
)

// EchoHandler is the simplest possible implementation of the gack.Handler
// interface, it unconditionally echoes any commands a user types directly back
// at them.
type EchoHandler struct{}

// CanHandle in this case always returns true. This causes this handler to
// unconditionally respond to any command that it is passed.
func (h *EchoHandler) CanHandle(sc *gack.SlashCommand) bool {
	return true
}

// Handle concatenates the command and text fields from the incoming
// SlashCommand to reform the text as the user originally typed it within
// Slack.
func (h *EchoHandler) Handle(sc *gack.SlashCommand) (string, error) {
	return fmt.Sprintf("%s %s", sc.Command, sc.Text), nil
}

func main() {
	mux := http.NewServeMux()

	srv := gack.NewServer(
		[]string{os.Getenv("SLACK_API_TOKEN")},
		[]gack.Handler{&EchoHandler{}},
	)

	mux.Handle("/", srv)
	http.ListenAndServe(":3000", mux)
}
