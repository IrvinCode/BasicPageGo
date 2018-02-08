package models

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("models: email address already in use")
	ErrInvalidCredentials = errors.New("models: invalid user credentials")
)

type Database struct{
	*sql.DB
}

func (db *Database) GetSnippet(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := db.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return s, nil

}

func (db *Database) LatestSnippets() (Snippets, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := db.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := Snippets{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (db *Database) InsertSnippet(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? SECOND))`

	result, err := db.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (db *Database) InsertUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password, admin, created)
VALUES(?, ?, ?, 0, UTC_TIMESTAMP())`

	_, err = db.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return ErrDuplicateEmail
		}
	}

	return err
}

func (db *Database) VerifyUser(email, password string) (int, error, bool) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	row := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email)
	rowAdmin := db.QueryRow("SELECT id, password FROM users WHERE email = ? and admin = '1'", email)

	err := rowAdmin.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		err := row.Scan(&id, &hashedPassword)
		if err == sql.ErrNoRows {
			return 0, ErrInvalidCredentials, false
		} else if err != nil {
			return 0, err, false
		}
		// Check whether the hashed password and plain-text password provided match.
		// If they don't, we return the ErrInvalidCredentials error.
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, ErrInvalidCredentials, false
		} else if err != nil {
			return 0, err, false
		}

		// Otherwise, the password is correct. Return the user ID.
		return id, nil, false

	} else if err != nil {
		return 0, err, false
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials, false
	} else if err != nil {
		return 0, err, false
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil, true
}

func (db *Database) InsertAdmin(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password, admin, created)
VALUES(?, ?, ?, '1', UTC_TIMESTAMP())`

	_, err = db.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return ErrDuplicateEmail
		}
	}

	return err
}

func (db *Database) DeleteSnippet(id string) (error) {
	stmt := `DELETE FROM snippets WHERE id = ?`

	_, err := db.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
