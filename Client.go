package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	message "Trabalho_2/proto"

	"google.golang.org/grpc"
)

func main() {
	// conect with master server
	conn, err := grpc.Dial("localhost:5040", grpc.WithInsecure())
	errMsg(err, "falha ao conectar com o Master server")

	defer conn.Close()

	c := message.NewCountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	file, err := os.Open("./Book.txt")
	errMsg(err, "Erro ao abrir arquivo")
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	t := string(b)
	errMsg(err, "Erro ao ler arquivo")

	r, err := c.Cont(ctx, &message.ClientRequest{Word: t})
	errMsg(err, "ocorreu um erro")

	A := r.GetMr()

	saida, err := os.Create("./saida.txt")
	errMsg(err, "Falha ao ler aquivo")
	defer saida.Close()

	for _, ocorrencia := range A {
		Ocorr := strconv.FormatInt(ocorrencia.Ocorrencia, 10)
		_, err := saida.WriteString(ocorrencia.Palavra + " : " + Ocorr + "\n")
		errMsg(err, "Falha ao escrever arquivo")

		saida.Sync()
		log.Printf("%v : %d", ocorrencia.Palavra, ocorrencia.Ocorrencia)
	}
}

func removeSpecialCharacters(s string) string {
	reg, err := regexp.Compile("[[:punct:]0-9]+")
	errMsg(err, "Erro ao Copila Regex")
	processedText := reg.ReplaceAllString(s, " ")

	return processedText
}

func errMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
	}
}
