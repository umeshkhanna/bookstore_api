package users

import (
	"bookstore_api/datasources/mysql/users_db"
	"bookstore_api/utils/date_utils"
	"bookstore_api/utils/errors"
	"fmt"
	"strings"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser    = "SELECT * FROM users WHERE id=?;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.ClientDB.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return errors.NewNotFoundError(fmt.Sprintf("User %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("Error when trying to get user %d : %s", user.Id, err.Error()))
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.ClientDB.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	user.DateCreated = date_utils.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), "users.email") {
			return errors.NewBadRequestError(fmt.Sprintf("Email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("Error when trying to save user: %s", err.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Error when trying to save user: %s", err.Error()))
	}
	user.Id = userId
	return nil
}
