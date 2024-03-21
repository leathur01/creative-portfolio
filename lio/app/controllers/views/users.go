package views

import (
	"creative-portfolio/lio/app/controllers/helpers"
	"creative-portfolio/lio/app/models"
	"database/sql"
	"errors"
	"strconv"

	"github.com/revel/revel"
)

type UserView struct {
	*revel.Controller
}

func (c UserView) Create(user *models.User) revel.Result {
	// http form only supports post and get
	// This code forward the request to the proper request handler
	method := c.Params.Form.Get("method")
	if method == "put" {
		return c.Update(user)
	} else if method == "delete" {
		return c.Delete()
	}

	data := make(map[string]interface{})

	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect("/users/new")
	}

	generatedId, err := models.InsertUser(*user)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/users/%d", generatedId)
}

func (c UserView) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	user, err := models.GetUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	c.ViewArgs["user"] = user
	return c.RenderTemplate("Users/show.html")
}

func (c UserView) GetAll() revel.Result {
	data := make(map[string]interface{})

	users, err := models.GetAllUsers()
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	c.ViewArgs["users"] = users
	return c.RenderTemplate("Users/index.html")
}

func (c UserView) Update(user *models.User) revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	// Check if the user exists
	temp, err := models.GetUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}
	user.Id = temp.Id

	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect("/users/%d/edit", user.Id)
	}

	err = models.UpdateUser(*user)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/users/%d", user.Id)
}

func (c UserView) Delete() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	err = models.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/users")

}

// Return a creat form or update form depend on the provided user-id
func (c UserView) Form() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	if id != "" {
		userId, err := strconv.Atoi(id)
		if err != nil {
			return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
		}

		user, err := models.GetUser(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return helpers.NotFoundResponse(data, c.Controller)
			}

			return helpers.ServerErrorResponse(data, err, c.Controller)
		}

		c.ViewArgs["user"] = user
	}

	return c.RenderTemplate("Users/form.html")
}
