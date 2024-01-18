package views

import (
	"creative-portfolio/app/models"
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

	html := c.Params.Query.Get("html")
	if html == "" {

		var input struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}

		err := c.Params.BindJSON(&input)
		if err != nil {
			return badRequestResponse(data, err.Error(), c.Controller)
		}

		if input.Name != nil {
			user.Name = *input.Name
		}

		if input.Email != nil {
			user.Email = *input.Email
		}

		user.Validate(c.Validation)
		if c.Validation.HasErrors() {
			return failedValidationResponse(data, c.Validation.Errors, c.Controller)
		}

		err = models.InsertUser(user)
		if err != nil {
			return serverErrorResponse(data, err, c.Controller)
		}

		return c.RenderJSON(user)
	} else if html == "true" {
		user.Name = c.Params.Form.Get("name")
		user.Email = c.Params.Form.Get("email")
		user.Validate(c.Validation)
		if c.Validation.HasErrors() {
			// TODO:
			// Redirect user to the form and display errors
			return failedValidationResponse(data, c.Validation.Errors, c.Controller)
		}

		err := models.InsertUser(user)
		if err != nil {
			return serverErrorResponse(data, err, c.Controller)
		}

		return c.Redirect("/users/?html=true")
	}

	return badRequestResponse(data, "Invalid html parameter", c.Controller)
}

func (c UserView) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return badRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	user, err := models.GettUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFoundResponse(data, c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)
	}

	html := c.Params.Query.Get("html")
	if html == "" {
		return c.RenderJSON(user)
	} else if html == "true" {
		c.ViewArgs["user"] = user
		return c.RenderTemplate("Users/show.html")
	}

	return badRequestResponse(data, "Invalid html parameter", c.Controller)
}

func (c UserView) GetAll() revel.Result {
	data := make(map[string]interface{})

	users, err := models.GetAllUsers()
	if err != nil {
		return serverErrorResponse(data, err, c.Controller)
	}

	html := c.Params.Query.Get("html")
	if html == "" {
		return c.RenderJSON(users)
	} else if html == "true" {
		c.ViewArgs["users"] = users
		return c.RenderTemplate("Users/index.html")
	}

	return badRequestResponse(data, "Invalid html parameter", c.Controller)
}

func (c UserView) Update() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return badRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	user, err := models.GettUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFoundResponse(data, c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)
	}

	html := c.Params.Query.Get("html")
	// Parsing JSON data
	if html == "" {
		var input struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}

		err = c.Params.BindJSON(&input)
		if err != nil {
			return badRequestResponse(data, err.Error(), c.Controller)
		}

		if input.Name != nil {
			user.Name = *input.Name
		}

		if input.Email != nil {
			user.Email = *input.Email
		}

		user.Validate(c.Validation)
		if c.Validation.HasErrors() {
			// TODO:
			// Redirect user to the form and display errors
			return failedValidationResponse(data, c.Validation.Errors, c.Controller)
		}

		err = models.UpdateUser(*user)
		if err != nil {
			return serverErrorResponse(data, err, c.Controller)
		}

		return c.RenderJSON(user)

	} else if html == "true" {
		user.Name = c.Params.Form.Get("name")
		user.Email = c.Params.Get("email")

		user.Validate(c.Validation)
		if c.Validation.HasErrors() {
			// TODO:
			// Redirect user to the form and display errors
			return failedValidationResponse(data, c.Validation.Errors, c.Controller)
		}

		err = models.UpdateUser(*user)
		if err != nil {
			return serverErrorResponse(data, err, c.Controller)
		}

		return c.Redirect("/users/%d/?html=true", user.Id)
	}

	return badRequestResponse(data, "Invalid html parameter", c.Controller)
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
			return badRequestResponse(data, "Invalid id parameter", c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)
	}

	html := c.Params.Query.Get("html")
	if html == "" {
		return c.RenderJSON(userId)
	} else if html == "true" {
		return c.Redirect("/users/?html=true")
	}

	return badRequestResponse(data, "Invalid html parameter", c.Controller)
}

// Return a creat form or update form depend on the provided user-id
func (c UserView) Form() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	if id != "" {
		userId, err := strconv.Atoi(id)
		if err != nil {
			return badRequestResponse(data, "Invalid id parameter", c.Controller)
		}

		user, err := models.GettUser(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return notFoundResponse(data, c.Controller)
			}

			return serverErrorResponse(data, err, c.Controller)
		}
		c.ViewArgs["user"] = user
	}

	return c.RenderTemplate("Users/form.html")
}
