## Documentação do Backend de Gerenciamento Universitário em Go
Este documento detalha o processo de configuração e desenvolvimento de um backend simples em Go para gerenciar alunos e matérias de uma universidade de Tecnologia da Informação. O projeto segue uma arquitetura em camadas (Models, Repositories, Services, Handlers) e utiliza SQLite como banco de dados, com a API REST exposta via gorilla/mux.

1. Configuração do Ambiente
Assumimos que você está utilizando o Ubuntu via WSL2 e tem o Go instalado e gerenciado via asdf.

1.1. Verificar a Versão do Go
Confirme a versão do Go instalada com go version.

1.2. Gerenciar o Go com asdf (Se necessário)
O asdf permite instalar e alternar entre diferentes versões do Go. Comandos como asdf plugin add golang, asdf install golang <versão> e asdf global golang <versão> são usados para gerenciar sua instalação do Go.

2. Início do Projeto
2.1. Criação do Diretório e Inicialização do Módulo Go
Inicie criando o diretório principal do seu espaço de trabalho Go e, em seguida, o diretório para o seu projeto college_api. Dentro dele, o comando go mod init college_api inicializa o módulo Go, criando o arquivo go.mod que gerencia as dependências do projeto.

2.2. Estrutura de Pastas
A organização do projeto é feita em pastas que representam cada camada da arquitetura, promovendo modularidade e clareza:

config/: Contém arquivos de configuração, como a inicialização do banco de dados.

models/: Define as estruturas de dados (modelos) que representam as entidades da sua aplicação (Aluno, Matéria).

repositories/: Abriga a lógica para interagir diretamente com o banco de dados.

services/: Contém a lógica de negócio e as validações, orquestrando as operações dos repositórios.

handlers/: Define as funções que recebem as requisições HTTP, processam-nas e enviam as respostas.

3. Implementação das Camadas
3.1. Camada de Modelos (models/)
Os arquivos nesta pasta (subject.go, student.go) definem as estruturas (structs) Go que representam as entidades Subject (matéria) e Student (aluno). Elas incluem campos como ID, Name, Enrollment, Year, etc., e as tags json para facilitar a serialização e desserialização para JSON nas operações da API.

3.2. Configuração do Banco de Dados (config/database.go)
O arquivo config/database.go é responsável por estabelecer a conexão com o banco de dados SQLite e criar as tabelas necessárias (students, subjects, e student_subjects para o relacionamento muitos-para-muitos). Ele expõe uma instância global do DB e funções para inicializar e fechar a conexão, garantindo que a aplicação possa persistir e recuperar dados.

3.3. Camada de Repositórios (repositories/)
Os arquivos subject_repository.go e student_repository.go nesta pasta implementam as operações de CRUD (Create, Read, Update, Delete) para Subject e Student diretamente no banco de dados. Eles contêm funções como CreateSubject, GetStudentByID, UpdateStudent, DeleteSubject, e métodos para gerenciar a associação entre alunos e matérias (AddSubjectToStudent, RemoveSubjectFromStudent). Essa camada isola a lógica de acesso a dados do restante da aplicação.

3.4. Camada de Serviços (services/)
Os arquivos subject_service.go e student_service.go contêm a lógica de negócio da aplicação. Eles recebem dados do handler, aplicam validações (por exemplo, se uma matéria já existe ou se um aluno tem uma matrícula válida), e orquestram as chamadas aos métodos dos repositórios. É a camada que garante a integridade e a consistência dos dados, tratando erros de forma mais amigável.

3.5. Camada de Handlers (handlers/)
Os arquivos subject_handler.go e student_handler.go são os controladores da API. Cada função handler (como CreateSubjectHandler, GetStudentByIDHandler) é responsável por:

Receber requisições HTTP (POST, GET, PUT, DELETE).

Extrair dados da requisição (JSON do corpo ou parâmetros da URL).

Chamar a lógica apropriada na camada de Serviços.

Construir e enviar a resposta HTTP de volta para o cliente, incluindo o corpo JSON e o código de status (ex: 200 OK, 201 Created, 400 Bad Request, 404 Not Found, 500 Internal Server Error).
Essa camada atua como a interface entre o mundo HTTP e a lógica interna da sua aplicação.

3.6. Arquivo Principal (main.go)
O main.go é o ponto de entrada da aplicação. Ele orquestra a inicialização de todas as camadas:

Inicia a conexão com o banco de dados.

Cria instâncias dos repositórios e serviços.

Cria instâncias dos handlers, injetando os serviços correspondentes.

Configura o roteador HTTP (gorilla/mux), mapeando as URLs da API para os handlers.

Inicia o servidor HTTP para escutar as requisições na porta 8080.

4. Dependências
Para que o projeto funcione, você precisará instalar as seguintes bibliotecas Go, que são gerenciadas pelo seu go.mod:

github.com/mattn/go-sqlite3: O driver para o banco de dados SQLite.

github.com/google/uuid: Para gerar IDs únicos (UUIDs) para alunos.

github.com/gorilla/mux: O roteador HTTP para sua API.

Essas dependências são instaladas com o comando go get <nome_da_dependencia>. Após adicionar todos os arquivos, o comando go mod tidy garante que seu go.mod e go.sum estejam sincronizados, adicionando e removendo as dependências conforme necessário.

5. Executando a Aplicação
Para iniciar o servidor backend, navegue até o diretório raiz do seu projeto (~/projetos_go/college_api) no terminal WSL2 e execute:

Bash

go run .
A mensagem Servidor iniciado em http://localhost:8080 indicará que sua API está no ar. Mantenha este terminal rodando enquanto você interage com a API.

6. Testando a API com curl
Para interagir com sua API, abra um segundo terminal WSL2 (ou use ferramentas como Postman/Insomnia) e envie requisições HTTP para os endpoints.

6.1. Endpoints de Matérias (/subjects)
Criar Matéria (POST): Envia os dados de uma nova matéria para o sistema.
curl -X POST -H "Content-Type: application/json" -d '{"id":"BSI101","name":"Programação Orientada a Objetos","year":1,"credits":8}' http://localhost:8080/subjects

Listar Todas as Matérias (GET): Retorna um array JSON com todas as matérias cadastradas.
curl http://localhost:8080/subjects

Buscar Matéria por ID (GET): Recupera os detalhes de uma matéria específica.
curl http://localhost:8080/subjects/{ID_DA_MATERIA}

Atualizar Matéria (PUT): Modifica os dados de uma matéria existente.
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Novo Nome Matéria"}' http://localhost:8080/subjects/{ID_DA_MATERIA}

Deletar Matéria (DELETE): Remove uma matéria do banco de dados.
curl -X DELETE http://localhost:8080/subjects/{ID_DA_MATERIA}

6.2. Endpoints de Alunos (/students)
Criar Aluno (POST): Adiciona um novo aluno, podendo associar matérias existentes (apenas pelo id).
curl -X POST -H "Content-Type: application/json" -d '{"enrollment":"20230001","name":"Cris Silva","current_year":1,"subjects":[{"id":"BSI101"}]}' http://localhost:8080/students

Listar Todos os Alunos (GET): Retorna um array JSON com todos os alunos, incluindo suas matérias.
curl http://localhost:8080/students

Buscar Aluno por ID (GET): Recupera os detalhes de um aluno específico, incluindo suas matérias.
curl http://localhost:8080/students/{ID_DO_ALUNO}

Atualizar Aluno (PUT): Modifica os dados principais de um aluno. Note que este endpoint não atualiza as matérias do aluno.
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Novo Nome Aluno","current_year":2}' http://localhost:8080/students/{ID_DO_ALUNO}

Deletar Aluno (DELETE): Remove um aluno do banco de dados.
curl -X DELETE http://localhost:8080/students/{ID_DO_ALUNO}

6.3. Endpoints de Relacionamento Aluno-Matéria
Associar Matéria a Aluno (POST): Adiciona uma matéria existente a um aluno específico.
curl -X POST http://localhost:8080/students/{ID_DO_ALUNO}/subjects/{ID_DA_MATERIA}

Remover Matéria de Aluno (DELETE): Remove a associação de uma matéria de um aluno.
curl -X DELETE http://localhost:8080/students/{ID_DO_ALUNO}/subjects/{ID_DA_MATERIA}

7. Limpando o Banco de Dados para Testes
Para resetar o banco de dados entre os testes e garantir que as operações de criação funcionem do zero, você pode simplesmente excluir o arquivo college.db do diretório raiz do seu projeto antes de cada execução:

Bash

rm college.db
go run .