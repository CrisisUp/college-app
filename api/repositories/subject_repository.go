package repositories

import (
	"college_api/config"
	"college_api/models"
	"database/sql"
	"log"
)

// SubjectRepository define as operações de CRUD para matérias.
type SubjectRepository struct {
	db *sql.DB
}

// NewSubjectRepository cria uma nova instância de SubjectRepository.
func NewSubjectRepository() *SubjectRepository {
	return &SubjectRepository{db: config.DB}
}

// CreateSubject insere uma nova matéria no banco de dados.
func (r *SubjectRepository) CreateSubject(subject *models.Subject) error {
	query := `INSERT INTO subjects (id, name, year, credits) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, subject.ID, subject.Name, subject.Year, subject.Credits)
	if err != nil {
		log.Printf("Erro ao criar matéria: %v", err)
		return err
	}
	return nil
}

// GetSubjectByID busca uma matéria pelo ID.
func (r *SubjectRepository) GetSubjectByID(id string) (*models.Subject, error) {
	subject := &models.Subject{}
	query := `SELECT id, name, year, credits FROM subjects WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Matéria não encontrada
		}
		log.Printf("Erro ao buscar matéria por ID: %v", err)
		return nil, err
	}
	return subject, nil
}

// GetAllSubjects busca todas as matérias.
func (r *SubjectRepository) GetAllSubjects() ([]models.Subject, error) {
	rows, err := r.db.Query(`SELECT id, name, year, credits FROM subjects`)
	if err != nil {
		log.Printf("Erro ao buscar todas as matérias: %v", err)
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{}
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits); err != nil {
			log.Printf("Erro ao escanear matéria: %v", err)
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

// UpdateSubject atualiza uma matéria existente.
func (r *SubjectRepository) UpdateSubject(subject *models.Subject) error {
	query := `UPDATE subjects SET name = ?, year = ?, credits = ? WHERE id = ?`
	result, err := r.db.Exec(query, subject.Name, subject.Year, subject.Credits, subject.ID)
	if err != nil {
		log.Printf("Erro ao atualizar matéria: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Nenhuma matéria encontrada para atualizar
	}
	return nil
}

// DeleteSubject deleta uma matéria pelo ID.
func (r *SubjectRepository) DeleteSubject(id string) error {
	query := `DELETE FROM subjects WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Erro ao deletar matéria: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Nenhuma matéria encontrada para deletar
	}
	return nil
}
