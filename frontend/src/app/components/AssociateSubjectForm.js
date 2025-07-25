// frontend/src/app/components/AssociateSubjectForm.js
'use client';

import { useState, useEffect } from 'react';

// Importa os serviços de API
import { studentService, subjectService } from '../services/apiService';

// Componente AssociateSubjectForm recebe students, subjects, styles e onAssociateSuccess como props
export default function AssociateSubjectForm({ students, subjects, styles, onAssociateSuccess }) {
    // Estados para Associação Aluno-Matéria
    const [selectedStudentId, setSelectedStudentId] = useState(''); // Inicializa como string vazia
    const [selectedSubjectId, setSelectedSubjectId] = useState(''); // Inicializa como string vazia
    const [associationMessage, setAssociationMessage] = useState('');
    const [associationMessageType, setAssociationMessageType] = useState('');

    // Efeito para atualizar as seleções iniciais quando as listas de alunos/matérias são carregadas
    useEffect(() => {
        if (students.length > 0 && selectedStudentId === '') {
            setSelectedStudentId(students[0].id);
        }
        if (subjects.length > 0 && selectedSubjectId === '') {
            setSelectedSubjectId(subjects[0].id);
        }
    }, [students, subjects]); // Depende de students e subjects

    async function handleAssociateSubject(e) {
        e.preventDefault();
        setAssociationMessage('');
        setAssociationMessageType('');

        if (!selectedStudentId || !selectedSubjectId) {
            setAssociationMessage('Selecione um aluno e uma matéria para associar.');
            setAssociationMessageType('error');
            return;
        }

        try {
            // Usa o studentService para adicionar a matéria
            await studentService.addSubject(selectedStudentId, selectedSubjectId);

            setAssociationMessage('Matéria associada ao aluno com sucesso!');
            setAssociationMessageType('success');
            if (onAssociateSuccess) {
                onAssociateSuccess(); // Chama o callback para recarregar alunos na Home
            }
        } catch (error) {
            console.error('Erro ao associar matéria:', error);
            setAssociationMessage(`Erro ao associar matéria: ${error.message}`);
            setAssociationMessageType('error');
        }
    }

    return (
        <>
            <hr style={styles.hr} />
            <h2 style={styles.h2}>Associar Matéria a Aluno</h2>
            <form onSubmit={handleAssociateSubject} style={styles.form}>
                <select
                    value={selectedStudentId}
                    onChange={(e) => setSelectedStudentId(e.target.value)}
                    required
                    style={styles.input}
                >
                    <option value="">Selecione um Aluno</option>
                    {Array.isArray(students) && students.map(student => (
                        <option key={student.id} value={student.id}>
                            {student.name} ({student.enrollment})
                        </option>
                    ))}
                </select>

                <select
                    value={selectedSubjectId}
                    onChange={(e) => setSelectedSubjectId(e.target.value)}
                    required
                    style={styles.input}
                >
                    <option value="">Selecione uma Matéria</option>
                    {Array.isArray(subjects) && subjects.map(subject => (
                        <option key={subject.id} value={subject.id}>
                            {subject.name} (Ano {subject.year})
                        </option>
                    ))}
                </select>

                <button type="submit" style={styles.button}>Associar Matéria</button>
                {associationMessage && (
                    <div style={{ ...styles.message, color: associationMessageType === 'error' ? '#c0392b' : '#27ae60', backgroundColor: associationMessageType === 'error' ? '#fde0dc' : '#d4edda', border: `1px solid ${associationMessageType === 'error' ? '#e74c3c' : '#28a745'}` }}>
                        {associationMessage}
                    </div>
                )}
            </form>
        </>
    );
}