package web

import (
	"encoding/gob"
)

func init() {
	// Register structs to inform encoding/gob how to handle it
	// so they can be stored in session data
	gob.Register(FlashMessage{})
	gob.Register(PendingRegistration{})
}
