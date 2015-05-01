gack
====

[![Build Status](https://travis-ci.org/paddycarey/gack.svg?branch=master)](https://travis-ci.org/paddycarey/gack)
[![GoDoc](https://godoc.org/github.com/paddycarey/gack?status.svg)](https://godoc.org/github.com/paddycarey/gack)

gack is a library designed to allow easy implementation of APIs that respond to [slash commands](https://api.slack.com/slash-commands) from the [Slack API](https://api.slack.com). gack provides the ability to define multiple handlers per endpoint/server instance (with a straightforward route-like behaviour), it has built in authentication using API tokens and knows how to unmarshall an incoming slash command.

gack provides a `Server` struct which implements the `http.Handler` interface (so it can be used almost anywhere). The user of the library must provide their own Handler implementations to control routing of and processing of incoming commands.

Several small example applications with simple `Handler` implementations are included within the [examples](examples) directory.

```go
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

func (h *EchoHandler) CanHandle(sc *gack.SlashCommand) bool {
	return true
}

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
```

See [GoDoc](https://godoc.org/github.com/paddycarey/gack) for full library documentation.


### Copyright & License

- Copyright © 2015 Patrick Carey (https://github.com/paddycarey)
- Licensed under the **MIT** license.
