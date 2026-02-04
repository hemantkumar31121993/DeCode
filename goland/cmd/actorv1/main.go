package main

import (
	"fmt"
	"goland/actor/v1"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type TestActorReciver struct {
	counter int
}

type getCounter string

const GetCounter getCounter = "GET_COUNTER"

func (t *TestActorReciver) Setup(a *actor.ReceiverActor) {

}

func (t *TestActorReciver) Receive(a *actor.ReceiverActor, msg actor.Postcard) {
	t.counter++
	switch msg.Message.(type) {
	case string:
		a.Logger().Sugar().Debugf("string message: %#v\n", msg.Message)
	case getCounter:
		a.Logger().Sugar().Debugf("counter: %v\n", t.counter)
	default:
	}
}

type TransactionType string

const PRODUCE TransactionType = "PRODUCE"
const CONSUME TransactionType = "CONSUME"
const COMSUME_ACK TransactionType = "CONSUME_ACK"

type Transaction struct {
	T       TransactionType
	Product string
}

type ProducerReceiver struct {
	counter  int
	consumer actor.Actor
}

func (*ProducerReceiver) Setup(a actor.Actor) {
	a.Schedule(Transaction{PRODUCE, ""}, 5*time.Second)
}

func (p *ProducerReceiver) Receive(a actor.Actor, pc actor.Postcard) {
	t := pc.Message.(Transaction)
	switch t.T {
	case PRODUCE:
		a.SendTo(p.consumer, Transaction{CONSUME, fmt.Sprintf("product:%d", p.counter)})
		p.counter++
	case COMSUME_ACK:
		a.Logger().Info("consumer " + pc.Sender.ID() + " consumed product " + t.Product)
	}
}

type ConsumerAck struct {
	Producer actor.Actor
	Product  string
}

type ConsumerReceiver struct{}

func (*ConsumerReceiver) Setup(actor.Actor) {

}

func (*ConsumerReceiver) Receive(a actor.Actor, msg actor.Postcard) {
	switch msg.Message.(type) {
	case Transaction:
		a.ScheduleOnce(ConsumerAck{msg.Sender, msg.Message.(Transaction).Product}, 3*time.Second)

	case ConsumerAck:
		mack := msg.Message.(ConsumerAck)
		a.SendTo(mack.Producer, Transaction{COMSUME_ACK, msg.Message.(ConsumerAck).Product})
	}
}

func main() {

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger, _ := zap.NewDevelopment()

	as := actor.NewActorSystem("testing", logger)

	// ta := as.Spawn("TestActor", &TestActorReciver{0})

	// as.Spawn("TestActorChildren", &actor.EmptyReceiver{})

	// ticker := time.NewTicker(1 * time.Second)

	// go func() {
	// 	for t := range ticker.C {
	// 		err := ta.Send(actor.MakeMessage(t.String(), nil))
	// 		if err != nil {
	// 			fmt.Printf("err: %#v\n", err)
	// 			return
	// 		}
	// 	}
	// }()

	// ta.Send(actor.MakeMessage("message", nil))
	// ta.Send(actor.MakeMessage(GetCounter, nil))
	// ta.Send(actor.MakeMessage(9, nil))
	// ta.Send(actor.MakeMessage(GetCounter, nil))

	consumerRRBalancer := as.Spawn(actor.NewRoundRobinBalancer("consumer balancer", logger.Named("consumerRRBalancer")))

	consumerRRBalancer.Spawn(actor.NewReceiverActor("consumer1", &ConsumerReceiver{}, logger.Named("consumer1")))
	consumerRRBalancer.Spawn(actor.NewReceiverActor("consumer2", &ConsumerReceiver{}, logger.Named("consumer2")))
	consumerRRBalancer.Spawn(actor.NewReceiverActor("consumer3", &ConsumerReceiver{}, logger.Named("consumer3")))
	consumerRRBalancer.Spawn(actor.NewReceiverActor("consumer4", &ConsumerReceiver{}, logger.Named("consumer4")))

	as.Spawn(actor.NewReceiverActor("producer", &ProducerReceiver{counter: 0, consumer: consumerRRBalancer}, logger.Named("producer")))

	go func() {
		<-sigChan
		as.Shutdown()
	}()

	as.Wait()
}
