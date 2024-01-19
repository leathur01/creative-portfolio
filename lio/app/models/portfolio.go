package models

import (
	"creative-portfolio/lio/app"
	"database/sql"
	"fmt"
	"time"

	"github.com/revel/revel"
)

type Portfolio struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created-at,omitempty"`
	User      *User     `json:"user,omitempty"`
}

func NewPortfolio() Portfolio {
	user := NewUser()

	return Portfolio{
		User: &user,
	}
}

func (p *Portfolio) String() string {
	return fmt.Sprintf("Portfolio(%s, %d)", p.Name, p.User.Id)
}

func (p *Portfolio) Validate(v *revel.Validation) {
	v.Check(p.Name,
		revel.Required{},
		revel.MaxSize{Max: 50},
	).Key("portfolio name")
}

func InsertPortfolio(p Portfolio) error {
	query := `
		INSERT INTO portfolio(name, user_id) 
		VALUES($1, $2);
	`

	args := []interface{}{p.Name, p.User.Id}
	_, err := app.DB.Exec(query, args...)
	return err
}

func GetPortfolio(id int) (*Portfolio, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT p.id, p.name, p.created_at, u.id, u.name, u.email
		FROM portfolio as p
		JOIN "user" as u on p.user_id = u.id
		WHERE p.id = $1;
	`

	args := []interface{}{id}
	portfolio := Portfolio{}
	user := NewUser()
	err := app.DB.QueryRow(query, args...).Scan(
		&portfolio.Id,
		&portfolio.Name,
		&portfolio.CreatedAt,
		&user.Id,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	portfolio.User = &user
	user.Portfolios = append(user.Portfolios, &portfolio)

	return &portfolio, nil
}

func GetAllPortfolios() ([]*Portfolio, error) {
	query := `
		SELECT id, name, created_at, user_id 
		FROM portfolio
		ORDER BY id ASC
	`

	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	portfolios := []*Portfolio{}
	for rows.Next() {
		portfolio := NewPortfolio()
		err := rows.Scan(
			&portfolio.Id,
			&portfolio.Name,
			&portfolio.CreatedAt,
			&portfolio.User.Id,
		)

		if err != nil {
			return nil, err
		}

		portfolios = append(portfolios, &portfolio)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return portfolios, nil
}

func GetAllPortfoliosOfUser(userId int) ([]*Portfolio, error) {
	if userId < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT p.id, p.name, p.created_at
		FROM portfolio as p
		WHERE p.user_id = $1;
	`

	rows, err := app.DB.Query(query, userId)
	if err != nil {
		return nil, err //no rows
	}

	defer rows.Close()

	user := NewUser()
	user.Id = userId
	portfolios := []*Portfolio{}
	for rows.Next() {
		portfolio := NewPortfolio()
		rows.Scan(
			&portfolio.Id,
			&portfolio.Name,
			&portfolio.CreatedAt,
		)

		if err != nil {
			return nil, err //server error
		}

		portfolio.User = &user
		portfolios = append(portfolios, &portfolio)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return portfolios, nil
}

func UpdatePortfolio(p Portfolio) error {
	query := `
		UPDATE portfolio
		SET name = $1
		WHERE id = $2;
	`

	args := []interface{}{p.Name, p.Id}
	_, err := app.DB.Exec(query, args...)
	return err
}

func DeletePortfolio(id int) error {
	if id < 1 {
		return sql.ErrNoRows
	}

	query := `
		DELETE
		FROM portfolio
		WHERE id = $1;
	`

	_, err := app.DB.Exec(query, id)
	return err
}
