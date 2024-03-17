package views

import (
	"bytes"
	"creative-portfolio/lio/app/controllers/helpers"
	"creative-portfolio/lio/app/models"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/twinj/uuid"
)

type CarouselView struct {
	*revel.Controller
}

func init() {
	revel.InterceptMethod(CarouselView.PopulateCarouselOrderCache, revel.BEFORE)
}

func (c CarouselView) PopulateCarouselOrderCache() revel.Result {
	if models.CurrentCarousels.Initialized {
		return nil
	}

	carousels, err := models.GetAllCarouselsForCaching()
	if err != nil {
		data := make(map[string]interface{})
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	for _, carousel := range carousels {
		if carousel.Order > 0 {
			models.CurrentCarousels.Carousels[carousel.Order] = carousel.Id
		}
	}

	models.CurrentCarousels.Initialized = true
	return nil
}

// TODO: Extract some code into the helper or service file
func (c CarouselView) Upload(carouselImage []byte, carousel *models.Carousel) revel.Result {
	data := make(map[string]interface{})

	method := c.Params.Form.Get("method")
	if method == "put" {
		return c.Update()
	} else if method == "delete" {
		return c.Delete()
	}

	contentType := http.DetectContentType(carouselImage)
	// Take the content type in the first part
	carousel.ContentType = strings.Split(contentType, "/")[0]

	// Take the file extension in the second part and discard the utf part
	extension := strings.Split(strings.Split(contentType, "/")[1], ";")[0]
	carousel.FileType = extension

	// Calculate file size
	size, err := io.Copy(io.Discard, bytes.NewReader(carouselImage))
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}
	carousel.FileSize = int(size)

	carousel.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return helpers.FailedValidationResponse(data, c.Validation.Errors, c.Controller)
	}

	uuid := uuid.NewV1()
	uuidFileName := uuid.String() + "." + carousel.FileType
	carousel.FilePath = uuidFileName

	filepath := "storage/image/" + carousel.FilePath
	err = os.WriteFile(filepath, carouselImage, 0644)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	oldCarouselId := models.CurrentCarousels.Carousels[carousel.Order]
	if oldCarouselId != 0 {
		err = models.UpdateCarousel(oldCarouselId, 0)
		if err != nil {
			return helpers.ServerErrorResponse(data, err, c.Controller)
		}
	}

	uploadedCarouselId, err := models.InsertCarousel(*carousel)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}
	models.CurrentCarousels.Carousels[carousel.Order] = uploadedCarouselId

	return c.Redirect("/carousels")
}

func (c CarouselView) Get() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	carouselId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	carousel, err := models.GetCarousel(carouselId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	// 0 value is used for null order.
	// So we need to populate the array up to 6 value
	allowedOrders := make([]int, models.CarouselLimitOnUI+1)
	for i := range allowedOrders {
		allowedOrders[i] = i
	}
	c.ViewArgs["allowedOrders"] = allowedOrders

	c.ViewArgs["carousel"] = carousel
	return c.RenderTemplate("Carousels/show.html")
}

func (c CarouselView) GetAll() revel.Result {
	data := make(map[string]interface{})

	carousels, err := models.GetAllCarousels()
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	c.ViewArgs["carousels"] = carousels
	return c.RenderTemplate("Carousels/index.html")
}

func (c CarouselView) Update() revel.Result {
	data := make(map[string]interface{})

	carouselId, err := strconv.Atoi(c.Params.Route.Get("id"))
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	intOrder, err := strconv.Atoi(c.Params.Form.Get("order"))
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	oldCarouselId := models.CurrentCarousels.Carousels[intOrder]
	if oldCarouselId != 0 {
		err = models.UpdateCarousel(oldCarouselId, 0)
		if err != nil {
			return helpers.ServerErrorResponse(data, err, c.Controller)
		}
	}

	err = models.UpdateCarousel(carouselId, intOrder)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}
	models.CurrentCarousels.Carousels[intOrder] = carouselId

	return c.Redirect("/carousels/%d", carouselId)
}

func (c CarouselView) Delete() revel.Result {
	data := make(map[string]interface{})

	id := c.Params.Route.Get("id")
	carouselId, err := strconv.Atoi(id)
	if err != nil {
		return helpers.BadRequestResponse(data, "Invalid id parameter", c.Controller)
	}

	carousel, err := models.GetCarousel(carouselId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.NotFoundResponse(data, c.Controller)
		}

		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	err = models.DeleteCarousel(carouselId)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	// Delete Cache
	for key, value := range models.CurrentCarousels.Carousels {
		if value == carouselId {
			models.CurrentCarousels.Carousels[key] = 0
		}
	}

	filePath := "storage/image/" + carousel.FilePath
	err = os.Remove(filePath)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.Redirect("/carousels")
}

func (c CarouselView) Form() revel.Result {
	return c.RenderTemplate("Carousels/form.html")
}
