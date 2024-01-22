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

type PortfolioView struct {
	*revel.Controller
}

func (c PortfolioView) Create(portfolio *models.Portfolio) revel.Result {
	// http form only supports post and get
	// This code forward the request to the proper request handler
	method := c.Params.Form.Get("method")
	if method == "put" {
		return c.Update(portfolio)
	} else if method == "delete" {
		return c.Delete()
	}

	data := make(map[string]interface{})

	user := models.NewUser()
	userId := c.Params.Form.Get("user-id")
	if userId != "" {
		var err error
		user.Id, err = strconv.Atoi(userId)
		if err != nil {
			// TODO: Redirect user to the form and display errors
			return helpers.BadRequestResponse(data, "invalid id parameter", c.Controller)
		}
	}
	portfolio.User = &user

	portfolio.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect("portfolios/new")
	}

	err := models.InsertPortfolio(*portfolio)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/portfolios")
}

func (c PortfolioView) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	portfolioId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	portfolio, err := models.GetPortfolio(portfolioId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	// Prevent cyclic json
	// We should not show null value when the data is actually exist
	// With null value, omitempty will not marshlle the field into json
	// This is equivilent to having a DTO
	portfolio.User.Portfolios = []*models.Portfolio{}

	c.ViewArgs["portfolio"] = portfolio
	return c.RenderTemplate("Portfolios/show.html")
}

func (c PortfolioView) GetAll() revel.Result {
	data := make(map[string]interface{})

	var portfolios []*models.Portfolio
	userId := c.Params.Query.Get("user-id")
	if userId != "" {
		userId, err := strconv.Atoi(userId)
		if err != nil {
			return helpers.BadRequestResponse(data, "Invalid user id", c.Controller)
		}

		portfolios, err = models.GetAllPortfoliosOfUser(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return helpers.NotFoundResponse(data, c.Controller)
			}
			return helpers.ServerErrorResponse(data, err, c.Controller)
		}

	} else {
		var err error
		portfolios, err = models.GetAllPortfolios()
		if err != nil {
			return helpers.ServerErrorResponse(data, err, c.Controller)
		}
	}

	users, err := models.GetAllUsers()
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	c.ViewArgs["users"] = users
	c.ViewArgs["portfolios"] = portfolios
	c.ViewArgs["userId"] = userId
	return c.RenderTemplate("Portfolios/index.html")
}

func (c PortfolioView) Update(portfolio *models.Portfolio) revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	portfolioId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	// Check if the portfolio exists
	temp, err := models.GetPortfolio(portfolioId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}
	portfolio.Id = temp.Id
	tempUser := models.NewUser()
	portfolio.User = &tempUser
	portfolio.User.Id = temp.User.Id

	// Prevent cyclic json
	// We should not show null value when the data is actually exist
	// With null value, omitempty will not marshlle the field into json
	// TODO: Refactor to DTO
	portfolio.User.Portfolios = []*models.Portfolio{}

	portfolio.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect("/portfolios/%d/edit", portfolio.Id)
	}

	err = models.UpdatePortfolio(*portfolio)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/portfolios/%d", portfolio.Id)
}

func (c PortfolioView) Delete() revel.Result {
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
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/portfolios")
}

func (c PortfolioView) Form() revel.Result {
	data := make(map[string]interface{})
	var users []*models.User

	// Return the form to create a portfolio of a specific user
	// if the user id is specified in the query parameters
	// The user's id of them form will be populated by this user id parameter
	id := c.Params.Query.Get("user-id")
	if id != "" {
		userId, err := strconv.Atoi(id)
		if err != nil {
			return helpers.BadRequestResponse(data, err.Error(), c.Controller)
		}

		user, err := models.GetUser(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return helpers.NotFoundResponse(data, c.Controller)
			}
			return helpers.ServerErrorResponse(data, err, c.Controller)
		}
		users = append(users, user)

	} else {
		id = c.Params.Route.Get("id")
		if id == "" {
			var err error
			users, err = models.GetAllUsers()
			if err != nil {
				return helpers.ServerErrorResponse(data, err, c.Controller)
			}
		} else {

			// Return the update form of a portfolio since there is a portfolio id
			portfolioId, err := strconv.Atoi(id)
			if err != nil {
				return helpers.BadRequestResponse(data, "invalid id parameter", c.Controller)
			}

			portfolio, err := models.GetPortfolio(portfolioId)
			if err != nil {
				return helpers.ServerErrorResponse(data, err, c.Controller)
			}

			users = append(users, portfolio.User)
			c.ViewArgs["portfolio"] = portfolio
		}
	}

	c.ViewArgs["users"] = users
	return c.RenderTemplate("Portfolios/form.html")
}
