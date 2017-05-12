// Package router combines routing and middleware handling in a single
// package.
package router

import (
	"net/http"
	"sync"

	"github.com/husobee/vestigo"
)

// Item implements the password hashing system.
type Item struct {
	r         *vestigo.Router
	routeList []string
	listMutex sync.RWMutex
}

// New returns a new instance of the router.
func New() *Item {
	s := new(Item)
	s.r = vestigo.NewRouter()
	return s
}

// Router returns the router.
func (s *Item) Router() http.Handler {
	return s.r
}

// SetNotFound sets the 404 handler.
func (s *Item) SetNotFound(fn http.HandlerFunc) {
	vestigo.CustomNotFoundHandlerFunc(fn)
}

// SetMethodNotAllowed sets the 405 handler.
func (s *Item) SetMethodNotAllowed(fn vestigo.MethodNotAllowedHandlerFunc) {
	vestigo.CustomMethodNotAllowedHandlerFunc(fn)
}

// Param returns the URL parameter.
func (s *Item) Param(r *http.Request, name string) string {
	return vestigo.Param(r, name)
}
