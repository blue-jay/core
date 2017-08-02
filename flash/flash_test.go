package flash_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/view"
)

// Session is an interface for typical sessions
type Session struct {
	flashes []interface{}
	mutex   sync.RWMutex
}

// Save mocks saving the session
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return nil
}

// Flashes retrieves the flashes
func (s *Session) Flashes(vars ...string) []interface{} {
	s.mutex.RLock()
	f := s.flashes
	// Clear the flashes
	s.flashes = s.flashes[:0]
	s.mutex.RUnlock()
	return f
}

// AddFlash adds a flash to the list
func (s *Session) AddFlash(f interface{}) {
	s.mutex.Lock()
	s.flashes = append(s.flashes, f)
	s.mutex.Unlock()
}

// TestFlashSession ensures flashes can be added to the session.
func TestFlashSession(t *testing.T) {
	text := "Success test."

	// Get the fake session
	sess := Session{}

	// Add flashes to the session
	sess.AddFlash(flash.Success(text))

	// Get the flashes
	flashes := sess.Flashes()

	if len(flashes) != 1 {
		t.Fatal("Expected 1 flash message.")
	}

	// Convert the flash
	f, ok := flashes[0].(flash.Message)

	if f.Class != "alert-success" {
		t.Fatalf("Flash class is: %v, should be: %v.", f.Class, "alert-success")
	}

	if f.Content != text {
		t.Fatalf("Flash message is: %v, should be: %v", f.Content, text)
	}

	if !ok {
		t.Fatal("Flashes missing from session.")
	}
}

// TestSendFlashes are available for AJAX.
func TestSendFlashes(t *testing.T) {
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

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Success test."

	// Get the fake session
	sess := &Session{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add flashes to the session
		sess.AddFlash(flash.Success(text))
		sess.AddFlash(text)

		// Send the flashes
		flash.SendFlashes(w, r, sess)
	})

	handler.ServeHTTP(w, r)

	actual := w.Body.String()
	expected := fmt.Sprintf(`[{"Content":"%v","Class":"%v"},{"Content":"%v","Class":"%v"}]`, text, "alert-success", text, "alert-box")

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}
