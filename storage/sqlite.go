package storage

import (
	"database/sql"
	"fmt"
	"todo-orion-bot/entity"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS storage(
	id INTEGER PRIMARY KEY,
	task TEXT UNIQUE,
	status TEXT,
	date TEXT,
	count TEXT,
	name TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_date ON storage(date);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(data entity.Param) (int64, error) {
	const op = "storage.sqlite.Save"
	count := ""
	stmt, err := s.db.Prepare("INSERT INTO storage(status, date, name, task, count) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(data.Status, data.Date, data.Name, data.Task, count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *Storage) UpdateCount(data entity.EventData, id int64) error {
	const op = "storage.sqlite.UpdateCount"
	stmt, err := s.db.Prepare("UPDATE storage SET count= :count, id= :id WHERE id = :id")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(sql.Named("count", data.Count), sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}

func (s *Storage) Update(data entity.Param, id int64) (int64, error) {
	const op = "storage.sqlite.Update"
	query := "UPDATE  storage SET id = :id, status = :status, date = :date, task = :task WHERE id = :id"
	res, err := s.db.Exec(query, sql.Named("id", id), sql.Named("date", data.Date), sql.Named("status", data.Status), sql.Named("task", data.Task))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) Delete(id int64) error {
	const op = "storage.sqlite.Delete"
	query := "DELETE FROM storage WHERE id> :id"
	_, err := s.db.Query(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Get(id int64) (entity.EventData, error) {
	const op = "storage.sqlite.Get"
	var Data entity.EventData

	stmt, err := s.db.Prepare("SELECT  task, date, count, status, name FROM storage WHERE id = ?")
	if err != nil {
		return Data, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	err = stmt.QueryRow(id).Scan(&Data.Task, &Data.Date, &Data.Count, &Data.Status, &Data.Name)
	if err != nil {
		return Data, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return Data, nil
}

func (s *Storage) GetID() ([]int64, error) {
	const op = "storage.sqlite.GetID"
	var id int64
	var resID = make([]int64, 0)
	rows, err := s.db.Query("select id from storage")
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("%s: execute statement: %w", op, err)
		}
		resID = append(resID, id)

	}
	return resID, nil
}
