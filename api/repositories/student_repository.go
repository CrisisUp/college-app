// repositories/student_repository.go
package repositories

import (
	"college_api/config"
	"college_api/models"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

// StudentRepository define as operações de CRUD para alunos.
type StudentRepository struct {
	db *sql.DB
}

// NewStudentRepository cria uma nova instância de StudentRepository.
func NewStudentRepository() *StudentRepository {
	return &StudentRepository{db: config.DB}
}

// CreateStudent insere um novo aluno no banco de dados.
func (r *StudentRepository) CreateStudent(student *models.Student) error {
	// Corrige: importa corretamente o pacote uuid
	// Adicione "github.com/google/uuid" ao import se ainda não estiver lá
	student.ID = uuid.New().String() // Gera um ID único para o aluno
	query := `INSERT INTO students (id, enrollment, name, current_year) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, student.ID, student.Enrollment, student.Name, student.CurrentYear)
	if err != nil {
		log.Printf("Erro ao criar aluno: %v", err)
		return err
	}

	// Insere as matérias do aluno na tabela de relacionamento
	for _, subject := range student.Subjects {
		err := r.AddSubjectToStudent(student.ID, subject.ID)
		if err != nil {
			log.Printf("Aviso: Erro ao adicionar matéria %s ao aluno %s: %v", subject.ID, student.ID, err)
			// Decida se o erro aqui deve parar a criação do aluno ou apenas logar
			// Por enquanto, vamos apenas logar e continuar
		}
	}
	return nil
}

// GetStudentByID busca um aluno pelo ID, incluindo suas matérias.
func (r *StudentRepository) GetStudentByID(id string) (*models.Student, error) {
	student := &models.Student{}
	query := `SELECT id, enrollment, name, current_year FROM students WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&student.ID, &student.Enrollment, &student.Name, &student.CurrentYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Aluno não encontrado
		}
		log.Printf("Erro ao buscar aluno por ID: %v", err)
		return nil, err
	}

	// Busca as matérias associadas a este aluno
	subjects, err := r.GetSubjectsByStudentID(student.ID)
	if err != nil {
		log.Printf("Erro ao buscar matérias para o aluno %s: %v", student.ID, err)
		return nil, err
	}
	student.Subjects = subjects

	return student, nil
}

// GetAllStudents busca todos os alunos, incluindo suas matérias.
func (r *StudentRepository) GetAllStudents() ([]models.Student, error) {
	rows, err := r.db.Query(`SELECT id, enrollment, name, current_year FROM students`)
	if err != nil {
		log.Printf("Erro ao buscar todos os alunos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		student := models.Student{}
		if err := rows.Scan(&student.ID, &student.Enrollment, &student.Name, &student.CurrentYear); err != nil {
			log.Printf("Erro ao escanear aluno: %v", err)
			return nil, err
		}

		// Busca as matérias para cada aluno
		subjects, err := r.GetSubjectsByStudentID(student.ID)
		if err != nil {
			log.Printf("Erro ao buscar matérias para o aluno %s: %v", student.ID, err)
			return nil, err
		}
		student.Subjects = subjects
		students = append(students, student)
	}
	return students, nil
}

// UpdateStudent atualiza um aluno existente.
// Esta função não atualiza as matérias diretamente, apenas os dados do aluno.
// O gerenciamento de matérias (adicionar/remover) deve ser feito por funções separadas.
func (r *StudentRepository) UpdateStudent(student *models.Student) error {
	query := `UPDATE students SET enrollment = ?, name = ?, current_year = ? WHERE id = ?`
	result, err := r.db.Exec(query, student.Enrollment, student.Name, student.CurrentYear, student.ID)
	if err != nil {
		log.Printf("Erro ao atualizar aluno: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Nenhum aluno encontrado para atualizar
	}
	return nil
}

// DeleteStudent deleta um aluno pelo ID.
// O ON DELETE CASCADE na tabela student_subjects garante que os relacionamentos sejam deletados.
func (r *StudentRepository) DeleteStudent(id string) error {
	query := `DELETE FROM students WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Erro ao deletar aluno: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Nenhum aluno encontrado para deletar
	}
	return nil
}

// AddSubjectToStudent adiciona uma matéria a um aluno.
func (r *StudentRepository) AddSubjectToStudent(studentID, subjectID string) error {
	// Primeiro, verifica se a matéria existe para garantir integridade
	subjectRepo := NewSubjectRepository()
	existingSubject, err := subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		return err // Erro ao buscar a matéria
	}
	if existingSubject == nil {
		return sql.ErrNoRows // Matéria não existe
	}

	query := `INSERT OR IGNORE INTO student_subjects (student_id, subject_id) VALUES (?, ?)`
	_, err = r.db.Exec(query, studentID, subjectID)
	if err != nil {
		log.Printf("Erro ao adicionar matéria %s ao aluno %s: %v", subjectID, studentID, err)
		return err
	}
	return nil
}

// RemoveSubjectFromStudent remove uma matéria de um aluno.
func (r *StudentRepository) RemoveSubjectFromStudent(studentID, subjectID string) error {
	query := `DELETE FROM student_subjects WHERE student_id = ? AND subject_id = ?`
	result, err := r.db.Exec(query, studentID, subjectID)
	if err != nil {
		log.Printf("Erro ao remover matéria %s do aluno %s: %v", subjectID, studentID, err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Relacionamento não encontrado
	}
	return nil
}

// GetSubjectsByStudentID busca todas as matérias associadas a um aluno.
func (r *StudentRepository) GetSubjectsByStudentID(studentID string) ([]models.Subject, error) {
	query := `
    SELECT s.id, s.name, s.year, s.credits
    FROM subjects s
    JOIN student_subjects ss ON s.id = ss.subject_id
    WHERE ss.student_id = ?`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		log.Printf("Erro ao buscar matérias para o aluno %s: %v", studentID, err)
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{}
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits); err != nil {
			log.Printf("Erro ao escanear matéria do aluno: %v", err)
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}
