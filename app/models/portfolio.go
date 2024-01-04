package models

import (
	"fmt"

	"github.com/revel/revel"
)

type Portfolio struct {
	Id   int
	Name string
	User User
}

func (p *Portfolio) String() string {
	return fmt.Sprintf("Portfolio(%s, %s)", p.Name, p.User.Name)
}

func (p *Portfolio) Validate(v *revel.Validation) {
	v.Check(p.Name,
		revel.Required{},
		revel.MaxSize{Max: 50},
	)

	v.Check(p.User.Id,
		revel.Required{})
}
