package actor

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Actor struct {
	id       string
	name     string
	closed   bool
	ch       chan ActorMessage
	poisonCh chan PoisonPill
	children map[string]*Actor
	receiver ActorReceiver
	logger   *zap.Logger
}

func (a *Actor) Logger() *zap.Logger {
	return a.logger
}

func (a *Actor) ID() string {
	return a.id
}

func (a *Actor) Name() string {
	return a.name
}

func (a *Actor) SendTo(t *Actor, msg interface{}) error {
	return t.send(ActorMessage{msg, a})
}

func (a *Actor) poison() {
	ch := make(chan int)
	defer close(ch)
	a.poisonCh <- PoisonPill{ch}
	<-ch
}

func (a *Actor) send(msg ActorMessage) error {
	if !a.closed {
		a.ch <- msg
		return nil
	} else {
		return fmt.Errorf("actor has been shutdown")
	}
}

func (a *Actor) Spawn(name string, receiver ActorReceiver) *Actor {
	id := makeID(name)
	sa := newActor(a.id+"/"+id, name, receiver, a.logger.Named(id))
	a.children[id] = sa
	return sa
}

func (a *Actor) Schedule(message interface{}, duration time.Duration) {
	t := time.NewTicker(duration)
	go func() {
		for range t.C {
			a.SendTo(a, message)
		}
	}()
}

func (a *Actor) ScheduleOnce(message interface{}, delay time.Duration) {
	t := time.NewTimer(delay)
	go func() {
		<-t.C
		a.SendTo(a, message)
	}()
}

func (a *Actor) start() {
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

func newActor(id, name string, receiver ActorReceiver, logger *zap.Logger) *Actor {
	ch := make(chan ActorMessage, 1024)
	a := &Actor{id: id,
		name:     name,
		ch:       ch,
		closed:   false,
		poisonCh: make(chan PoisonPill),
		children: make(map[string]*Actor),
		receiver: receiver,
		logger:   logger}
	go a.start()
	go receiver.Setup(a)
	return a
}
