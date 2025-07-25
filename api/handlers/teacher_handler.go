// handlers/teacher_handler.go
package handlers

import (
	"college_api/models"
	"college_api/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TeacherHandler gerencia as requisições HTTP para professores.
type TeacherHandler struct {
	service *services.TeacherService
}

// NewTeacherHandler cria uma nova instância de TeacherHandler.
func NewTeacherHandler(s *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: s}
}

// CreateTeacherHandler lida com a criação de um novo professor.
// POST /teachers
func (h *TeacherHandler) CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTeacher(&teacher); err != nil {
		log.Printf("Erro ao criar professor no serviço: %v", err)
		http.Error(w, "Erro ao criar professor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(teacher)
}

// GetTeacherByIDHandler lida com a busca de um professor por ID.
// GET /teachers/{id}
func (h *TeacherHandler) GetTeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	teacher, err := h.service.GetTeacherByID(id)
	if err != nil {
		if err.Error() == "professor não encontrado" { // Erro personalizado do serviço
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao buscar professor no serviço: %v", err)
		http.Error(w, "Erro ao buscar professor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// GetAllTeachersHandler lida com a busca de todos os professores.
// GET /teachers
func (h *TeacherHandler) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.service.GetAllTeachers()
	if err != nil {
		log.Printf("Erro ao buscar todos os professores no serviço: %v", err)
		http.Error(w, "Erro ao buscar professores: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

// UpdateTeacherHandler lida com a atualização de um professor existente.
// PUT /teachers/{id}
func (h *TeacherHandler) UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	teacher.ID = id // Garante que o ID da URL seja usado

	if err := h.service.UpdateTeacher(&teacher); err != nil {
		if err.Error() == "professor não encontrado para atualização" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao atualizar professor no serviço: %v", err)
		http.Error(w, "Erro ao atualizar professor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teacher)
}

// DeleteTeacherHandler lida com a exclusão de um professor por ID.
// DELETE /teachers/{id}
func (h *TeacherHandler) DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteTeacher(id); err != nil {
		if err.Error() == "professor não encontrado para exclusão" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao deletar professor: %v", err)
		http.Error(w, "Erro ao deletar professor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204 No Content para deleção bem-sucedida
}
