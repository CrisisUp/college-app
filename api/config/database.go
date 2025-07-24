// config/database.go
package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite3
)

// DB é a instância global do banco de dados.
var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados SQLite e cria as tabelas.
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados SQLite estabelecida com sucesso!")

	createTables()
}

// createTables cria as tabelas Student e Subject se elas não existirem.
func createTables() {
	createStudentsTableSQL := `
    CREATE TABLE IF NOT EXISTS students (
        id TEXT PRIMARY KEY,
        enrollment TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        current_year INTEGER NOT NULL
    );`

	createSubjectsTableSQL := `
    CREATE TABLE IF NOT EXISTS subjects (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        year INTEGER NOT NULL,
        credits INTEGER NOT NULL
    );`

	// Tabela para relacionamento muitos-para-muitos entre alunos e matérias
	createStudentSubjectsTableSQL := `
    CREATE TABLE IF NOT EXISTS student_subjects (
        student_id TEXT NOT NULL,
        subject_id TEXT NOT NULL,
        PRIMARY KEY (student_id, subject_id),
        FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE
    );`

	_, err := DB.Exec(createStudentsTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela students: %v", err)
	}
	_, err = DB.Exec(createSubjectsTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela subjects: %v", err)
	}
	_, err = DB.Exec(createStudentSubjectsTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela student_subjects: %v", err)
	}

	log.Println("Tabelas verificadas/criadas com sucesso!")
}

// CloseDB fecha a conexão com o banco de dados.
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Conexão com o banco de dados SQLite fechada.")
	}
}
