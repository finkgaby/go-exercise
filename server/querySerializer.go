package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"log"
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
	queryJSON, _ := json.Marshal(query.query)
	js.Publish(fmt.Sprintf("(%s_%s)", pubSubjectName, query.subjectId), queryJSON)

	log.Printf("Published queryJSON:%s to subjectName:%q", string(queryJSON), pubSubjectName)
	js.Subscribe(fmt.Sprintf("(%s_%s)", subSubjectName, query.subjectId), func(msg *nats.Msg) {
		msg.Ack()
		var query Query
		err := json.Unmarshal(msg.Data, &query)
		checkErr(err)
		log.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		serializedQuery = query.query
	}, nats.Durable("monitor"), nats.ManualAck())

	return serializedQuery
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
