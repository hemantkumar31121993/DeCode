package actor

type ActorReceiver interface {
	// Setup puts initial set of messages for the self. Use it for scheduled messages.
	Setup(a *Actor)

	// Receive processes the messages.
	// a is the reference to the actor.
	// am is the received message.
	Receive(a *Actor, am ActorMessage)
}

type EmptyReceiver struct{}

func (*EmptyReceiver) Receive(*Actor, ActorMessage) {

}

func (*EmptyReceiver) Setup(*Actor) {}
