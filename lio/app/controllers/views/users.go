package views

import (
	"creative-portfolio/lio/app/controllers/helpers"
	"creative-portfolio/lio/app/models"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/revel/revel"
)

type UserView struct {
	*revel.Controller
}

func (c UserView) Create() revel.Result {
	// http form only supports post and get
	// This code forward the request to the proper request handler
	method := c.Params.Form.Get("method")
	if method == "put" {
		return c.Update()
	} else if method == "delete" {
		return c.Delete()
	}

	data := make(map[string]interface{})
	user := models.NewUser()

	user.Name = c.Params.Form.Get("name")
	user.Email = c.Params.Form.Get("email")
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// TODO:
		// Redirect user to the form and display errors
		return helpers.FailedValidationResponse(data, c.Validation.Errors, c.Controller)
	}

	err := models.InsertUser(user)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/users")
}

func (c UserView) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	user, err := models.GettUser(userId)
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

func (c UserView) Update() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	user, err := models.GettUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	user.Name = c.Params.Form.Get("name")
	user.Email = c.Params.Get("email")

	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// TODO:
		// Redirect user to the form and display errors
		return helpers.FailedValidationResponse(data, c.Validation.Errors, c.Controller)
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
		c.Response.Status = http.StatusBadRequest
		data["error"] = "Invalid Id parameter"
		return c.RenderJSON(data)
	}

	err = models.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
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

		user, err := models.GettUser(userId)
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
