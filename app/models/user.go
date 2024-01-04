package models

import (
	"creative-portfolio/app"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/revel/revel"
)

type User struct {
	Id         int
	Name       string
	Email      string
	CreatedAt  time.Time
	Portfolios []Portfolio
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Name)
}

func (u *User) Validate(v *revel.Validation) {
	v.Check(u.Name,
		revel.Required{},
		revel.MaxSize{Max: 15},
		revel.MinSize{Min: 4},
	)

	v.Check(u.Email,
		revel.Required{},
		revel.Match{Regexp: emailRegex},
	)
}

func InsertUser(u User) error {
	query := `
		INSERT INTO "user"(name, email) 
		VALUES($1, $2);
	`
	args := []interface{}{u.Name, u.Email}
	_, err := app.DB.Exec(query, args...)
	return err
}

func GettUser(id int) (*User, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT id, name, email, created_at
		FROM "user"
		WHERE id = $1;
	`
	user := User{}
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

func GetAllUser() ([]*User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM "user"
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

	return users, err
}

func UpdateUser(u User) error {
	query := `
		UPDATE "user" 
		SET name = $1, email = $2
		WHERE id = $3;
	`

	args := []interface{}{u.Name, u.Email, u.Id}
	_, err := app.DB.Exec(query, args...)
	if err != nil {
		return err //no row case
	}

	return nil
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
	if err != nil {
		return err
	}

	return nil
}

// func SelectUserWithPortfolios(id int) []*User {

// }
