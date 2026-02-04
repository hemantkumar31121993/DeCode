package actor

import "fmt"

// Postcard is the container for the communication message.
type Postcard struct {

	// Message contains the actual message.
	Message interface{}

	// Sender is the reference to the sender of the message.
	Sender Actor
}

func (p Postcard) String() string {
	return fmt.Sprintf("actor.Postcard{Message: %#v, Sender: %s}", p.Message, p.Sender.ID())
}

type PoisonPill struct {
	waitCh chan int
}
