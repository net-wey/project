package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/netwey/feedbox/trainee-rbs-feedbox/internal/storage/postgres/Entity"
)

// ErrDeveloperNotFound returns when developer not found in storage
var ErrDeveloperNotFound = errors.New("developer not found")

type Storage struct {
	db *sql.DB
}

func New(url string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func createTables(db *sql.DB) error {
	const op = "storage.postgres.createTables"

	tables := []string{
		`CREATE TABLE IF NOT EXISTS developers (
			id SERIAL PRIMARY KEY,
			firstname VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_developers_firstname ON developers(firstname);
		CREATE INDEX IF NOT EXISTS idx_developers_lastname ON developers(last_name);`,

		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);`,

		`CREATE TABLE IF NOT EXISTS reports (
			id SERIAL PRIMARY KEY,
			developer_id INTEGER NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_reports_developer ON reports(developer_id);`,

		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			report_id INTEGER NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
			project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			developer_note TEXT,
			estimate_planed INTEGER NOT NULL,
			estimate_progress INTEGER NOT NULL,
			start_timestamp TIMESTAMP NOT NULL,
			end_timestamp TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_report ON tasks(report_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) SaveDeveloper(developer Entity.Developer) (uint, error) {
	const op = "storage.postgres.SaveDeveloper"

	stmt, err := s.db.Prepare(`
		INSERT INTO developers (
			firstname, 
			last_name, 
			deleted_at
		) VALUES ($1, $2, $3) 
		RETURNING id, created_at, modified_at`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(
		developer.Firstname,
		developer.LastName,
		developer.DeletedAt,
	).Scan(&id, &developer.CreatedAt, &developer.ModifiedAt)
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetDeveloper(id uint) (Entity.Developer, error) {
	const op = "storage.postgres.GetDeveloper"

	stmt, err := s.db.Prepare(`
		SELECT id, firstname, last_name, created_at, modified_at, deleted_at 
		FROM developers 
		WHERE id = $1 AND deleted_at IS NULL`)
	if err != nil {
		return Entity.Developer{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var developer Entity.Developer
	err = stmt.QueryRow(id).Scan(
		&developer.ID,
		&developer.Firstname,
		&developer.LastName,
		&developer.CreatedAt,
		&developer.ModifiedAt,
		&developer.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entity.Developer{}, ErrDeveloperNotFound
		}
		return Entity.Developer{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return developer, nil
}

