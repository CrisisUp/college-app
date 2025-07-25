// frontend/src/app/components/TeacherManager.js
'use client';

import { useState } from 'react'; // Removido useEffect pois o fetch inicial será no page.js

// Importa o serviço de professor
import { teacherService } from '../services/apiService';

export default function TeacherManager({ styles, teachers, fetchTeachers }) {
    // Estados para o formulário de cadastro de Professores
    const [teacherName, setTeacherName] = useState('');
    const [teacherDepartment, setTeacherDepartment] = useState('');
    const [teacherMessage, setTeacherMessage] = useState('');
    const [teacherMessageType, setTeacherMessageType] = useState('');

    // Estados para a funcionalidade de edição de Professores
    const [editingTeacher, setEditingTeacher] = useState(null);
    const [editTeacherName, setEditTeacherName] = useState('');
    const [editTeacherDepartment, setEditTeacherDepartment] = useState('');

    // As funções fetchTeachers são passadas como prop do page.js,
    // então não precisamos de useEffect aqui para buscar a lista inicial.

    async function handleTeacherSubmit(e) {
        e.preventDefault();
        setTeacherMessage('');
        setTeacherMessageType('');

        try {
            // Usa o teacherService para criar o professor
            const newTeacher = await teacherService.create({
                name: teacherName,
                department: teacherDepartment,
            });

            console.log('Professor criado:', newTeacher);
            setTeacherMessage(`Professor "${newTeacher.name}" cadastrado com sucesso! Registro: ${newTeacher.registry}`);
            setTeacherMessageType('success');

            setTeacherName('');
            setTeacherDepartment('');
            fetchTeachers(); // Chama a função passada como prop para recarregar a lista globalmente
        } catch (error) {
            console.error('Erro ao cadastrar professor:', error);
            setTeacherMessage(`Erro ao cadastrar professor: ${error.message}`);
            setTeacherMessageType('error');
        }
    }

    async function handleDeleteTeacher(id, name) {
        if (window.confirm(`Tem certeza que deseja deletar o professor ${name}?`)) {
            try {
                // Usa o teacherService para deletar o professor
                await teacherService.delete(id);

                setTeacherMessage(`Professor "${name}" deletado com sucesso!`);
                setTeacherMessageType('success');
                fetchTeachers(); // Recarrega a lista
            } catch (error) {
                console.error('Erro ao deletar professor:', error);
                setTeacherMessage(`Erro ao deletar professor: ${error.message}`);
                setTeacherMessageType('error');
            }
        }
    }

    // FUNÇÃO handleUpdateTeacher (inicia a edição)
    function handleUpdateTeacher(teacher) {
        setEditingTeacher(teacher);
        setEditTeacherName(teacher.name);
        setEditTeacherDepartment(teacher.department);
        setTeacherMessage('');
    }

    // FUNÇÃO handleCancelEditTeacher (cancela a edição)
    function handleCancelEditTeacher() {
        setEditingTeacher(null);
        setEditTeacherName('');
        setEditTeacherDepartment('');
        setTeacherMessage('');
    }

    // FUNÇÃO handleUpdateTeacherSubmit (envia a atualização)
    async function handleUpdateTeacherSubmit(e) {
        e.preventDefault();
        setTeacherMessage('');
        setTeacherMessageType('');

        if (!editingTeacher) return;

        try {
            // Usa o teacherService para atualizar o professor
            const updatedTeacher = await teacherService.update(editingTeacher.id, {
                registry: editingTeacher.registry, // Manter o registro original
                name: editTeacherName,
                department: editTeacherDepartment,
            });

            console.log('Professor atualizado:', updatedTeacher);
            setTeacherMessage(`Professor "${updatedTeacher.name}" atualizado com sucesso!`);
            setTeacherMessageType('success');
            setEditingTeacher(null); // Sai do modo de edição
            fetchTeachers(); // Recarrega a lista
        } catch (error) {
            console.error('Erro ao atualizar professor:', error);
            setTeacherMessage(`Erro ao atualizar professor: ${error.message}`);
            setTeacherMessageType('error');
        }
    }

    return (
        <>
            <h2 style={styles.h2}>Cadastrar Novo Professor</h2>
            <form onSubmit={handleTeacherSubmit} style={styles.form}>
                <input
                    type="text"
                    placeholder="Nome do Professor"
                    value={teacherName}
                    onChange={(e) => setTeacherName(e.target.value)}
                    required
                    style={styles.input}
                />
                <input
                    type="text"
                    placeholder="Departamento (ex: Computação)"
                    value={teacherDepartment}
                    onChange={(e) => setTeacherDepartment(e.target.value)}
                    required
                    style={styles.input}
                />
                <button type="submit" style={styles.button}>Cadastrar Professor</button>
                {teacherMessage && (
                    <div style={{ ...styles.message, color: teacherMessageType === 'error' ? '#c0392b' : '#27ae60', backgroundColor: teacherMessageType === 'error' ? '#fde0dc' : '#d4edda', border: `1px solid ${teacherMessageType === 'error' ? '#e74c3c' : '#28a745'}` }}>
                        {teacherMessage}
                    </div>
                )}
            </form>

            {editingTeacher && (
                <div style={styles.editFormContainer}>
                    <h2 style={styles.h2}>Editar Professor: {editingTeacher.name} ({editingTeacher.registry})</h2>
                    <form onSubmit={handleUpdateTeacherSubmit} style={styles.form}>
                        <input
                            type="text"
                            placeholder="Nome do Professor"
                            value={editTeacherName}
                            onChange={(e) => setEditTeacherName(e.target.value)}
                            required
                            style={styles.input}
                        />
                        <input
                            type="text"
                            placeholder="Departamento"
                            value={editTeacherDepartment}
                            onChange={(e) => setEditTeacherDepartment(e.target.value)}
                            required
                            style={styles.input}
                        />
                        <div style={styles.formButtons}>
                            <button type="submit" style={{...styles.button, ...styles.updateButton}}>Salvar Alterações</button>
                            <button type="button" onClick={handleCancelEditTeacher} style={{...styles.button, ...styles.cancelButton}}>Cancelar</button>
                        </div>
                        {teacherMessage && (
                            <div style={{ ...styles.message, color: teacherMessageType === 'error' ? '#c0392b' : '#27ae60', backgroundColor: teacherMessageType === 'error' ? '#fde0dc' : '#d4edda', border: `1px solid ${teacherMessageType === 'error' ? '#e74c3c' : '#28a745'}` }}>
                                {teacherMessage}
                            </div>
                        )}
                    </form>
                </div>
            )}

            <h2 style={styles.h2}>Lista de Professores</h2>
            <ul style={styles.ul}>
                {Array.isArray(teachers) && teachers.length > 0 ? (
                    teachers.map((teacher) => (
                        <li key={teacher.id} style={styles.li}>
                            <strong>ID:</strong> {teacher.id}<br />
                            <strong>Registro:</strong> {teacher.registry}<br />
                            <strong>Nome:</strong> {teacher.name}<br />
                            <strong>Departamento:</strong> {teacher.department}
                            <div style={styles.cardButtons}>
                                <button onClick={() => handleUpdateTeacher(teacher)} style={{...styles.button, ...styles.updateButton}}>Atualizar</button>
                                <button onClick={() => handleDeleteTeacher(teacher.id, teacher.name)} style={{...styles.button, ...styles.deleteButton}}>Deletar</button>
                            </div>
                        </li>
                    ))
                ) : (
                    <li style={styles.li}>Nenhum professor cadastrado ainda.</li>
                )}
            </ul>
        </>
    );
}