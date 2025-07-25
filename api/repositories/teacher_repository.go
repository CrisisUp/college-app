// api/repositories/teacher_repository.go
package repositories

import (
	"college_api/config"
	"college_api/models"
	"database/sql" // Adicionar import para fmt
	"log"
	// Não precisa importar uuid aqui se o serviço já gera o ID
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository() *TeacherRepository {
	return &TeacherRepository{db: config.DB}
}

// CreateTeacher insere um novo professor no banco de dados.
// O ID e Registry já devem vir preenchidos do Service.
func (r *TeacherRepository) CreateTeacher(teacher *models.Teacher) error {
	query := `INSERT INTO teachers (id, registry, name, department) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, teacher.ID, teacher.Registry, teacher.Name, teacher.Department)
	if err != nil {
		log.Printf("Erro ao criar professor: %v", err)
		return err
	}
	return nil
}

// GetTeacherByID busca um professor pelo ID.
func (r *TeacherRepository) GetTeacherByID(id string) (*models.Teacher, error) {
	teacher := &models.Teacher{}
	query := `SELECT id, registry, name, department FROM teachers WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.Registry, &teacher.Name, &teacher.Department)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Erro ao buscar professor por ID: %v", err)
		return nil, err
	}
	return teacher, nil
}

// GetAllTeachers busca todos os professores.
func (r *TeacherRepository) GetAllTeachers() ([]models.Teacher, error) {
	rows, err := r.db.Query(`SELECT id, registry, name, department FROM teachers`)
	if err != nil {
		log.Printf("Erro ao buscar todos os professores: %v", err)
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		teacher := models.Teacher{}
		if err := rows.Scan(&teacher.ID, &teacher.Registry, &teacher.Name, &teacher.Department); err != nil {
			log.Printf("Erro ao escanear professor: %v", err)
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

// UpdateTeacher atualiza um professor existente.
func (r *TeacherRepository) UpdateTeacher(teacher *models.Teacher) error {
	query := `UPDATE teachers SET registry = $1, name = $2, department = $3 WHERE id = $4`
	result, err := r.db.Exec(query, teacher.Registry, teacher.Name, teacher.Department, teacher.ID)
	if err != nil {
		log.Printf("Erro ao atualizar professor: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteTeacher deleta um professor pelo ID.
func (r *TeacherRepository) DeleteTeacher(id string) error {
	query := `DELETE FROM teachers WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Erro ao deletar professor: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetLastRegistryForDepartment busca o maior número de registro para o departamento especificado.
// Retorna o registro como string e um erro, se houver.
// Retorna "" e nil se não houver registros para o departamento.
func (r *TeacherRepository) GetLastRegistryForDepartment(departmentCode string) (string, error) {
	var lastRegistry sql.NullString
	query := `
		SELECT registry FROM teachers
		WHERE registry LIKE $1 || '-%'
		ORDER BY registry DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query, departmentCode).Scan(&lastRegistry)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Nenhuma matrícula encontrada para este ano e turno
		}
		log.Printf("Erro ao buscar último registro para o departamento %s: %v", departmentCode, err)
		return "", err
	}

	if lastRegistry.Valid {
		return lastRegistry.String, nil
	}
	return "", nil
}
