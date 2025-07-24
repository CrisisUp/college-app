// services/student_service.go
package services

import (
	"college_api/models"
	"college_api/repositories"
	"errors"
	"fmt"
)

// StudentService define a interface para as operações de serviço de alunos.
type StudentService struct {
	studentRepo *repositories.StudentRepository
	subjectRepo *repositories.SubjectRepository // Para validar matérias
}

// NewStudentService cria uma nova instância de StudentService.
func NewStudentService(sr *repositories.StudentRepository, subjR *repositories.SubjectRepository) *StudentService {
	return &StudentService{
		studentRepo: sr,
		subjectRepo: subjR,
	}
}

// CreateStudent adiciona um novo aluno após validações.
func (s *StudentService) CreateStudent(student *models.Student) error {
	// Validação: Matrícula, nome e ano atual são obrigatórios
	if student.Enrollment == "" || student.Name == "" || student.CurrentYear == 0 {
		return errors.New("matrícula, nome e ano atual do aluno são obrigatórios")
	}

	// Validação: Matrícula já existe (assumindo que seja UNIQUE)
	// Para isso, precisaríamos de um método GetStudentByEnrollment no repositório.
	// Por simplicidade, vamos pular essa validação por enquanto ou assumir que o DB cuida dela.
	// Se a matrícula for UNIQUE no DB, o s.studentRepo.CreateStudent retornará erro.

	// Validação: Verificar se todas as matérias existem
	for i, sub := range student.Subjects {
		if sub.ID == "" {
			return errors.New("ID da matéria não pode ser vazio")
		}
		foundSub, err := s.subjectRepo.GetSubjectByID(sub.ID)
		if err != nil {
			return fmt.Errorf("erro ao verificar matéria %s: %w", sub.ID, err)
		}
		if foundSub == nil {
			return fmt.Errorf("matéria com ID %s não encontrada", sub.ID)
		}
		// Atualiza a matéria no modelo do aluno com os dados completos do DB
		student.Subjects[i] = *foundSub
	}

	return s.studentRepo.CreateStudent(student)
}

// GetStudentByID busca um aluno pelo ID.
func (s *StudentService) GetStudentByID(id string) (*models.Student, error) {
	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aluno: %w", err)
	}
	if student == nil {
		return nil, errors.New("aluno não encontrado")
	}
	return student, nil
}

// GetAllStudents busca todos os alunos.
func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	students, err := s.studentRepo.GetAllStudents()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos os alunos: %w", err)
	}
	return students, nil
}

// UpdateStudent atualiza um aluno existente após validações.
func (s *StudentService) UpdateStudent(student *models.Student) error {
	if student.ID == "" {
		return errors.New("ID do aluno é obrigatório para atualização")
	}
	if student.Enrollment == "" || student.Name == "" || student.CurrentYear == 0 {
		return errors.New("matrícula, nome e ano atual do aluno são obrigatórios para atualização")
	}

	// Validação: o aluno deve existir para ser atualizado
	existingStudent, err := s.studentRepo.GetStudentByID(student.ID)
	if err != nil {
		return fmt.Errorf("erro ao verificar aluno para atualização: %w", err)
	}
	if existingStudent == nil {
		return errors.New("aluno não encontrado para atualização")
	}

	// A atualização de matérias associadas não é feita aqui diretamente,
	// mas sim pelas funções AddSubjectToStudent/RemoveSubjectFromStudent.
	// O Slice student.Subjects que vem no DTO de update será ignorado aqui
	// e precisa de handlers e serviços separados para ser modificado.

	return s.studentRepo.UpdateStudent(student)
}

// DeleteStudent deleta um aluno pelo ID.
func (s *StudentService) DeleteStudent(id string) error {
	if id == "" {
		return errors.New("ID do aluno é obrigatório para exclusão")
	}
	// Validação: o aluno deve existir para ser deletado
	existingStudent, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar aluno para exclusão: %w", err)
	}
	if existingStudent == nil {
		return errors.New("aluno não encontrado para exclusão")
	}

	return s.studentRepo.DeleteStudent(id)
}

// AddSubjectToStudent adiciona uma matéria a um aluno, com validações.
func (s *StudentService) AddSubjectToStudent(studentID, subjectID string) error {
	if studentID == "" || subjectID == "" {
		return errors.New("IDs de aluno e matéria são obrigatórios")
	}

	// Valida se o aluno existe
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return fmt.Errorf("erro ao verificar aluno: %w", err)
	}
	if student == nil {
		return errors.New("aluno não encontrado")
	}

	// Valida se a matéria existe
	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		return fmt.Errorf("erro ao verificar matéria: %w", err)
	}
	if subject == nil {
		return errors.New("matéria não encontrada")
	}

	// Verifica se a matéria já está associada ao aluno
	for _, sub := range student.Subjects {
		if sub.ID == subjectID {
			return errors.New("matéria já associada a este aluno")
		}
	}

	return s.studentRepo.AddSubjectToStudent(studentID, subjectID)
}

// RemoveSubjectFromStudent remove uma matéria de um aluno, com validações.
func (s *StudentService) RemoveSubjectFromStudent(studentID, subjectID string) error {
	if studentID == "" || subjectID == "" {
		return errors.New("IDs de aluno e matéria são obrigatórios")
	}

	// Valida se o aluno existe
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return fmt.Errorf("erro ao verificar aluno: %w", err)
	}
	if student == nil {
		return errors.New("aluno não encontrado")
	}

	// Não é estritamente necessário validar se a matéria existe aqui,
	// pois o repositório cuidará de não remover o que não existe.
	// Mas podemos validar se a associação existe para dar um erro mais específico.
	found := false
	for _, sub := range student.Subjects {
		if sub.ID == subjectID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("matéria não associada a este aluno")
	}

	return s.studentRepo.RemoveSubjectFromStudent(studentID, subjectID)
}
