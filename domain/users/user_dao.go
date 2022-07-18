package users

import (
	"bookstore_api/datasources/mysql/users_db"
	"bookstore_api/logger"
	"bookstore_api/utils/errors"
	"fmt"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.ClientDB.Prepare(queryGetUser)
	if err != nil {
		logger.Error("Error occurred during the prepare get user statement", err)
		return errors.NewInternalServerError("Database Error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("Error occurred during the get user by id", err)
		return errors.NewInternalServerError("Database Error")
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.ClientDB.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("Error occurred during the prepare save user statement", err)
		return errors.NewInternalServerError("Database Error")
	}
	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("Error occurred during the save user", saveErr)
		return errors.NewInternalServerError("Database Error")
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("Error occurred during the get user last inserted id", err)
		return errors.NewInternalServerError("Database Error")
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.ClientDB.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("Error occurred during the prepare update user", err)
		return errors.NewInternalServerError("Database Error")
	}
	defer stmt.Close()
	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updateErr != nil {
		logger.Error("Error occurred during the update user by id", updateErr)
		return errors.NewInternalServerError("Database Error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.ClientDB.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("Error occurred during the prepare delete user statement", err)
		return errors.NewInternalServerError("Database Error")
	}
	defer stmt.Close()
	_, deleteErr := stmt.Exec(user.Id)
	if deleteErr != nil {
		logger.Error("Error occurred during the delete user by id", deleteErr)
		return errors.NewInternalServerError("Database Error")
	}
	return nil
}

func (user *User) Search(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.ClientDB.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("Error occurred during the prepare get user statement", err)
		return nil, errors.NewInternalServerError("Database Error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("Error occurred during the search user statement", err)
		return nil, errors.NewInternalServerError("Database Error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("Error occurred during the get user by status", err)
			return nil, errors.NewInternalServerError("Database Error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
