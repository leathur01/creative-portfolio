package models

import (
	"creative-portfolio/lio/app"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/revel/revel"
)

type User struct {
	Id         int          `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Email      string       `json:"email,omitempty"`
	CreatedAt  time.Time    `json:"created-at,omitempty"`
	Portfolios []*Portfolio `json:"portfolios,omitempty"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewUser() User {
	return User{
		Portfolios: []*Portfolio{},
	}
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s, %d)", u.Name, u.Id)
}

func (user *User) Validate(v *revel.Validation) {
	v.Required(user.Name)
	v.MaxSize(user.Name, 50)
	v.MinSize(user.Name, 3)

	v.Required(user.Email)
	v.Match(user.Email, emailRegex).Message("The email address appears to be invalid")
}

func InsertUser(u User) (int, error) {
	query := `
		INSERT INTO "user"(name, email) 
		VALUES($1, $2)
		RETURNING id;
	`
	var createdUsesrId int
	args := []interface{}{u.Name, u.Email}
	err := app.DB.QueryRow(query, args...).Scan(
		&createdUsesrId,
	)

	return createdUsesrId, err
}

func GetUser(id int) (*User, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT id, name, email, created_at
		FROM "user"
		WHERE id = $1;
	`
	user := NewUser()
	err := app.DB.QueryRow(query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]*User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM "user"
		ORDER by id ASC
	`
	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUser(u User) error {
	query := `
		UPDATE "user" 
		SET name = $1, email = $2
		WHERE id = $3;
	`

	args := []interface{}{u.Name, u.Email, u.Id}
	_, err := app.DB.Exec(query, args...)
	return err
}

func DeleteUser(id int) error {
	if id < 1 {
		return sql.ErrNoRows
	}

	query := `
		DELETE
		FROM "user"
		WHERE id = $1;
	`

	_, err := app.DB.Exec(query, id)
	return err
}
