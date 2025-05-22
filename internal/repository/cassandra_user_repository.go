// filepath: /home/luoxinyu01/github/todolist/internal/repository/cassandra_user_repository.go
package repository

import (
	"errors"
	"fmt"
	"todolist/internal/models"

	"github.com/gocql/gocql"
)

type CassandraUserRepository struct {
	session *gocql.Session
}

func NewCassandraUserRepository(session *gocql.Session) *CassandraUserRepository {
	return &CassandraUserRepository{session: session}
}

func (repo *CassandraUserRepository) AddUser(user models.User) error {
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}

	if !user.Active {
		return errors.New("user must be active upon creation")
	}

	query := "INSERT INTO users (username, password, active) VALUES (?, ?, ?)"
	if err := repo.session.Query(query, user.Username, user.Password, user.Active).Exec(); err != nil {
		return err
	}
	return nil
}

func (repo *CassandraUserRepository) GetUser(username string) (models.User, error) {
	var user models.User
	query := "SELECT username, password, active FROM users WHERE username = ?"
	if err := repo.session.Query(query, username).Scan(&user.Username, &user.Password, &user.Active); err != nil {
		if err == gocql.ErrNotFound {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (repo *CassandraUserRepository) Authenticate(username, password string) bool {
	var user models.User // Use the existing User model

	query := "SELECT password, active FROM users WHERE username = ?"
	if err := repo.session.Query(query, username).Scan(&user.Password, &user.Active); err != nil {
		if err == gocql.ErrNotFound {
			// User not found
			return false
		}
		fmt.Println("Error querying user for authentication:", err) // Consider using a proper logger
		return false
	}
	if !user.Active {
		// User is not active
		return false
	}

	return user.Password == password
}

func (repo *CassandraUserRepository) DeactivateUser(username string) error {
	query := "UPDATE users SET active = false WHERE username = ?"
	if err := repo.session.Query(query, username).Exec(); err != nil {
		return err
	}
	return nil
}
