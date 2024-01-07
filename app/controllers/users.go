package controllers

import (
	"creative-portfolio/app/models"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/revel/revel"
)

type Users struct {
	*revel.Controller
}

func (c Users) Create() revel.Result {
	data := make(map[string]interface{})

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	err := c.Params.BindJSON(&input)
	if err != nil {
		return badRequestResponse(data, err.Error(), c.Controller)
	}

	user := models.NewUser()
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

	err = models.InsertUser(user)
	if err != nil {
		return serverErrorResponse(data, c.Controller)
	}

	return c.RenderJSON(user)
}

func (c Users) Get() revel.Result {
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

		return serverErrorResponse(data, c.Controller)
	}

	return c.RenderJSON(user)
}

func (c Users) GetAll() revel.Result {
	data := make(map[string]interface{})

	users, err := models.GetAllUsers()
	if err != nil {
		return serverErrorResponse(data, c.Controller)
	}

	return c.RenderJSON(users)
}

func (c Users) Update() revel.Result {
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

		return serverErrorResponse(data, c.Controller)
	}

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	// Parsing data
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
		return serverErrorResponse(data, c.Controller)
	}

	return c.RenderJSON(user)
}

func (c Users) Delete() revel.Result {
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

		return serverErrorResponse(data, c.Controller)
	}

	return c.RenderJSON(userId)
}
