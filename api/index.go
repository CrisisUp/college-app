// api/index.go (NOVO NOME E CONTEÚDO)
package handler // O pacote deve ser 'handler' para Vercel Functions

import (
	"college_api/config" // Importa o config do seu módulo Go
	"college_api/handlers"
	"college_api/repositories"
	"college_api/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// O roteador Mux precisa ser uma variável global ou ser inicializado uma vez
// para que não seja re-inicializado em cada invocação da função serverless.
var router *mux.Router
var initOnce bool = false // Flag para garantir que a inicialização ocorra apenas uma vez

// Handler é a função de entrada para a Vercel Function.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Inicializa o roteador e as dependências APENAS UMA VEZ
	if !initOnce {
		initAPI()
		initOnce = true
	}
	// Servir a requisição usando o roteador inicializado
	router.ServeHTTP(w, r)
}

// initAPI inicializa todas as dependências da aplicação
func initAPI() {
	// A DATABASE_URL será definida via variável de ambiente da Vercel.
	config.InitDB() // Inicializa o banco de dados PostgreSQL
	// NOTE: defer config.CloseDB() não é usado em Serverless Functions
	// A conexão é mantida viva pela plataforma entre invocações.

	log.Println("Backend da universidade inicializando para Vercel Function...")

	// --- Inicializando Repositórios e Serviços ---
	subjectRepo := repositories.NewSubjectRepository()
	studentRepo := repositories.NewStudentRepository()
	teacherRepo := repositories.NewTeacherRepository()

	subjectService := services.NewSubjectService(subjectRepo)
	studentService := services.NewStudentService(studentRepo, subjectRepo)
	teacherService := services.NewTeacherService(teacherRepo)

	// --- Inicializando Handlers ---
	subjectHandler := handlers.NewSubjectHandler(subjectService)
	studentHandler := handlers.NewStudentHandler(studentService)
	teacherHandler := handlers.NewTeacherHandler(teacherService)

	// --- Configurando o Roteador Mux ---
	router = mux.NewRouter()

	// Rotas para Matérias
	router.HandleFunc("/subjects", subjectHandler.CreateSubjectHandler).Methods("POST")
	router.HandleFunc("/subjects", subjectHandler.GetAllSubjectsHandler).Methods("GET")
	router.HandleFunc("/subjects/{id}", subjectHandler.GetSubjectByIDHandler).Methods("GET")
	router.HandleFunc("/subjects/{id}", subjectHandler.UpdateSubjectHandler).Methods("PUT")
	router.HandleFunc("/subjects/{id}", subjectHandler.DeleteSubjectHandler).Methods("DELETE")

	// Rotas para Alunos
	router.HandleFunc("/students", studentHandler.CreateStudentHandler).Methods("POST")
	router.HandleFunc("/students", studentHandler.GetAllStudentsHandler).Methods("GET")
	router.HandleFunc("/students/{id}", studentHandler.GetStudentByIDHandler).Methods("GET")
	router.HandleFunc("/students/{id}", studentHandler.UpdateStudentHandler).Methods("PUT")
	router.HandleFunc("/students/{id}", studentHandler.DeleteStudentHandler).Methods("DELETE")

	// Rotas para associação Aluno-Matéria
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}", studentHandler.AddSubjectToStudentHandler).Methods("POST")
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}", studentHandler.RemoveSubjectFromStudentHandler).Methods("DELETE")

	// --- ROTAS PARA PROFESSORES ---
	router.HandleFunc("/teachers", teacherHandler.CreateTeacherHandler).Methods("POST")
	router.HandleFunc("/teachers", teacherHandler.GetAllTeachersHandler).Methods("GET")
	router.HandleFunc("/teachers/{id}", teacherHandler.GetTeacherByIDHandler).Methods("GET")
	router.HandleFunc("/teachers/{id}", teacherHandler.UpdateTeacherHandler).Methods("PUT")
	router.HandleFunc("/teachers/{id}", teacherHandler.DeleteTeacherHandler).Methods("DELETE")

	// --- Configuração do CORS ---
	// Em Vercel Functions, o CORS deve ser tratado pelo 'vercel.json' nos headers,
	// mas é bom ter no código também como fallback ou para testes locais.
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            false, // Defina como false em produção
	})

	// Aplica o middleware CORS ao seu roteador
	router.Use(mux.MiddlewareFunc(corsHandler.Handler)) // Usa mux.MiddlewareFunc para integrar o handler como middleware

	log.Println("Backend da universidade inicializado com sucesso para Vercel Function!")

	// Remover o http.ListenAndServe pois a Vercel Function não é um servidor tradicional
	// log.Fatal(srv.ListenAndServe())
}
