package controllers

import (
	"creative-portfolio/app/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/revel/revel"
)

type Users struct {
	*revel.Controller
}

func (c Users) Create() revel.Result {
	user := models.NewUser()

	//Parsing data from the view
	user.Name = "p minh"
	user.Email = "phamhongthai@gmail.com"

	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		// TODO:
		//Redirect to the creating form
	}

	err := models.InsertUser(user)
	if err != nil {
		revel.AppLog.Fatal(err.Error())
	}

	return c.RenderJSON(user)
}

func (c Users) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		data["error"] = "Query parameter id must be int"
		return c.RenderJSON(data)
	}

	user, err := models.GettUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Response.Status = http.StatusNotFound
			data["error"] = "Resource not found"
			revel.AppLog.Error(err.Error())
			return c.RenderJSON(data)
		}

		c.Response.Status = http.StatusInternalServerError
		data["error"] = err.Error()
		return c.RenderJSON(data)
	}

	return c.RenderJSON(user)
}

func (c Users) GetAll() revel.Result {
	data := make(map[string]interface{})

	users, err := models.GetAllUsers()
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		data["error"] = "Server error"
		revel.AppLog.Fatal(err.Error())
		return c.RenderJSON(data)
	}

	return c.RenderJSON(users)
}

func (c Users) Update() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		data["error"] = "Invalid Id parameter"
		return c.RenderJSON(data)
	}

	user, err := models.GettUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Response.Status = http.StatusNotFound
			data["error"] = "Resource not found"
			revel.AppLog.Error(err.Error())
			return c.RenderJSON(data)
		}

		c.Response.Status = http.StatusInternalServerError
		data["error"] = err.Error()
		return c.RenderJSON(data)
	}

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	// Parsing data
	err = c.Params.BindJSON(&input)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		badRequest := errors.As(err, &syntaxError) || errors.As(err, &invalidUnmarshalError) || errors.As(err, &unmarshalTypeError)
		if badRequest {
			c.Response.Status = http.StatusBadRequest
			data["error"] = err.Error()
			revel.AppLog.Error(err.Error())
			return c.RenderJSON(data)
		}

		c.Response.Status = http.StatusInternalServerError
		data["error"] = "Server error"
		return c.RenderJSON(data)
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
		c.Response.Status = http.StatusBadRequest
		data["error"] = err.Error()
		return c.RenderJSON(data)
	}

	err = models.UpdateUser(*user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Response.Status = http.StatusNotFound
			data["error"] = "Resource not found"
			revel.AppLog.Error(err.Error())
			return c.RenderJSON(data)
		}

		c.Response.Status = http.StatusInternalServerError
		data["error"] = err.Error()
		return c.RenderJSON(data)
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
			c.Response.Status = http.StatusNotFound
			data["error"] = "Resource not found"
			revel.AppLog.Error(err.Error())
			return c.RenderJSON(data)
		}

		c.Response.Status = http.StatusInternalServerError
		data["error"] = err.Error()
		return c.RenderJSON(data)
	}

	return c.RenderJSON(userId)
}
