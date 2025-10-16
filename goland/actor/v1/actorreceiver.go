package actor

type Receiver interface {
	// Setup puts initial set of messages for the self. Use it for scheduled messages.
	Setup(a *ReceiverActor)

	// Receive processes the messages.
	// a is the reference to the actor.
	// am is the received message.
	Receive(a *ReceiverActor, am ActorMessage)
}

type EmptyReceiver struct{}

func (*EmptyReceiver) Receive(*ReceiverActor, ActorMessage) {

}

func (*EmptyReceiver) Setup(*ReceiverActor) {}
