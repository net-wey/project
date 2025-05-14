package postgres

import (
	"database/sql"
	"fmt"
	"goproject/internal/storage/postgres/Entity"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(url string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveDeveloper(developer Entity.Developer) (uint, error) {
	const op = "storage.postgres.SaveDeveloper"

	if developer.Firstname == "" || developer.LastName == "" {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidDeveloperData)
	}

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

func (s *Storage) GetDeveloperById(id uint) (Entity.Developer, error) {
	const op = "storage.postgres.GetDeveloperById"

	stmt, err := s.db.Prepare(`		SELECT id, firstname, last_name, created_at, modified_at, deleted_at 
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
		if err == sql.ErrNoRows {
			return Entity.Developer{}, fmt.Errorf("%s: %w", op, ErrDeveloperNotFound)
		}
		return Entity.Developer{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return developer, nil
}

func (s *Storage) GetDevelopers() ([]Entity.Developer, error) {
	const op = "storage.postgres.GetDevelopers"

	stmt, err := s.db.Prepare(`
		SELECT id, firstname, last_name, created_at, modified_at, deleted_at 
		FROM developers WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	developers := make([]Entity.Developer, 0)
	for rows.Next() {
		var developer Entity.Developer
		if err := rows.Scan(
			&developer.ID,
			&developer.Firstname,
			&developer.LastName,
			&developer.CreatedAt,
			&developer.ModifiedAt,
			&developer.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		developers = append(developers, developer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return developers, nil
}

func (s *Storage) DeleteDeveloper(id uint) error {
	const op = "storage.postgres.DeleteDeveloper"

	stmt, err := s.db.Prepare(`
		UPDATE developers SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrDeveloperNotFound)
	}

	return nil
}

func (s *Storage) SaveTask(task Entity.Task) (uint, error) {
	const op = "storage.postgres.SaveTask"

	if task.Name == "" || task.EstimatePlaned <= 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidTaskData)
	}

	stmt, err := s.db.Prepare(
		`INSERT INTO tasks(
			report_id,
			project_id,
			name,
			developer_note,
			estimate_planed,
			estimate_progress,
			start_timestamp,
			end_timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(
		task.ReportID,
		task.ProjectID,
		task.Name,
		task.DeveloperNote,
		task.EstimatePlaned,
		task.EstimateProgress,
		task.StartTimestamp,
		task.EndTimestamp,
	).Scan(&id, &task.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return id, nil
}
