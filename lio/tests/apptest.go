package tests

import (
	"creative-portfolio/lio/app/models"
	"fmt"

	"github.com/revel/revel"
	"github.com/revel/revel/testing"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) After() {
	println("Tear down")
}

func (t *AppTest) InsertPortfolio() {
	// Test insertintg portfolio belonging to user with id = 3
	user := models.NewUser()
	user.Id = 3

	portfolio := models.NewPortfolio()
	portfolio.Name = "media"
	portfolio.User = &user

	err := models.InsertPortfolio(portfolio)
	if err != nil {
		revel.AppLog.Fatal(err.Error())
	}

	revel.AppLog.Info("Insert successfully")

}

func (t *AppTest) GetPortfolio() {
	// Test getting portfolio without a related user
	portfolio, err := models.GetPortfolio(2)
	if err != nil {
		revel.AppLog.Fatal(err.Error())
	}

	revel.AppLog.Info(portfolio.String())

	// Test getting portfolion with the related user
	portfolio, err = models.GetPortfolio(2)
	if err != nil {
		revel.AppLog.Fatal(err.Error())
	}

	revel.AppLog.Info(portfolio.String())
}

func (t *AppTest) GetAllPortfolios() {
	portfolios := []*models.Portfolio{}
	portfolios, err := models.GetAllPortfolios()
	if err != nil {
		revel.AppLog.Fatal(err.Error())
	}

	revel.AppLog.Info(fmt.Sprintf(`%+v`, portfolios))
}

func (t *AppTest) GetAllPortfoliosOfUser() {
	userId := 6
	portfolios, err := models.GetAllPortfoliosOfUser(userId)
	if err != nil {
		revel.AppLog.Fatal(err.Error())
	}

	revel.AppLog.Info(fmt.Sprintf(`%+v`, portfolios))
}
