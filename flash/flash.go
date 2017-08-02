// Package flash provides one-time messages for the user.
package flash

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	// danger is a bootstrap class.
	danger = "alert-danger"
	// success is a bootstrap class.
	success = "alert-success"
	// notice is a bootstrap class.
	notice = "alert-info"
	// warning is a bootstrap class.
	warning = "alert-warning"
	// standard is the default class.
	standard = "alert-box"
)

// Message is the flash information.
type Message struct {
	Content string
	Class   string
}

// Success returns a success flash message.
func Success(message string) Message {
	return Message{
		Content: message,
		Class:   success,
	}
}

// Danger returns a danger flash message.
func Danger(message string) Message {
	return Message{
		Content: message,
		Class:   danger,
	}
}

// Notice returns a notice flash message.
func Notice(message string) Message {
	return Message{
		Content: message,
		Class:   notice,
	}
}

// Warning returns a warning flash message.
func Warning(message string) Message {
	return Message{
		Content: message,
		Class:   warning,
	}
}

// Standard returns a standard flash message.
func Standard(message string) Message {
	return Message{
		Content: message,
		Class:   standard,
	}
}

// Session is an interface for typical sessions.
type Session interface {
	Save(*http.Request, http.ResponseWriter) error
	Flashes(vars ...string) []interface{}
}

func init() {
	// Magic goes here to allow serializing maps in securecookie
	// http://golang.org/pkg/encoding/gob/#Register
	// Source: http://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(Message{})
}

// SendFlashes allows retrieval of flash messages for using with Ajax.
func SendFlashes(w http.ResponseWriter, r *http.Request, sess Session) {
	flashes := PeekFlashes(w, r, sess)
	sess.Save(r, w)

	// There is no way for marshal to fail since it's a static type.
	js, _ := json.Marshal(flashes)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// PeekFlashes returns the flashes without destroying them.
func PeekFlashes(w http.ResponseWriter, r *http.Request, sess Session) []Message {
	var v []Message

	// Get the flashes for the template.
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v = make([]Message, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Message:
				v[i] = f.(Message)
			default:
				v[i] = Standard(fmt.Sprint(f))
			}

		}
	}

	return v
}
