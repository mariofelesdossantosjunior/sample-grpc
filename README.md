# 📡 sample-grpc

Exemplo de implementação de serviço e cliente gRPC utilizando Protocol Buffers.

## 🛠️ Tecnologias

* **gRPC** – framework RPC de alto desempenho com HTTP/2 e suporte a múltiplas linguagens ([grpc.io][1], [github.com][2])
* **Protocol Buffers** – Interface Definition Language (IDL) e formato de serialização padrão ([grpc.io][1])

## 🎯 Visão Geral

Este projeto demonstra:

* **Definição da interface com `.proto`**
* **Servidor gRPC** capaz de atender chamadas remotas
* **Cliente gRPC** que consome o serviço do servidor

Exemplo básico inspirado no famoso "Hello World".

## 🚀 Começando

1. Clone o repositório

   ```bash
   git clone https://github.com/mariofelesdossantosjunior/sample-grpc.git
   cd sample-grpc
   ```
2. Gere o código a partir do `.proto`

   ```bash
   protoc --go_out=. --go-grpc_out=. proto/sample.proto
   ```
3. Execute o servidor

   ```bash
   go run cmd/server/main.go
   ```
4. Execute o cliente em outra janela

   ```bash
   go run cmd/client/main.go
   ```

## 📁 Estrutura do Projeto

```
├── cmd/
│   ├── server/        # Executável do servidor gRPC
│   └── client/        # Cliente gRPC que consome o serviço
├── proto/             # Arquivos .proto + código gerado
├── pkg/               # Lógica de negócio e handlers
├── go.mod             # Dependências Go
└── README.md
```

## 🧪 Testes

(Adicione testes unitários e integração, se houver)

* Para executar os testes:

```bash
go test ./...
```

## 🧩 Possíveis Extensões

* 🚦 Autenticação via tokens JWT
* 📶 Streaming (unário, servidor, cliente, bidirecional)
* 🛡️ TLS/SSL — segurança end-to-end
* 🎛️ Interceptadores (logging, métricas)
* 📤 gRPC-Gateway — expor REST a partir do gRPC

## 📚 Recursos

* Introdução ao gRPC com Protocol Buffers ([github.com][3])
* Exemplos oficiais em várias linguagens (C++, Go, Python…) ([github.com][4])

## 📝 Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

---

**Dicas para uso em projetos reais**:

* Utilize interceptadores gRPC para tratamento de erros, logging e métricas.
* Integre autenticação segura e TLS.
* Considere usar **gRPC-Gateway** para oferecer compatibilidade com APIs REST.
* Explore streaming e interceptadores — essenciais para aplicações robustas.

