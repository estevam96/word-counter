package main

import (
	"fmt"
	"log"
	"time"

	message "Trabalho_2/proto"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

type result struct {
	r map[string]int64
}

func main() {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	errMsg(err, " Falha ao conectar")
	defer conn.Close()

	ch, err := conn.Channel()
	errMsg(err, "Falha ao abrir canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"score-result",
		false,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao declarar Fila de resposta")

	fim, err := ch.QueueDeclare(
		"fim",
		false,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao declara fial ")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao Consumir a fila")

	const duration = 5 * time.Second
	timer := time.NewTimer(duration)

	cont := newWords()
	forever := make(chan bool)

	go func() {
		for {
			res := message.Result{}

			select {
			case d := <-msgs:
				timer.Reset(duration)
				err := proto.Unmarshal(d.Body, &res)
				errMsg(err, "Falha ao Deserializar")
				for w, o := range res.Found {
					cont.r[w] = cont.r[w] + o
				}
				for w, o := range res.Found {
					fmt.Printf(" => %s : %d\n", w, o)
				}
			case <-timer.C:
				println("timeout")
				rf := message.ResultFinal{
					Found: cont.r,
				}

				data, err := proto.Marshal(&rf)
				errMsg(err, "Falha ao serializar result")
				if len(data) > 0 {
					println("enviou")
					err = ch.Publish(
						"",
						fim.Name,
						false,
						false,
						amqp.Publishing{
							ContentType: "application/protobuf",
							Body:        data,
						})
					errMsg(err, "Falha ao publicar")

				}
				cont = newWords()
			}
		}
	}()
	<-forever
}

func newWords() *result {
	return &result{r: map[string]int64{}}
}

func errMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
	}
}
