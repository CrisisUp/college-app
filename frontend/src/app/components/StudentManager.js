// frontend/src/app/components/StudentManager.js
'use client';

import { useState } from 'react'; // Removido useEffect pois o fetch inicial será no page.js

const API_BASE_URL = 'http://localhost:8080'; // A API_BASE_URL agora é gerenciada dentro de apiService.js

// Importa o serviço de aluno
import { studentService } from '../services/apiService';

// Componente StudentManager recebe students, subjects, styles e fetchStudents como props
export default function StudentManager({ students, subjects, styles, fetchStudents }) {
    // Estados para o formulário de cadastro de Alunos
    const [studentName, setStudentName] = useState('');
    const [currentYear, setCurrentYear] = useState('');
    const [studentShift, setStudentShift] = useState('M');
    const [studentMessage, setStudentMessage] = useState('');
    const [studentMessageType, setStudentMessageType] = useState('');

    // Estados para a funcionalidade de edição de Alunos
    const [editingStudent, setEditingStudent] = useState(null); // Armazena o aluno em edição
    const [editStudentName, setEditStudentName] = useState(''); // Nome no form de edição
    const [editCurrentYear, setEditCurrentYear] = useState(''); // Ano no form de edição
    const [editStudentShift, setEditStudentShift] = useState(''); // Turno no form de edição

    async function handleStudentSubmit(e) {
        e.preventDefault();
        setStudentMessage('');
        setStudentMessageType('');

        const firstYearSubjects = subjects.filter(sub => sub.year === 1);
        const shuffledSubjects = firstYearSubjects.sort(() => 0.5 - Math.random());
        const selectedSubjectIds = shuffledSubjects.slice(0, 5).map(sub => ({ id: sub.id }));

        if (selectedSubjectIds.length < 5 && firstYearSubjects.length > 0) {
            setStudentMessage('Aviso: Não há 5 matérias suficientes do primeiro ano para associar.');
            setStudentMessageType('error');
        } else if (firstYearSubjects.length === 0) {
             setStudentMessage('Aviso: Nenhuma matéria do primeiro ano disponível para associação automática.');
             setStudentMessageType('error');
        }

        try {
            // Usa o studentService para criar o aluno
            const newStudent = await studentService.create({
                name: studentName,
                current_year: parseInt(currentYear, 10),
                shift: studentShift,
                subjects: selectedSubjectIds,
            });

            console.log('Aluno criado:', newStudent);
            setStudentMessage(`Aluno "${newStudent.name}" cadastrado com sucesso! Matrícula: ${newStudent.enrollment}`);
            setStudentMessageType('success');
            
            setStudentName('');
            setCurrentYear('');
            setStudentShift('M');

            fetchStudents(); // Chama a função passada como prop para recarregar a lista globalmente
        } catch (error) {
            console.error('Erro ao cadastrar aluno:', error);
            setStudentMessage(`Erro ao cadastrar aluno: ${error.message}`);
            setStudentMessageType('error');
        }
    }

    async function handleDeleteStudent(id, name) {
        if (window.confirm(`Tem certeza que deseja deletar o aluno ${name}?`)) {
            try {
                // Usa o studentService para deletar o aluno
                await studentService.delete(id);

                setStudentMessage(`Aluno "${name}" deletado com sucesso!`);
                setStudentMessageType('success');
                fetchStudents(); // Recarrega a lista
            } catch (error) {
                console.error('Erro ao deletar aluno:', error);
                setStudentMessage(`Erro ao deletar aluno: ${error.message}`);
                setStudentMessageType('error');
            }
        }
    }

    // FUNÇÃO handleUpdateStudent (inicia a edição)
    function handleUpdateStudent(student) {
        setEditingStudent(student);
        setEditStudentName(student.name);
        setEditCurrentYear(student.current_year);
        setEditStudentShift(student.shift);
        setStudentMessage('');
    }

    // FUNÇÃO handleCancelEditStudent (cancela a edição)
    function handleCancelEditStudent() {
        setEditingStudent(null);
        setEditStudentName('');
        setEditCurrentYear('');
        setEditStudentShift('M');
        setStudentMessage('');
    }

    // FUNÇÃO handleUpdateStudentSubmit (envia a atualização)
    async function handleUpdateStudentSubmit(e) {
        e.preventDefault();
        setStudentMessage('');
        setStudentMessageType('');

        if (!editingStudent) return;

        try {
            // Usa o studentService para atualizar o aluno
            const updatedStudent = await studentService.update(editingStudent.id, {
                enrollment: editingStudent.enrollment,
                name: editStudentName,
                current_year: parseInt(editCurrentYear, 10),
                shift: editStudentShift,
                subjects: editingStudent.subjects?.map(sub => ({id: sub.id})) || [],
            });

            console.log('Aluno atualizado:', updatedStudent);
            setStudentMessage(`Aluno "${updatedStudent.name}" atualizado com sucesso!`);
            setStudentMessageType('success');
            setEditingStudent(null); // Sai do modo de edição
            fetchStudents(); // Recarrega a lista
        } catch (error) {
            console.error('Erro ao atualizar aluno:', error);
            setStudentMessage(`Erro ao atualizar aluno: ${error.message}`);
            setStudentMessageType('error');
        }
    }

    return (
        <>
            <h2 style={styles.h2}>Cadastrar Novo Aluno</h2>
            <form onSubmit={handleStudentSubmit} style={styles.form}>
                <div style={styles.inlineFormGroup}>
                    <input
                        type="text"
                        placeholder="Nome do Aluno"
                        value={studentName}
                        onChange={(e) => setStudentName(e.target.value)}
                        required
                        style={{ ...styles.input, flex: '2' }}
                    />
                    <input
                        type="number"
                        placeholder="Ano Atual"
                        value={currentYear}
                        onChange={(e) => setCurrentYear(e.target.value)}
                        required
                        min="1"
                        style={{ ...styles.input, flex: '1' }}
                    />
                    <select
                        value={studentShift}
                        onChange={(e) => setStudentShift(e.target.value)}
                        required
                        style={styles.input}
                    >
                        <option value="M">Manhã</option>
                        <option value="T">Tarde</option>
                        <option value="N">Noite</option>
                    </select>
                </div>
                <button type="submit" style={styles.button}>Cadastrar Aluno</button>
                {studentMessage && (
                    <div style={{ ...styles.message, color: studentMessageType === 'error' ? '#c0392b' : '#27ae60', backgroundColor: studentMessageType === 'error' ? '#fde0dc' : '#d4edda', border: `1px solid ${studentMessageType === 'error' ? '#e74c3c' : '#28a745'}` }}>
                        {studentMessage}
                    </div>
                )}
            </form>

            {/* Formulário de Edição de Aluno (condicional) */}
            {editingStudent && (
                <div style={styles.editFormContainer}>
                    <h2 style={styles.h2}>Editar Aluno: {editingStudent.name} ({editingStudent.enrollment})</h2>
                    <form onSubmit={handleUpdateStudentSubmit} style={styles.form}>
                        <input
                            type="text"
                            placeholder="Nome do Aluno"
                            value={editStudentName}
                            onChange={(e) => setEditStudentName(e.target.value)}
                            required
                            style={styles.input}
                        />
                        <input
                            type="number"
                            placeholder="Ano Atual"
                            value={editCurrentYear}
                            onChange={(e) => setEditCurrentYear(e.target.value)}
                            required
                            min="1"
                            style={styles.input}
                        />
                         <select
                            value={editStudentShift}
                            onChange={(e) => setEditStudentShift(e.target.value)}
                            required
                            style={styles.input}
                        >
                            <option value="M">Manhã</option>
                            <option value="T">Tarde</option>
                            <option value="N">Noite</option>
                        </select>
                        <div style={styles.formButtons}>
                            <button type="submit" style={{...styles.button, ...styles.updateButton}}>Salvar Alterações</button>
                            <button type="button" onClick={handleCancelEditStudent} style={{...styles.button, ...styles.cancelButton}}>Cancelar</button>
                        </div>
                        {studentMessage && (
                            <div style={{ ...styles.message, color: studentMessageType === 'error' ? '#c0392b' : '#27ae60', backgroundColor: studentMessageType === 'error' ? '#fde0dc' : '#d4edda', border: `1px solid ${studentMessageType === 'error' ? '#e74c3c' : '#28a745'}` }}>
                                {studentMessage}
                            </div>
                        )}
                    </form>
                </div>
            )}

            <h2 style={styles.h2}>Lista de Alunos</h2>
            <ul style={styles.ul}>
                {Array.isArray(students) && students.length > 0 ? (
                    students.map((student) => (
                        <li key={student.id} style={styles.li}>
                            <strong>ID:</strong> {student.id}<br />
                            <strong>Matrícula:</strong> {student.enrollment}<br />
                            <strong>Nome:</strong> {student.name}<br />
                            <strong>Ano Atual:</strong> {student.current_year}<br />
                            <strong>Turno:</strong> {student.shift}<br />
                            <strong>Matérias:</strong> {
                                student.subjects?.length > 0
                                    ? student.subjects.map(sub => `${sub.name} (${sub.id})`).join(', ')
                                    : 'Nenhuma matéria associada.'
                            }
                            <div style={styles.cardButtons}>
                                <button onClick={() => handleUpdateStudent(student)} style={{...styles.button, ...styles.updateButton}}>Atualizar</button>
                                <button onClick={() => handleDeleteStudent(student.id, student.name)} style={{...styles.button, ...styles.deleteButton}}>Deletar</button>
                            </div>
                        </li>
                    ))
                ) : (
                    <li style={styles.li}>Nenhum aluno cadastrado ainda.</li>
                )}
            </ul>
        </>
    );
}