# Word Counter Golang

*Requer docker e docker compose

## Preparando
Para a executar é necessário possui a lib:
* [grpc](https://godoc.org/google.golang.org/grpc) lib em Golang para trabalhar com grpc
* [protoc-gen-go](https://github.com/golang/protobuf/) lib para uso do  protocol buffers em Golang
* [amqp](https://github.com/streadway/amqp) Cliente Rabbitmq para Golang

Use os comando no terminal
* para instalar o grpc use:\
  ```go get -u google.golang.org/grpc```

* para instalar o protoco-gen-go use:\
  ```go get -u github.com/golang/protobuf/protoc-gen-go```

* para instalar o amqp use:\
  ```go get -u github.com/streadway/amqp```

Agora é necessário compilar o Read.proto que esta na pasta proto.

Execute o comando:\
  ```cd word_counter```

logo apos\
  ```protoc -I . proto/Read.proto --go_out=plugins=grpc:.```

execute um container com o rabbitmq. Na pasta do projeto execute o comando:\
  ```docker-compose up```

## Executando projeto

Inicie executando o MapReducer.go e não o Client.go (execute o Client.go por ultimo) pois o grpc necessita que o servidor rpc já tenha iniciado.\
```go run MapReducer.go```

Em seguida execute na ordem que preferir o Separator.go, Counter.go ou e o Result.go.

O Counter.go recebe três argumentos (pode ser os três, apenas dois ou somente, mas pelo menos um argumento deve ser passado).

Argumento | Descrição
--------| --------
a- h | realiza contagem de todas as palavras que possui as iniciais ente a até h (seja minúscula ou maiúscula).
i-p | realiza contagem de todas as palavras que possui as iniciais ente i até p (seja minúscula ou maiúscula).
q-z | realiza contagem de todas as palavras que possui as iniciais ente q até z (seja minúscula ou maiúscula), e também caracteres especiais.


O Resultado será impresso no arquivo chamado saida.txt