package actor

import (
	"sync"

	"go.uber.org/zap"
)

type ActorSystem struct {
	root   Actor
	name   string
	logger *zap.Logger
	wg     sync.WaitGroup
}

func (as *ActorSystem) Spawn(a Actor) Actor {
	return as.root.Spawn(a)
}

func (as *ActorSystem) Shutdown() {
	// ch := make(chan int)
	// defer close(ch)
	// as.root.poisonCh <- PoisonPill{ch}
	// <-ch
	as.root.poison()
	as.wg.Done()
}

func (as *ActorSystem) Wait() {
	as.wg.Wait()
}

func NewActorSystem(name string, logger *zap.Logger) *ActorSystem {
	log := logger.Named("ActorSystem[" + name + "]")
	log.Debug("starting the actor system")
	as := &ActorSystem{name: name, root: NewReceiverActor("root", &EmptyReceiver{}, log.Named("root")), logger: log}
	as.wg.Add(1)
	return as
}
