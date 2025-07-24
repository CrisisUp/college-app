// models/student.go
package models

// Student representa um aluno na universidade.
type Student struct {
	ID          string    `json:"id"`           // ID único do aluno (gerado, ex: UUID)
	Enrollment  string    `json:"enrollment"`   // Matrícula do aluno (ex: "20230001")
	Name        string    `json:"name"`         // Nome completo do aluno
	CurrentYear int       `json:"current_year"` // Ano atual do aluno na universidade (ex: 1, 2, 3, 4)
	Subjects    []Subject `json:"subjects"`     // Matérias que o aluno está cursando/cursou
}
