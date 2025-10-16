package actor

import (
	"time"

	"go.uber.org/zap"
)

type Actor interface {
	Logger() *zap.Logger
	ID() string
	Name() string
	SendTo(t Actor, msg interface{}) error
	send(msg ActorMessage) error
	poison()
	Spawn(a Actor) Actor
	Schedule(message interface{}, duration time.Duration)
	ScheduleOnce(message interface{}, delay time.Duration)
	start()
}
