package db

import (
	"fmt"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Image    string
	Sex      string
	PassHash string
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

func (m *Manager) GetUser(ID string) (*User, error) {
	rows, err := m.DB.Query("SELECT * FROM players WHERE players.id=?", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	user := &User{}
	err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Sex, &user.PassHash)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %v", err)
	}
	return user, nil
}

func (m *Manager) SelectUsers() ([]*User, error) {
	rows, err := m.DB.Query("SELECT * FROM players")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Sex, &user.PassHash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *Manager) CreateUser(user *User) error {
	_, err := m.DB.Exec(
		"INSERT INTO players VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Name, user.Email, user.Image, user.Sex, user.PassHash,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (m *Manager) UpdateUser(user *User) error {
	_, err := m.DB.Exec(
		"UPDATE players SET VALUES (?, ?, ?, ?, ?, ?) WHERE players.id=?",
		user.ID, user.Name, user.Email, user.Image, user.Sex, user.PassHash, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (m *Manager) DeleteUser(user *User) error {
	_, err := m.DB.Exec("DELETE FROM players WHERE players.id=?", user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
