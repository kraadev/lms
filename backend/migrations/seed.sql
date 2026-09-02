-- ==============================================================================
-- Learning Management System (LMS) - PostgreSQL Seed Data
-- ==============================================================================

-- 1. USERS (passwords: 'admin123' for admin, 'password123' for teachers & students)
INSERT INTO users (id, name, email, password_hash, role, created_at, updated_at) VALUES
(1, 'System Administrator', 'admin@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'Budi Santoso (Teacher)', 'teacher1@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'teacher', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'Siti Rahma (Teacher)', 'teacher2@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'teacher', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Andi Pratama', 'student1@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'student', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, 'Bunga Citra', 'student2@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'student', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(6, 'Candra Wijaya', 'student3@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'student', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(7, 'Dewi Lestari', 'student4@lms.local', '$2a$10$w8571TsqPqK645y7kO032ODK5Z6wGj9Gf4L2hMvR5Q7sV3N9W6s.m', 'student', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Reset users serial sequence
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- 2. CLASSES
INSERT INTO classes (id, name, description, teacher_id, academic_year, status, created_at, updated_at) VALUES
(1, 'Kelas Backend Go & Microservices', 'Mempelajari arsitektur backend Go, REST API, WebSocket, dan PostgreSQL.', 2, '2026/2027 Ganjil', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'Kelas Frontend & UI Architecture', 'Mempelajari arsitektur SPA modern, state management, dan styling clean.', 3, '2026/2027 Ganjil', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('classes_id_seq', (SELECT MAX(id) FROM classes));

-- 3. CLASS MEMBERS
INSERT INTO class_members (class_id, user_id, joined_at) VALUES
(1, 4, CURRENT_TIMESTAMP), -- Student 1 in Class 1
(1, 5, CURRENT_TIMESTAMP), -- Student 2 in Class 1
(2, 6, CURRENT_TIMESTAMP), -- Student 3 in Class 2
(2, 7, CURRENT_TIMESTAMP)  -- Student 4 in Class 2
ON CONFLICT (class_id, user_id) DO NOTHING;

-- 4. MATERIALS
INSERT INTO materials (id, class_id, teacher_id, title, description, content, external_url, created_at, updated_at) VALUES
(1, 1, 2, 'Modul 01: Pengenalan Concurrency Go', 'Memahami Goroutines, Channels, dan WaitGroups.', '# Concurrency in Go\n\nGo memiliki model konkurensi bawaan bernama CSP.', 'https://go.dev/tour/concurrency/1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('materials_id_seq', (SELECT MAX(id) FROM materials));

-- 5. ASSIGNMENTS
INSERT INTO assignments (id, class_id, teacher_id, title, description, deadline, max_score, created_at, updated_at) VALUES
(1, 1, 2, 'Tugas 1: Implementasi WebSocket Hub di Go', 'Buat concurrency-safe WebSocket room broker menggunakan channels dan mutex.', CURRENT_TIMESTAMP + INTERVAL '7 days', 100.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('assignments_id_seq', (SELECT MAX(id) FROM assignments));

-- 6. QUIZZES
INSERT INTO quizzes (id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at) VALUES
(1, 1, 2, 'Kuis 1: Konsep Dasar Go', 'Uji pemahaman dasar sintaksis, pointer, dan interface di Go.', 30, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP + INTERVAL '14 days', 3, 'published', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('quizzes_id_seq', (SELECT MAX(id) FROM quizzes));

-- 7. QUIZ QUESTIONS & OPTIONS
INSERT INTO quiz_questions (id, quiz_id, question, type, points, order_index) VALUES
(1, 1, 'Apa keyword untuk memulai sebuah Goroutine baru di Go?', 'multiple_choice', 10.0, 1),
(2, 1, 'Di Go, map aman untuk dibaca dan ditulis secara konkuren tanpa mutex.', 'true_false', 10.0, 2)
ON CONFLICT (id) DO NOTHING;

SELECT setval('quiz_questions_id_seq', (SELECT MAX(id) FROM quiz_questions));

INSERT INTO quiz_options (id, question_id, option_text, is_correct) VALUES
(1, 1, 'go', TRUE),
(2, 1, 'async', FALSE),
(3, 1, 'thread', FALSE),
(4, 1, 'spawn', FALSE),
(5, 2, 'True', FALSE),
(6, 2, 'False', TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval('quiz_options_id_seq', (SELECT MAX(id) FROM quiz_options));

-- 8. SAMPLE MESSAGES
INSERT INTO messages (class_id, user_id, message, created_at) VALUES
(1, 2, 'Selamat datang di kelas Backend Go! Silakan periksa materi Modul 01.', CURRENT_TIMESTAMP);

-- 9. NOTIFICATIONS
INSERT INTO notifications (user_id, type, title, message, is_read, created_at) VALUES
(4, 'assignment', 'Tugas Baru Tersedia', 'Guru telah menerbitkan Tugas 1: Implementasi WebSocket Hub di Go.', FALSE, CURRENT_TIMESTAMP);
