package actor

import (
	"time"

	"go.uber.org/zap"
)

type IActor interface {
	Logger() *zap.Logger
	ID() string
	Name() string
	SendTo(t IActor, msg interface{})
	send(msg ActorMessage)
	Spawn(name string, receiver ActorReceiver)
	Schedule(message interface{}, duration time.Duration)
	ScheduleOnce(message interface{}, delay time.Duration)
	start()
}
