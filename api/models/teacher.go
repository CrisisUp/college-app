// models/teacher.go
package models

// Teacher representa um professor na universidade.
type Teacher struct {
	ID         string `json:"id"`         // ID único do professor (gerado, ex: UUID)
	Registry   string `json:"registry"`   // Registro único do professor (ex: "PROF001")
	Name       string `json:"name"`       // Nome completo do professor
	Department string `json:"department"` // Departamento do professor (ex: "Ciência da Computação")
}
