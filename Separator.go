package main

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	message "Trabalho_2/proto"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

func main() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && unicode.IsSpace(c)
	}
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	errMsg(err, "Falha ao conectar com o RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	errMsg(err, "Falha ao Abrir canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"word",
		true,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao declar a fila")

	err = ch.Qos(
		1,
		0,
		false,
	)
	errMsg(err, "Falha ao criar Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao registra o consumidor")

	ch2, err := conn.Channel()
	errMsg(err, "Falha ao criar canal 2")
	defer conn.Close()

	err = ch2.ExchangeDeclare(
		"conting",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	errMsg(err, "Falha ao declara fila 2")

	forever := make(chan bool)

	for d := range msgs {
		res := message.Worker{}

		err = proto.Unmarshal([]byte((d).Body), &res)

		for _, word := range strings.FieldsFunc(string(res.Word), f) {
			palavra := message.Separador{
				Word: word,
			}
			data, err := proto.Marshal(&palavra)
			errMsg(err, "Erro ao serializar data")
			println(" => " + word)
			if equalRegex(string(word[0]), "[abcdefghABCDEFGH]") {
				publish(ch2, "a-h", []byte(data))
			} else if equalRegex(string(word[0]), "[ijklmnopIJKLMNOP]") {

				publish(ch2, "i-p", []byte(data))
			} else {
				publish(ch2, "q-z", []byte(data))
			}

		}

		d.Ack(false)
	}

	<-forever
}

func equalRegex(palavra string, regex string) bool {
	re := regexp.MustCompile(regex)
	return re.Match([]byte(palavra))
}

func removeSpecialCharacters(s string) string {

	reg, err := regexp.Compile("[[:punct:]0-9]+")
	errMsg(err, "Erro ao Copila Regex")
	processedText := reg.ReplaceAllString(s, " ")
	return processedText
}

func publish(ch *amqp.Channel, routing string, data []byte) {

	err := ch.Publish(
		"conting",
		routing,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        data,
		})
	errMsg(err, "Erro ao publicar")
}

func errMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
	}
}
