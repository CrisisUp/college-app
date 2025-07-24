// main.go
package main

import (
	"college_api/config"
	"college_api/handlers"
	"college_api/repositories"
	"college_api/services"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Inicializa o banco de dados SQLite.
	config.InitDB("./college.db")
	defer config.CloseDB() // Garante que a conexão com o DB seja fechada ao final da execução.

	log.Println("Backend da universidade iniciado!")

	// --- Inicializando Repositórios e Serviços ---
	subjectRepo := repositories.NewSubjectRepository()
	studentRepo := repositories.NewStudentRepository()

	subjectService := services.NewSubjectService(subjectRepo)
	studentService := services.NewStudentService(studentRepo, subjectRepo)

	// --- Inicializando Handlers ---
	subjectHandler := handlers.NewSubjectHandler(subjectService)
	studentHandler := handlers.NewStudentHandler(studentService)

	// --- Configurando o Roteador Mux ---
	router := mux.NewRouter()

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

	// --- Iniciando o Servidor HTTP ---
	addr := ":8080" // A porta em que o servidor irá escutar
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Servidor iniciado em http://localhost%s", addr)
	// log.Fatal irá logar o erro e encerrar a aplicação se houver um problema ao iniciar o servidor
	log.Fatal(srv.ListenAndServe())

	// A função main não será mais encerrada automaticamente como antes,
	// pois o servidor http.ListenAndServe() é um processo de bloqueio.
	// Ele só será encerrado se houver um erro fatal ou se o processo for manualmente interrompido (Ctrl+C).
}
