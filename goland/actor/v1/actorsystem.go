package actor

import (
	"sync"

	"go.uber.org/zap"
)

type ActorSystem struct {
	root   *Actor
	name   string
	logger *zap.Logger
	wg     sync.WaitGroup
}

func (as *ActorSystem) Spawn(name string, receiver ActorReceiver) *Actor {
	return as.root.Spawn(name, receiver)
}

func (as *ActorSystem) Shutdown() {
	ch := make(chan int)
	defer close(ch)
	as.root.poisonCh <- PoisonPill{ch}
	<-ch
	as.wg.Done()
}

func (as *ActorSystem) Wait() {
	as.wg.Wait()
}

func NewActorSystem(name string, logger *zap.Logger) *ActorSystem {
	log := logger.Named("ActorSystem[" + name + "]")
	log.Debug("starting the actor system")
	as := &ActorSystem{name: name, root: newActor("/root", "root", &EmptyReceiver{}, log.Named("root")), logger: log}
	as.wg.Add(1)
	return as
}
