package main

import (
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAuthentication(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/")

		validCSRFToken := extractCSRFToken(t, body)

		code, header, _ := ts.postHX(t, "/stop", validCSRFToken)

		assert.Equal(t, code, http.StatusOK)
		assert.Equal(t, header.Get("HX-Redirect"), "/user/login")
	})
}
