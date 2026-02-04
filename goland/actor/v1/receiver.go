package actor

type SelfActor struct {
	Actor
}

type Receiver interface {

	// Setup puts initial set of messages for the self. Use it for scheduled messages.
	Setup(a SelfActor)

	// Receive processes the messages.
	// a is the reference to the actor.
	// am is the received message.
	Receive(a SelfActor, p Postcard)
}

type EmptyReceiver struct{}

func (*EmptyReceiver) Receive(SelfActor, Postcard) {

}

func (*EmptyReceiver) Setup(SelfActor) {}
