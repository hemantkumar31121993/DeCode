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
	as.root.poison()
	as.wg.Done()
}

func (as *ActorSystem) Wait() {
	as.wg.Wait()
}

func NewActorSystem(name string, logger *zap.Logger) *ActorSystem {
	log := logger.Named("ActorSystem[" + name + "]")
	log.Debug("starting the actor system")
	root := NewReceiverActor("root", &EmptyReceiver{}, log.Named("root"))
	as := &ActorSystem{name: name, root: root, logger: log}
	go root.start()
	go root.setup()
	as.wg.Add(1)
	return as
}
