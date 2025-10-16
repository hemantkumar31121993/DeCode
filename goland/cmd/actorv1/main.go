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

func (t *TestActorReciver) Receive(a *actor.ReceiverActor, msg actor.ActorMessage) {
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

type Producer struct {
	counter  int
	consumer actor.Actor
}

func (*Producer) Setup(a actor.Actor) {
	a.Schedule(Transaction{PRODUCE, ""}, 5*time.Second)
}

func (p *Producer) Receive(a actor.Actor, msg actor.ActorMessage) {
	t := msg.Message.(Transaction)
	switch t.T {
	case PRODUCE:
		// p.consumer.Send(actor.MakeMessage(Transaction{CONSUME, fmt.Sprintf("product:%d", p.counter)}, a))
		a.SendTo(p.consumer, Transaction{CONSUME, fmt.Sprintf("product:%d", p.counter)})
		p.counter++
	case COMSUME_ACK:
		a.Logger().Info("consumer consumed product " + t.Product)
	}
}

type ConsumerAck struct {
	Producer actor.Actor
	Product  string
}

type Consumer struct {
}

func (*Consumer) Setup(actor.Actor) {

}

func (*Consumer) Receive(a actor.Actor, msg actor.ActorMessage) {
	switch msg.Message.(type) {
	case Transaction:
		a.ScheduleOnce(ConsumerAck{msg.Sender, msg.Message.(Transaction).Product}, 3*time.Second)

	case ConsumerAck:
		mack := msg.Message.(ConsumerAck)
		// msg.Message.(ConsumerAck).Producer.Send(actor.MakeMessage(Transaction{COMSUME_ACK, msg.Message.(ConsumerAck).Product}, a))
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

	consumer := as.Spawn(actor.NewReceiverActor("consumer", &Consumer{}, logger.Named("consumer")))

	as.Spawn(actor.NewReceiverActor("producer", &Producer{0, consumer}, logger.Named("producer")))

	go func() {
		<-sigChan
		as.Shutdown()
	}()

	as.Wait()
}
