// frontend/src/app/services/apiService.js
const API_BASE_URL = 'http://localhost:8080';

// Função auxiliar para lidar com as respostas da API
async function handleResponse(response) {
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Erro desconhecido da API' }));
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorData.message || response.statusText}`);
    }
    return response.json();
}

// --- Funções de Serviço para Alunos ---
export const studentService = {
    getAll: async () => {
        const response = await fetch(`${API_BASE_URL}/students`);
        return handleResponse(response);
    },
    create: async (studentData) => {
        const response = await fetch(`${API_BASE_URL}/students`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(studentData),
        });
        return handleResponse(response);
    },
    update: async (id, studentData) => {
        const response = await fetch(`${API_BASE_URL}/students/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(studentData),
        });
        return handleResponse(response);
    },
    delete: async (id) => {
        const response = await fetch(`${API_BASE_URL}/students/${id}`, {
            method: 'DELETE',
        });
        if (!response.ok) { // DELETE 204 No Content não tem body
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return null; // Retorna null ou um indicador de sucesso, já que 204 não tem JSON
    },
    // Funções para associação de matérias
    addSubject: async (studentId, subjectId) => {
        const response = await fetch(`${API_BASE_URL}/students/${studentId}/subjects/${subjectId}`, {
            method: 'POST',
        });
        return handleResponse(response);
    },
    removeSubject: async (studentId, subjectId) => {
        const response = await fetch(`${API_BASE_URL}/students/${studentId}/subjects/${subjectId}`, {
            method: 'DELETE',
        });
        if (!response.ok) { // DELETE 204 No Content não tem body
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return null;
    }
};

// --- Funções de Serviço para Professores ---
export const teacherService = {
    getAll: async () => {
        const response = await fetch(`${API_BASE_URL}/teachers`);
        return handleResponse(response);
    },
    create: async (teacherData) => {
        const response = await fetch(`${API_BASE_URL}/teachers`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(teacherData),
        });
        return handleResponse(response);
    },
    update: async (id, teacherData) => {
        const response = await fetch(`${API_BASE_URL}/teachers/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(teacherData),
        });
        return handleResponse(response);
    },
    delete: async (id) => {
        const response = await fetch(`${API_BASE_URL}/teachers/${id}`, {
            method: 'DELETE',
        });
        if (!response.ok) { // DELETE 204 No Content não tem body
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return null;
    }
};

// --- Funções de Serviço para Matérias ---
export const subjectService = {
    getAll: async () => {
        const response = await fetch(`${API_BASE_URL}/subjects`);
        return handleResponse(response);
    }
    // Adicione create, update, delete se for implementar no frontend
};