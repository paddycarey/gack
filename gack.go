// Package gack is a library designed to allow easy implementation of slash
// commands using the Slack API.
//
// gack provides a Server struct which implements the http.Handler interface
// (so it can be used almost anywhere). The user of the library must provide
// their own Handler implementations to control routing of and processing of
// incoming commands.
package gack

import (
	"fmt"
	"net/http"
)

// Handler is an interface that should be implemented by any structure designed
// to process and respond to a slash command from Slack.
type Handler interface {
	// CanHandle checks the SlashCommand that has been passed in to determine
	// if this Handler can (or should) actually process it. Each Handler
	// implementation should choose how best to decide if it's able to handle
	// the command. Some Handlers might simply check the Command or the TeamID
	// that were passed match some predefined values, whilst others might
	// return true unconditionally (e.g. a fall-through handler that fires if
	// no others have matched).
	CanHandle(*SlashCommand) bool

	// Handle is called when processing a SlashCommand. Normally this method is
	// called by a Server instance after checking that CanHandle returns true.
	// Implementations should use the Handle method to process an incoming
	// SlashCommand and generate a response. If an error is returned, it is
	// stringified and returned to the user who sent the command. If a string
	// is returned, it is sent to the originating user verbatim. Handle is not
	// required to return a string or an error (it can return "", nil), in this
	// case the implementation may want to send a notification to the user
	// using a different mechanism e.g. the Slack web API or a Bot user.
	Handle(*SlashCommand) (string, error)
}

// SlashCommand holds all details of the command entered by a user within
// Slack. When the user types a "slash command", Slack will send an HTTP POST
// request containing the entered text plus additional metadata. The body of
// this HTTP request is unmarshalled into a SlashCommand struct.
type SlashCommand struct {
	ChannelID   string
	ChannelName string
	Command     string
	TeamDomain  string
	TeamID      string
	Text        string
	Token       string
	UserID      string
	UserName    string
}

// Server is a struct which implements the http.Handler interface. It is
// designed to act as the main API endpoint for all slash commands sent from
// Slack.
type Server struct {
	tokens   map[string]bool
	handlers []Handler
}

// NewServer creates a new Server instance to handle incoming slash commands.
//
// The first argument is a slice containing any API tokens that should be used
// to authenticate incoming commands. If you don't pass at least one token, you
// will not be able to accept any incoming commands.
//
// The second argument is a slice containing all Handlers that should be used
// to handle incoming commands. Each handler will be tried in sequence, when a
// matching handler for a given command is found the command will be processed
// and a response will (most of the time) be returned.
func NewServer(tokens []string, handlers []Handler) *Server {
	t := make(map[string]bool)
	for _, v := range tokens {
		t[v] = true
	}
	return &Server{
		tokens:   t,
		handlers: handlers,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	sc, err := s.parseSlashCommand(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	// authenticate the incoming request by checking that the provided token is
	// already in our s.tokens map. If the token isn't recognised then the
	// request will be rejected before any processing takes place.
	if !s.tokens[sc.Token] {
		fmt.Fprint(w, "Invalid API token. Check configuration.")
		return
	}

	// loop over all command handlers until we find one that can handle the
	// given command. Once we find a matching command handler it is run
	// immediately and any returned text is passed verbatim back to the client
	// that initiated the request.
	var handleErr error
	var handleOut string
	for _, h := range s.handlers {
		if !h.CanHandle(sc) {
			continue
		}
		handleOut, handleErr = h.Handle(sc)
		break
	}

	if handleErr != nil {
		fmt.Fprint(w, handleErr.Error())
	} else {
		fmt.Fprint(w, handleOut)
	}
}

// parseSlashCommand extracts form data (whether from the body or query string)
// from an incoming HTTP request and unmarshalls it into a SlashCommand struct.
func (s *Server) parseSlashCommand(r *http.Request) (*SlashCommand, error) {

	sc := &SlashCommand{}

	// parse data from the HTTP request
	err := r.ParseForm()
	if err != nil {
		return sc, err
	}

	// load form data into a SlashCommand struct
	sc.ChannelID = r.Form.Get("channel_id")
	sc.ChannelName = r.Form.Get("channel_name")
	sc.Command = r.Form.Get("command")
	sc.TeamDomain = r.Form.Get("team_domain")
	sc.TeamID = r.Form.Get("team_id")
	sc.Text = r.Form.Get("text")
	sc.Token = r.Form.Get("token")
	sc.UserID = r.Form.Get("user_id")
	sc.UserName = r.Form.Get("user_name")

	return sc, nil
}
