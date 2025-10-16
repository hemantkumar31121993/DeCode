package actor

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ReceiverActor struct {
	id       string
	name     string
	closed   bool
	ch       chan ActorMessage
	poisonCh chan PoisonPill
	children map[string]Actor
	receiver Receiver
	logger   *zap.Logger
}

func (a *ReceiverActor) Logger() *zap.Logger {
	return a.logger
}

func (a *ReceiverActor) ID() string {
	return a.id
}

func (a *ReceiverActor) Name() string {
	return a.name
}

func (a *ReceiverActor) SendTo(t Actor, msg interface{}) error {
	return t.send(ActorMessage{msg, a})
}

func (a *ReceiverActor) poison() {
	ch := make(chan int)
	defer close(ch)
	a.poisonCh <- PoisonPill{ch}
	<-ch
}

func (a *ReceiverActor) send(msg ActorMessage) error {
	if !a.closed {
		a.ch <- msg
		return nil
	} else {
		return fmt.Errorf("actor has been shutdown")
	}
}

func (a *ReceiverActor) Spawn(sa Actor) Actor {
	a.children[sa.ID()] = sa
	return sa
}

func (a *ReceiverActor) Schedule(message interface{}, duration time.Duration) {
	t := time.NewTicker(duration)
	go func() {
		for range t.C {
			a.SendTo(a, message)
		}
	}()
}

func (a *ReceiverActor) ScheduleOnce(message interface{}, delay time.Duration) {
	t := time.NewTimer(delay)
	go func() {
		<-t.C
		a.SendTo(a, message)
	}()
}

func (a *ReceiverActor) start() {
	a.logger.Debug("Starting")

	for {
		select {
		case msg := <-a.poisonCh:
			a.logger.Debug("Got the PoisonPill")
			a.closed = true
			a.ch = nil
			a.poisonCh = nil
			ch := make(chan int)
			defer close(ch)
			// for _, ca := range a.children {
			// 	ca.poisonCh <- PoisonPill{ch}
			// }
			// time.Sleep(1 * time.Second)
			// for range a.children {
			// 	<-ch
			// }
			for _, ca := range a.children {
				ca.poison()
			}
			a.children = nil
			a.logger.Debug("Terminating")
			msg.waitCh <- 0
			return
		case msg := <-a.ch:
			a.logger.Sugar().Debugf("got the message: %s", msg)
			a.receiver.Receive(a, msg)
		}
	}
}

func NewReceiverActor(name string, receiver Receiver, logger *zap.Logger) Actor {
	ch := make(chan ActorMessage, 1024)
	id := makeID(name)
	a := &ReceiverActor{id: id,
		name:     name,
		ch:       ch,
		closed:   false,
		poisonCh: make(chan PoisonPill),
		children: make(map[string]Actor),
		receiver: receiver,
		logger:   logger}
	go a.start()
	go receiver.Setup(a)
	return a
}
