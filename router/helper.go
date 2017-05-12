package router

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

// Chain returns an array of middleware.
func Chain(c ...alice.Constructor) []alice.Constructor {
	return c
}

// ChainHandler returns a handler of chained middleware.
func ChainHandler(h http.Handler, c ...alice.Constructor) http.Handler {
	return alice.New(c...).Then(h)
}

// record stores the method and path.
func (s *Item) record(method, path string) {
	s.listMutex.Lock()
	s.routeList = append(s.routeList, fmt.Sprintf("%v\t%v", method, path))
	s.listMutex.Unlock()
}

// RouteList returns a list of the HTTP methods and paths.
func (s *Item) RouteList() []string {
	s.listMutex.RLock()
	list := s.routeList
	s.listMutex.RUnlock()
	return list
}

// Delete is a shortcut for router.Handle("DELETE", path, handle).
func (s *Item) Delete(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	s.record("DELETE", path)
	s.r.Delete(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
}

// Get is a shortcut for router.Handle("GET", path, handle).
func (s *Item) Get(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	s.record("GET", path)
	s.r.Get(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
}

// Patch is a shortcut for router.Handle("PATCH", path, handle).
func (s *Item) Patch(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	s.record("PATCH", path)
	s.r.Patch(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
}

// Post is a shortcut for router.Handle("POST", path, handle).
func (s *Item) Post(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	s.record("POST", path)
	s.r.Post(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
}

// Put is a shortcut for router.Handle("PUT", path, handle).
func (s *Item) Put(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	s.record("PUT", path)
	s.r.Put(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
}
