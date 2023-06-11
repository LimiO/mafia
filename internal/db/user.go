package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Filter func(*User)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Sex      string `json:"sex"`
	PassHash string `json:"pass_hash"`
}

func (m *Manager) CreateUserTable() error {
	_, err := m.DB.Exec(`CREATE TABLE IF NOT EXISTS players (
    	id text PRIMARY KEY,
    	name text,
    	email text,
    	image text,
    	sex text,
    	pass_hash text
   	);`)
	if err != nil {
		return fmt.Errorf("failed to create table players: %v", err)
	}
	return nil
}

func commitTxn(tx *sql.Tx) {
	func() {
		err := tx.Commit()
		if err != nil {
			fmt.Printf("failed to commit txn: %v", err)
		}
	}()
}

func (m *Manager) GetUser(ID string, filters ...Filter) (*User, error) {
	row := m.DB.QueryRow("SELECT * FROM players WHERE players.id=?", ID)
	if row == nil {
		return nil, fmt.Errorf("failed to get user")
	}
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Sex, &user.PassHash)
	if err != nil {
		return nil, nil
	}
	for _, filter := range filters {
		filter(&user)
	}
	return &user, nil
}

func (m *Manager) SelectUsers(filters ...Filter) ([]*User, error) {
	rows, err := m.DB.Query("SELECT * FROM players")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Sex, &user.PassHash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}
		for _, filter := range filters {
			filter(&user)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (m *Manager) CreateUser(user *User) error {
	tx, err := m.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx")
	}
	defer commitTxn(tx)
	_, err = m.DB.Exec(
		"INSERT INTO players VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Name, user.Email, user.Image, user.Sex, user.PassHash,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	_, err = m.DB.Exec(
		"INSERT INTO stats VALUES (?, 0, 0, 0)", user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to create stats: %v", err)
	}
	return nil
}

func (m *Manager) UpdateUser(user *User) error {
	tx, err := m.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx")
	}
	defer commitTxn(tx)
	_, err = m.DB.Exec(
		`UPDATE players SET name = ?,
                   				  email = ?,
                   				  image = ?,
                   				  sex = ?
                   WHERE id = ?`,
		user.Name, user.Email, user.Image, user.Sex, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (m *Manager) DeleteUser(user *User) error {
	tx, err := m.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx")
	}
	defer commitTxn(tx)
	_, err = m.DB.Exec("DELETE FROM players WHERE id = ?", user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
