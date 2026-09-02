package database

import (
	"fmt"
	"log"
)

func RunMigrations(db *DB) error {
	var primaryKeyType string
	var boolType string
	var boolDefaultFalse string

	if db.Driver == "postgres" {
		primaryKeyType = "BIGSERIAL PRIMARY KEY"
		boolType = "BOOLEAN"
		boolDefaultFalse = "DEFAULT FALSE"
	} else {
		primaryKeyType = "INTEGER PRIMARY KEY AUTOINCREMENT"
		boolType = "BOOLEAN"
		boolDefaultFalse = "DEFAULT 0"
	}

	queries := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users (
			id %s,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			avatar_url TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS classes (
			id %s,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			teacher_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			academic_year VARCHAR(50) NOT NULL DEFAULT '',
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS class_members (
			id %s,
			class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(class_id, user_id)
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS materials (
			id %s,
			class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
			teacher_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			file_path TEXT,
			external_url TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS assignments (
			id %s,
			class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
			teacher_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			attachment_path TEXT,
			deadline TIMESTAMP NOT NULL,
			max_score DOUBLE PRECISION NOT NULL DEFAULT 100.0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS assignment_submissions (
			id %s,
			assignment_id BIGINT NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
			student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			text_answer TEXT,
			file_path TEXT,
			submitted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			score DOUBLE PRECISION,
			feedback TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'submitted',
			UNIQUE(assignment_id, student_id)
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS quizzes (
			id %s,
			class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
			teacher_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			duration_minutes INT NOT NULL DEFAULT 60,
			start_at TIMESTAMP NOT NULL,
			end_at TIMESTAMP NOT NULL,
			max_attempts INT NOT NULL DEFAULT 1,
			status VARCHAR(50) NOT NULL DEFAULT 'published',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS quiz_questions (
			id %s,
			quiz_id BIGINT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
			question TEXT NOT NULL,
			type VARCHAR(50) NOT NULL DEFAULT 'multiple_choice',
			points DOUBLE PRECISION NOT NULL DEFAULT 10.0,
			order_index INT NOT NULL DEFAULT 0
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS quiz_options (
			id %s,
			question_id BIGINT NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
			option_text TEXT NOT NULL,
			is_correct %s NOT NULL %s
		);`, primaryKeyType, boolType, boolDefaultFalse),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS quiz_attempts (
			id %s,
			quiz_id BIGINT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
			student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			submitted_at TIMESTAMP,
			score DOUBLE PRECISION,
			status VARCHAR(50) NOT NULL DEFAULT 'in_progress'
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS quiz_answers (
			id %s,
			attempt_id BIGINT NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
			question_id BIGINT NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
			selected_option_id BIGINT REFERENCES quiz_options(id) ON DELETE SET NULL,
			text_answer TEXT,
			is_correct %s,
			earned_points DOUBLE PRECISION
		);`, primaryKeyType, boolType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS meetings (
			id %s,
			class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
			teacher_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			room_name VARCHAR(255) NOT NULL UNIQUE,
			type VARCHAR(50) NOT NULL DEFAULT 'video',
			status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
			started_at TIMESTAMP,
			ended_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS messages (
			id %s,
			class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			message TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS notifications (
			id %s,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			is_read %s NOT NULL %s,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`, primaryKeyType, boolType, boolDefaultFalse),

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_class_members_lookup ON class_members(class_id, user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_class_created ON messages(class_id, created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_assignments_class ON assignments(class_id);`,
		`CREATE INDEX IF NOT EXISTS idx_materials_class ON materials(class_id);`,
		`CREATE INDEX IF NOT EXISTS idx_quizzes_class ON quizzes(class_id);`,
		`CREATE INDEX IF NOT EXISTS idx_meetings_class ON meetings(class_id);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, is_read);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("migration error executing [%s]: %w", query, err)
		}
	}

	log.Println("[MIGRATIONS] Database schema migrations completed successfully.")
	return nil
}
