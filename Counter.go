package main

import (
	"fmt"
	"log"
	"os"
	"time"

	message "Trabalho_2/proto"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

type cont struct {
	found map[string]int64
}

func main() {

	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	errMsg(err, "Falha ao Conectar ao Rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	errMsg(err, "Erro ao abrir canal")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"conting",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao declara Exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	errMsg(err, "Faile a declara fila")

	if len(os.Args) < 2 {
		log.Printf("Informe os tipos que deseja receber: [a-h] [i-p] [q-z]")
		os.Exit(0)
	}

	for _, a := range os.Args[1:] {
		err := ch.QueueBind(
			q.Name,
			a,
			"conting",
			false,
			nil)
		errMsg(err, "Falha ao buscar fila")
	}

	word, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	r, err := ch.QueueDeclare(
		"score-result",
		false,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao declarar Fila de resposta")

	w := newWords()
	errMsg(err, "Falha ao consumir a fila")

	forever := make(chan bool)

	const duration = 3 * time.Second
	timer := time.NewTimer(duration)

	go func() {

		for {
			select {
			case d := <-word:
				timer.Reset(duration)
				res := message.Separador{}

				err = proto.Unmarshal([]byte(d.Body), &res)
				errMsg(err, "Falha ao desserilizar proto")

				log.Printf(" =>  %s\n ", res.Word)

				w.found[res.Word]++

			case <-timer.C:
				result := message.Result{
					Found: w.found,
				}
				data, err := proto.Marshal(&result)
				errMsg(err, "Falha ao serializar result")

				err = ch.Publish(
					"",
					r.Name,
					false,
					false,
					amqp.Publishing{
						ContentType: "application/protobuf",
						Body:        data,
					})
				errMsg(err, "Falha ao publicar")

				w = newWords()
			}
		}

	}()

	<-forever
	println("acabou o for")
}

func errMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
	}
}

func newWords() *cont {
	return &cont{found: map[string]int64{}}
}

func imprime(p *cont) {
	for word, count := range p.found {
		if count > 1 {
			fmt.Printf("%s : %d \n		", word, count)
		}
	}
}
