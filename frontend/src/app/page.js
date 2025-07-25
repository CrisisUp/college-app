// frontend/src/app/page.js
'use client';

import { useState, useEffect } from 'react';
// Importa os componentes filhos
import StudentManager from './components/StudentManager';
import TeacherManager from './components/TeacherManager';
import AssociateSubjectForm from './components/AssociateSubjectForm';
// Importa os novos serviços de API
import { studentService, teacherService, subjectService } from './services/apiService';


// A API_BASE_URL agora é gerenciada dentro de apiService.js

export default function Home() {
    // Estados globais que serão passados como props para os componentes filhos
    const [students, setStudents] = useState([]);
    const [teachers, setTeachers] = useState([]);
    const [subjects, setSubjects] = useState([]);

    // Funções de busca globais que usam os novos serviços de API
    async function fetchStudentsGlobal() { // Renomeado para evitar conflito com StudentManager.js
        try {
            const data = await studentService.getAll();
            setStudents(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error('Erro global ao buscar alunos:', error);
            setStudents([]);
        }
    }

    async function fetchTeachersGlobal() { // Renomeado
        try {
            const data = await teacherService.getAll();
            setTeachers(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error('Erro global ao buscar professores:', error);
            setTeachers([]);
        }
    }

    async function fetchSubjectsGlobal() { // Renomeado
        try {
            const data = await subjectService.getAll();
            setSubjects(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error('Erro global ao buscar matérias:', error);
            setSubjects([]);
        }
    }

    // Efeito para buscar todos os dados iniciais quando o componente Home é montado
    useEffect(() => {
        fetchStudentsGlobal();
        fetchTeachersGlobal();
        fetchSubjectsGlobal();
    }, []); // Array de dependências vazio para rodar apenas uma vez na montagem

    return (
        <main style={styles.main}>
            <div style={styles.container}>
                <h1 style={styles.h1}>Gerenciador Universitário</h1>

                {/* Renderiza o componente StudentManager, passando dados e funções */}
                <StudentManager
                    students={students}
                    subjects={subjects}
                    styles={styles}
                    fetchStudents={fetchStudentsGlobal} // Passa a função global de recarregamento
                    studentService={studentService} // Passa o serviço de aluno
                />
                <hr style={styles.hr} />

                {/* Renderiza o componente TeacherManager, passando dados e funções */}
                <TeacherManager
                    teachers={teachers}
                    styles={styles}
                    fetchTeachers={fetchTeachersGlobal} // Passa a função global de recarregamento
                    teacherService={teacherService} // Passa o serviço de professor
                />
                <hr style={styles.hr} />

                {/* Renderiza o componente AssociateSubjectForm, passando dados e callback */}
                <AssociateSubjectForm
                    students={students}
                    subjects={subjects}
                    styles={styles}
                    onAssociateSuccess={fetchStudentsGlobal} // Callback para recarregar alunos
                    studentService={studentService} // Passa o serviço de aluno (para associar/desassociar)
                    subjectService={subjectService} // Passa o serviço de matéria (se fosse necessário mais que getAll)
                />

            </div>
        </main>
    );
}

// Estilos básicos para o componente (MANTÉM TODOS OS ESTILOS AQUI)
const styles = {
    main: {
        fontFamily: 'Arial, sans-serif',
        margin: '0',
        backgroundColor: '#e0f2f7',
        minHeight: '100vh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'flex-start',
        padding: '20px 0',
        color: '#333',
    },
    container: {
        backgroundColor: 'white',
        padding: '30px',
        borderRadius: '10px',
        boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
        maxWidth: '800px',
        width: '100%',
    },
    h1: { color: '#2c3e50', textAlign: 'center', marginBottom: '30px' },
    h2: { color: '#34495e', marginTop: '25px', marginBottom: '15px', borderBottom: '1px solid #eee', paddingBottom: '5px' },
    ul: { listStyle: 'none', padding: '0' },
    li: {
        backgroundColor: '#f9f9f9',
        marginBottom: '10px',
        padding: '15px',
        borderRadius: '8px',
        border: '1px solid #ddd',
        boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
        position: 'relative',
        paddingBottom: '40px',
    },
    form: { display: 'flex', flexDirection: 'column', gap: '10px', marginBottom: '20px' },
    input: { padding: '10px', border: '1px solid #a0c4ff', borderRadius: '5px', fontSize: '16px' },
    button: { padding: '12px 18px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer', fontSize: '16px', fontWeight: 'bold' },
    message: { marginTop: '15px', padding: '12px', borderRadius: '5px', textAlign: 'center', fontWeight: 'bold' },
    inlineFormGroup: {
        display: 'flex',
        gap: '10px',
        marginBottom: '10px',
        alignItems: 'center',
        flexWrap: 'wrap',
    },
    hr: { border: '0', height: '1px', backgroundColor: '#ccc', margin: '40px 0' },
    cardButtons: {
        position: 'absolute',
        bottom: '10px',
        right: '10px',
        display: 'flex',
        gap: '8px',
    },
    updateButton: {
        backgroundColor: '#28a745',
        padding: '8px 12px',
        fontSize: '14px',
    },
    deleteButton: {
        backgroundColor: '#dc3545',
        padding: '8px 12px',
        fontSize: '14px',
    },
    editFormContainer: {
        backgroundColor: '#f0f8ff',
        padding: '20px',
        borderRadius: '8px',
        border: '1px dashed #a0c4ff',
        marginBottom: '20px',
    },
    formButtons: {
        display: 'flex',
        gap: '10px',
        justifyContent: 'flex-end',
        marginTop: '10px',
    },
    cancelButton: {
        backgroundColor: '#6c757d',
    }
};