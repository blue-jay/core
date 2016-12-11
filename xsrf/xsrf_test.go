package xsrf_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

// TestModify ensures flashes are added to the view.
func TestModify(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "test",
		Children: []string{},
	}

	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Apply the flash modifier
	viewInfo.SetModifiers(
		xsrf.Token,
	)

	// Set up the session cookie store
	s.SetupConfig()

	// Decode the string
	key, err := base64.StdEncoding.DecodeString(s.AuthKey)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := viewInfo.New()

		err := v.Render(w, r)
		if err != nil {
			t.Fatalf("Should not get error: %v", err)
		}
	})

	// Configure the middleware
	cs := csrf.Protect([]byte(key),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("invalidHandler should not be called.")
		})),
		csrf.FieldName("_token"),
		csrf.Secure(s.Options.Secure),
	)(handler)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Get the session
	sess, err := s.Instance(r)
	if err != nil {
		t.Fatal(err)
	}

	cs.ServeHTTP(w, r)

	log.Println(r)

	actual := w.Body.String()
	expected := fmt.Sprintf(`<div>%v</div>`, sess.Values["_token"])

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}
