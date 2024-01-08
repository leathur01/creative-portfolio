package controllers

import (
	"creative-portfolio/app/models"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/revel/revel"
)

type Portfolios struct {
	*revel.Controller
}

func (c Portfolios) Create() revel.Result {
	data := make(map[string]interface{})

	var input struct {
		Name   string `json:"name"`
		UserId int    `json:"user-id"`
	}

	err := c.Params.BindJSON(&input)
	if err != nil {
		return badRequestResponse(data, err.Error(), c.Controller)
	}

	user, err := models.GettUser(input.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFoundResponse(data, c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)

	}

	user.Id = input.UserId

	portfolio := models.NewPortfolio()
	portfolio.User = user
	portfolio.Name = input.Name

	portfolio.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return failedValidationResponse(data, c.Validation.Errors, c.Controller)
	}

	err = models.InsertPortfolio(portfolio)
	if err != nil {
		return serverErrorResponse(data, err, c.Controller)
	}

	return c.RenderJSON(portfolio)
}

func (c Portfolios) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	portfolioId, err := strconv.Atoi(id)
	if err != nil {
		return badRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	portfolio, err := models.GetPortfolio(portfolioId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFoundResponse(data, c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)
	}

	// Prevent cyclic json
	// We should not show null value when the data is actually exist
	// With null value, omitempty will not marshlle the field into json
	// This is equivilent to having a DTO
	portfolio.User.Portfolios = []*models.Portfolio{}

	return c.RenderJSON(portfolio)
}

func (c Portfolios) GetAll() revel.Result {
	data := make(map[string]interface{})

	var portfolios []*models.Portfolio
	userIdQuery := c.Params.Query.Get("user-id")
	if userIdQuery != "" {
		userId, err := strconv.Atoi(userIdQuery)
		if err != nil {
			return badRequestResponse(data, "Invalid user id", c.Controller)
		}

		portfolios, err = models.GetAllPortfoliosOfUser(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return notFoundResponse(data, c.Controller)
			}
			return serverErrorResponse(data, err, c.Controller)
		}
	} else {
		var err error
		portfolios, err = models.GetAllPortfolios()
		if err != nil {
			return serverErrorResponse(data, err, c.Controller)
		}
	}

	html := c.Params.Query.Get("html")
	if html == "" {
		return c.RenderJSON(portfolios)
	} else if html == "true" {
		c.ViewArgs["portfolios"] = portfolios
		return c.RenderTemplate("Portfolios/index.html")
	}

	return badRequestResponse(data, "Invalid html parameter", c.Controller)
}

func (c Portfolios) Update() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	portfolioId, err := strconv.Atoi(id)
	if err != nil {
		return badRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	portfolio, err := models.GetPortfolio(portfolioId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFoundResponse(data, c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)
	}

	// Prevent cyclic json
	// We should not show null value when the data is actually exist
	// With null value, omitempty will not marshlle the field into json
	// This is equivilent to having a DTO
	portfolio.User.Portfolios = []*models.Portfolio{}

	var input struct {
		Name *string `json:"name"`
	}

	err = c.Params.BindJSON(&input)
	if err != nil {
		return badRequestResponse(data, err.Error(), c.Controller)
	}

	if input.Name != nil {
		portfolio.Name = *input.Name
	}

	portfolio.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return failedValidationResponse(data, c.Validation.Errors, c.Controller)
	}

	err = models.UpdatePortfolio(*portfolio)
	if err != nil {
		return serverErrorResponse(data, err, c.Controller)
	}

	return c.RenderJSON(portfolio)
}

func (c Portfolios) Delete() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	portfolioId, err := strconv.Atoi(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		data["error"] = "Invalid Id parameter"
		return c.RenderJSON(data)
	}

	err = models.DeletePortfolio(portfolioId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return badRequestResponse(data, "Invalid id parameter", c.Controller)
		}

		return serverErrorResponse(data, err, c.Controller)
	}

	return nil
}
