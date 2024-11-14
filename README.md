# library-go-api-rest
API feita em Go para gerenciar um site de empréstimos de livros de uma biblioteca.

## Requisitos
* [Go](https://go.dev/dl/)

## Instalação
1. Clone o repositório do Github.
```
git clone https://github.com/LucasPeixer/library-go-api-rest
cd library-go-api-rest
```
2. Instale as dependências.
```
go mod tidy
```

## Variáveis de ambiente
* `DB_DSN`: URL para conexão com o banco de dados;
* `JWT_KEY`: Chave secreta para autenticação JWT;
* `PORT`: Opcional, utilizada para atender as requisições HTTP.

## Banco de dados 
A API requer conexão com um banco de dados **PostgreSQL**, seja ele local ou na nuvem.

Será necessário definir a **variável de ambiente** `DB_DSN`, que contém as informações de conexão com o banco de dados. 
Você pode criar essa variável em um arquivo .env ou defini-la diretamente no seu ambiente.

O valor da variável `DB_DSN` segue o seguinte formato:
```
postgres://user:password@host:port/database
```

* `user`: nome de usuário do banco de dados;
* `senha`: senha do banco de dados;
* `host`: endereço do servidor onde o banco de dados está hospedado (por exemplo, localhost ou um endereço na nuvem);
* `port`: número da porta do PostgreSQL (geralmente 5432);
* `database`: nome do banco de dados que será utilizado pela API.

Exemplo de DSN para conexão:
```
postgresql://postgres:NDAED123dDGFr@gravely-incredible-emerald.data-1.use1.tembo.io:5432/postgres
```

Feito isso, será necessário criar as tabelas do banco de dados. Acesse o arquivo `DDL.sql` e execute-o.
