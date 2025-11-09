package db

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/kinedwi/ACS_t/internal/face"
)

type DB struct {
	conn *sql.DB
}

type User struct {
	ID         int
	Name       string
	Descriptor face.Descriptor
}

// New - открывает/создаёт БД
func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	schema := `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  descriptor BLOB NOT NULL
);
CREATE TABLE IF NOT EXISTS logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  event TEXT,
  ts DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS alerts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  message TEXT,
  ts DATETIME DEFAULT CURRENT_TIMESTAMP,
  resolved INTEGER DEFAULT 0
);
`
	if _, err := conn.Exec(schema); err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}

// serialize Descriptor with gob
func encodeDescriptor(d face.Descriptor) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode([]float32(d)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeDescriptor(b []byte) (face.Descriptor, error) {
	var arr []float32
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&arr); err != nil {
		return nil, err
	}
	return face.Descriptor(arr), nil
}

func (d *DB) AddUser(name string, desc face.Descriptor) error {
	b, err := encodeDescriptor(desc)
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(`INSERT INTO users (name, descriptor) VALUES (?, ?)`, name, b)
	if err == nil {
		fmt.Printf("Добавлен пользователь %s\n", name)
	}
	return err
}

func (d *DB) AllUsers() ([]User, error) {
	rows, err := d.conn.Query(`SELECT id, name, descriptor FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []User
	for rows.Next() {
		var id int
		var name string
		var blob []byte
		if err := rows.Scan(&id, &name, &blob); err != nil {
			return nil, err
		}
		desc, err := decodeDescriptor(blob)
		if err != nil {
			return nil, err
		}
		res = append(res, User{ID: id, Name: name, Descriptor: desc})
	}
	return res, nil
}

func (d *DB) LogEvent(name, event string) {
	d.conn.Exec(`INSERT INTO logs (name, event) VALUES (?, ?)`, name, event)
	fmt.Printf("Лог: %s - %s\n", name, event)
}

func (d *DB) AddAlert(msg string) {
	d.conn.Exec(`INSERT INTO alerts (message) VALUES (?)`, msg)
	fmt.Printf("АЛЕРТ: %s\n", msg)
}

// FindBestMatch: перебирает всех пользователей, возвращает имя и расстояние
func (d *DB) FindBestMatch(desc face.Descriptor) (string, float32, bool, error) {
	users, err := d.AllUsers()
	if err != nil {
		return "", 0, false, err
	}
	bestName := ""
	bestDist := float32(1e9)
	found := false
	for _, u := range users {
		dist := face.Distance(desc, u.Descriptor)
		if dist < bestDist {
			bestDist = dist
			bestName = u.Name
			found = true
		}
	}
	return bestName, bestDist, found, nil
}
