package main

import (
	"net/http"
	"net/url"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAuthentication(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/")

		validCSRFToken := extractCSRFTokenHx(t, body)

		code, header, _ := ts.postHX(t, "/stop", validCSRFToken)

		assert.Equal(t, code, http.StatusOK)
		assert.Equal(t, header.Get("HX-Redirect"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		validCSRFToken := extractCSRFTokenForm(t, body)

		form := url.Values{}
		form.Add("username", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", validCSRFToken)

		code, header, _ := ts.postForm(t, "/user/login", form)
		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/")

		code, header, _ = ts.postHX(t, "/stop", validCSRFToken)
		assert.Equal(t, code, http.StatusOK)
		assert.Equal(t, header.Get("Hx-Redirect"), "")
	})

}
