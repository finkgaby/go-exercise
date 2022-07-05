package server

import (
	"Exercise/server/commons"
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
	Query     string
	SubjectId string
}

func QuerySerialize(queryToSerialize string) string {
	log.Println("Connect to NATS")
	nc, _ := nats.Connect(commons.NatsUrl)
	log.Println("Creates JetStreamContext")
	js, err := nc.JetStream()
	checkErr(err)

	query := Query{
		Query:     queryToSerialize,
		SubjectId: uuid.New().String(),
	}
	queryJSON, err := json.Marshal(query)
	js.Publish(pubSubjectName, queryJSON)

	log.Printf("Published queryJSON:%s to subjectName:%q", string(queryJSON), pubSubjectName)

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
			var query Query
			err := json.Unmarshal(msg.Data, &query)
			checkErr(err)
			log.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
			return query.Query
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
