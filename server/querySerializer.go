package server

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

const (
	pubSubjectName = "QUERY.unserialized"
	subSubjectName = "QUERY.serialized"
)

type Query struct {
	query     string
	subjectId string
}

func QuerySerialize(queryToSerialize string) string {
	var serializedQuery string
	log.Println("Connect to NATS")
	nc, _ := nats.Connect("demo.nats.io")
	log.Println("Creates JetStreamContext")
	js, err := nc.JetStream()
	checkErr(err)

	query := Query{
		query:     queryToSerialize,
		subjectId: uuid.New().String(),
	}
	queryJSON, err := json.Marshal(query.query)
	js.Publish(pubSubjectName, queryJSON)

	log.Printf("Published queryJSON:%s to subjectName:%q", string(queryJSON), pubSubjectName)
	//js.Subscribe(subSubjectName, func(msg *nats.Msg) {
	//	msg.Ack()
	//	var query string
	//	err := json.Unmarshal(msg.Data, &query)
	//	checkErr(err)
	//	log.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
	//	serializedQuery = query
	//}, nats.Durable("monitor"), nats.ManualAck())

	sub, _ := js.PullSubscribe(subSubjectName, "queryReviewSubscriber", nats.PullMaxWaiting(1))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			checkErr(err)
		default:
		}
		msgs, _ := sub.Fetch(1, nats.Context(ctx))
		for _, msg := range msgs {
			msg.Ack()
			err := json.Unmarshal(msg.Data, &serializedQuery)
			checkErr(err)
			log.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
			return serializedQuery
		}
	}

	return serializedQuery
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
