package actor

import "fmt"

// ActorMessage is the container for the communication message.
type ActorMessage struct {

	// Message contains the actual message.
	Message interface{}

	// Sender is the reference to the sender of the message.
	Sender *Actor
}

func (am ActorMessage) String() string {
	return fmt.Sprintf("actor.ActorMessage{Message: %#v, Sender: %s}", am.Message, am.Sender.id)
}

type PoisonPill struct {
	waitCh chan int
}
