package actor

type Receiver interface {

	// Setup puts initial set of messages for the self. Use it for scheduled messages.
	Setup(a Actor)

	// Receive processes the messages.
	// a is the reference to the actor.
	// am is the received message.
	Receive(a Actor, p Postcard)
}

type EmptyReceiver struct{}

func (*EmptyReceiver) Receive(Actor, Postcard) {

}

func (*EmptyReceiver) Setup(Actor) {}
