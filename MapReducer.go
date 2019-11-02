package main

import (
	"context"
	"log"
	"net"
	"regexp"
	"strings"

	message "Trabalho_2/proto"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Cont(ctx context.Context, in *message.ClientRequest) (*message.MasterResponse, error) {

	lines := strings.Split(in.GetWord(), "\n")

	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	errMsg(err, "Falha ao connecta ao RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	errMsg(err, "Falha ao criar canal")
	defer ch.Close()

	fila, err := ch.QueueDeclare(
		"word",
		true,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao cria fila")

	q, err := ch.QueueDeclare(
		"fim",
		false,
		false,
		false,
		false,
		nil,
	)
	errMsg(err, "Falha ao declara fial ")

	for _, l := range lines {
		m := removeSpecialCharacters(l)
		palavra := message.Worker{
			Word: m,
		}

		data, _ := proto.Marshal(&palavra)

		err = ch.Publish(
			"",
			fila.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/protobuf",
				Body:         data,
			})
		errMsg(err, "Falha ao publicar message")
	}

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

	mresult := message.ResultFinal{}
	err = proto.Unmarshal((<-msgs).Body, &mresult)

	resposta := message.MasterResponse{}
	for a, b := range mresult.Found {
		test := message.MapResponse{
			Palavra:    a,
			Ocorrencia: b,
		}
		resposta.Mr = append(resposta.Mr, &test)
	}

	return &resposta, nil
}

func main() {

	// GRPC
	lis, err := net.Listen("tcp", ":5040")
	errMsg(err, "erro ao conectar")

	s := grpc.NewServer(grpc.MaxRecvMsgSize(1024*5000), grpc.MaxSendMsgSize(1024*5000))
	message.RegisterCountServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("falhou ao servir: %v", err)
	}
}

func errMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
	}
}

func removeSpecialCharacters(s string) string {
	reg, err := regexp.Compile("[[:punct:]0-9]+")
	errMsg(err, "Erro ao Copila Regex")
	processedText := reg.ReplaceAllString(s, " ")

	return processedText
}
