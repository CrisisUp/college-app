// api/services/teacher_service.go
package services

import (
	"college_api/models"
	"college_api/repositories"
	"errors"
	"fmt"
	"log"     // Importar para log.Printf
	"strconv" // Importar para strconv.Atoi
	"strings" // Importar para strings.ToUpper, strings.ReplaceAll

	"github.com/google/uuid" // Importar para uuid.New().String()
)

// TeacherService define a interface para as operações de serviço de professores.
type TeacherService struct {
	repo *repositories.TeacherRepository
}

// NewTeacherService cria uma nova instância de TeacherService.
func NewTeacherService(repo *repositories.TeacherRepository) *TeacherService {
	return &TeacherService{repo: repo}
}

// CreateTeacher adiciona um novo professor com registro gerado automaticamente.
func (s *TeacherService) CreateTeacher(teacher *models.Teacher) error {
	// 1. Validação de campos essenciais do frontend
	if teacher.Name == "" || teacher.Department == "" {
		return errors.New("nome e departamento do professor são obrigatórios")
	}

	// --- NOVO: Gerar o ID único do professor (interno) ---
	teacher.ID = uuid.New().String() // Gera o UUID aqui

	// --- NOVO: Lógica para gerar o Registro do professor ---
	// 2. Gerar o código do departamento padronizado
	departmentCode := strings.ToUpper(strings.ReplaceAll(teacher.Department, " ", ""))
	if len(departmentCode) > 4 { // Limita o código para algo razoável
		departmentCode = departmentCode[:4]
	} else if len(departmentCode) == 0 {
		return errors.New("departamento não pode ser vazio para gerar registro")
	}

	// 3. Buscar o último registro para este departamento
	lastRegistry, err := s.repo.GetLastRegistryForDepartment(departmentCode)
	if err != nil {
		return fmt.Errorf("erro ao buscar último registro para o departamento: %w", err)
	}

	// 4. Gerar o novo número sequencial
	newSequence := 1
	if lastRegistry != "" {
		// Espera formato CODIGO-NNN (ex: COMP-001)
		// Extrai a parte numérica do registro (os últimos 3 dígitos)
		parts := strings.Split(lastRegistry, "-")
		if len(parts) > 1 {
			seqStr := parts[len(parts)-1] // Pega a última parte (o número)
			lastSequence, err := strconv.Atoi(seqStr)
			if err == nil {
				newSequence = lastSequence + 1
			} else {
				log.Printf("Aviso: Não foi possível converter sequência '%s' para int. Reiniciando sequência para 1. Erro: %v", seqStr, err)
			}
		}
	}

	// 5. Formata o novo registro (ex: COMP-001)
	teacher.Registry = fmt.Sprintf("%s-%03d", departmentCode, newSequence)

	// O repositório agora salvará o professor com o ID e Registro gerados
	return s.repo.CreateTeacher(teacher)
}

// GetTeacherByID busca um professor pelo ID.
func (s *TeacherService) GetTeacherByID(id string) (*models.Teacher, error) {
	teacher, err := s.repo.GetTeacherByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar professor: %w", err)
	}
	if teacher == nil {
		return nil, errors.New("professor não encontrado")
	}
	return teacher, nil
}

// GetAllTeachers busca todos os professores.
func (s *TeacherService) GetAllTeachers() ([]models.Teacher, error) {
	teachers, err := s.repo.GetAllTeachers()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos os professores: %w", err)
	}
	return teachers, nil
}

// UpdateTeacher atualiza um professor existente após validações.
func (s *TeacherService) UpdateTeacher(teacher *models.Teacher) error {
	if teacher.ID == "" {
		return errors.New("ID do professor é obrigatório para atualização")
	}
	if teacher.Name == "" || teacher.Department == "" {
		return errors.New("nome e departamento do professor são obrigatórios para atualização")
	}

	existingTeacher, err := s.repo.GetTeacherByID(teacher.ID)
	if err != nil {
		return fmt.Errorf("erro ao verificar professor para atualização: %w", err)
	}
	if existingTeacher == nil {
		return errors.New("professor não encontrado para atualização")
	}

	// Atualiza apenas os campos permitidos (nome e departamento)
	existingTeacher.Name = teacher.Name
	existingTeacher.Department = teacher.Department
	// O registro (Registry) não é atualizado por aqui, pois é gerado na criação.
	// Se Registry precisar ser atualizado, seria um método de negócio específico.
	// Garanta que o Registry original seja mantido (não substituído por vazio)
	teacher.Registry = existingTeacher.Registry // Atribui o registro existente ao professor no DTO de entrada para que o repositório não o apague

	return s.repo.UpdateTeacher(existingTeacher)
}

// DeleteTeacher deleta um professor pelo ID.
func (s *TeacherService) DeleteTeacher(id string) error {
	if id == "" {
		return errors.New("ID do professor é obrigatório para exclusão")
	}
	existingTeacher, err := s.repo.GetTeacherByID(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar professor para exclusão: %w", err)
	}
	if existingTeacher == nil {
		return errors.New("professor não encontrado para exclusão")
	}

	return s.repo.DeleteTeacher(id)
}
