package gack

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type unconditionalHandler struct {
	canHandle bool
	handle    string
}

func (u *unconditionalHandler) CanHandle(*SlashCommand) bool {
	return u.canHandle
}

func (u *unconditionalHandler) Handle(*SlashCommand) (string, error) {
	return u.handle, nil
}

func testBody(overrides url.Values) *bytes.Buffer {
	vals := &url.Values{}
	vals.Set("channel_id", "C2147483705")
	vals.Set("channel_name", "test")
	vals.Set("command", "/weather")
	vals.Set("team_domain", "example")
	vals.Set("team_id", "T0001")
	vals.Set("text", "94070")
	vals.Set("token", "gIkuvaNzQIHg97ATvDxqgjtO")
	vals.Set("user_id", "U2147483697")
	vals.Set("user_name", "Steve")

	if overrides != nil {
		for k, _ := range overrides {
			vals.Set(k, overrides.Get(k))
		}
	}

	return bytes.NewBuffer([]byte(vals.Encode()))
}

func testRequest(overrides url.Values, srv *Server) *httptest.ResponseRecorder {
	r, err := http.NewRequest("POST", "/", testBody(overrides))
	r.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded; param=value",
	)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	return w
}

func TestInvalidAPITokens(t *testing.T) {
	tokens := []string{
		"asdadwqd",
		"dasdsdc",
		"2312312w",
		"21wj12o8",
		"",
		" ",
		"0okm9iujh8iuyhgv nhygv bhgbv nb hb \n\noudo8d",
	}

	for _, token := range tokens {
		vals := url.Values{}
		vals.Set("token", token)
		srv := NewServer([]string{"aaa"}, []Handler{})
		w := testRequest(vals, srv)

		respBody := string(w.Body.Bytes())
		if respBody != "Invalid API token. Check configuration." {
			t.Errorf("Unexpected response: %s", respBody)
			break
		}
	}
}

func TestValidAPITokens(t *testing.T) {
	tokens := []string{
		"aaa",
		"dasdsdc",
		"",
	}

	for _, token := range tokens {
		vals := url.Values{}
		vals.Set("token", token)
		srv := NewServer([]string{"aaa", "dasdsdc", ""}, []Handler{})
		w := testRequest(vals, srv)

		respBody := string(w.Body.Bytes())
		if respBody != "" {
			t.Errorf("Unexpected response: \"%s\"", respBody)
			break
		}
	}
}
