package actor

import (
	"time"

	"go.uber.org/zap"
)

type RoundRobinBalancer struct {
	id          string
	name        string
	closed      bool
	poisonCh    chan PoisonPill
	downstream  []Actor
	activeIndex int
	logger      *zap.Logger
}

func (r *RoundRobinBalancer) Logger() *zap.Logger {
	return r.logger
}

func (r *RoundRobinBalancer) ID() string {
	return r.id
}

func (r *RoundRobinBalancer) Name() string {
	return r.name
}

func (r *RoundRobinBalancer) Closed() bool {
	return r.closed
}

func (r *RoundRobinBalancer) SendTo(t Actor, msg interface{}) error {
	r.logger.Debug("RoundRobinBalancer does not send mesaage to another actor")
	return nil
}

func (r *RoundRobinBalancer) queue(msg Postcard) error {
	for ad := r.downstream[r.activeIndex]; !ad.Closed(); {
		r.activeIndex = (r.activeIndex + 1) % len(r.downstream)
		ad.queue(msg)
		return nil
	}
	return nil
}

func (r *RoundRobinBalancer) poison() {
	ch := make(chan int)
	defer close(ch)
	r.poisonCh <- PoisonPill{ch}
	<-ch
}

func (r *RoundRobinBalancer) Spawn(a Actor) Actor {
	r.downstream = append(r.downstream, a)
	go a.start()
	go a.setup()
	return a
}

func (r *RoundRobinBalancer) Schedule(message interface{}, duration time.Duration) {

}

func (r *RoundRobinBalancer) ScheduleOnce(message interface{}, delay time.Duration) {

}

func (r *RoundRobinBalancer) start() {
	r.logger.Debug("Starting")

	msg := <-r.poisonCh
	r.logger.Debug("Got the PoisonPill")
	r.closed = true
	r.poisonCh = nil
	// ch := make(chan int)
	// defer close(ch)
	// for _, ca := range a.children {
	// 	ca.poisonCh <- PoisonPill{ch}
	// }
	// time.Sleep(1 * time.Second)
	// for range a.children {
	// 	<-ch
	// }
	for _, da := range r.downstream {
		da.poison()
	}
	r.downstream = nil
	r.logger.Debug("Terminating")
	msg.waitCh <- 0
}

func (r *RoundRobinBalancer) setup() {

}

func NewRoundRobinBalancer(name string, logger *zap.Logger) *RoundRobinBalancer {
	id := makeID(name)
	rr := &RoundRobinBalancer{
		id:          id,
		name:        name,
		closed:      false,
		poisonCh:    make(chan PoisonPill),
		downstream:  []Actor{},
		activeIndex: 0,
		logger:      logger,
	}
	return rr
}
