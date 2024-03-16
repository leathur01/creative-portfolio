package views

import (
	"bytes"
	"creative-portfolio/lio/app/controllers/helpers"
	"creative-portfolio/lio/app/models"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/revel/revel"
	"github.com/twinj/uuid"
)

type CarouselView struct {
	*revel.Controller
}

// TODO: Extract some code into the helper or service file
func (c CarouselView) Upload(carouselImage []byte, carousel *models.Carousel) revel.Result {
	data := make(map[string]interface{})

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

	filepath := "storage/Image/" + carousel.FilePath
	err = os.WriteFile(filepath, carouselImage, 0644)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	_, err = models.InsertCarousel(*carousel)
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	return c.RenderJSON(carousel)
}

func (c CarouselView) GetAll() revel.Result {
	data := make(map[string]interface{})

	carousels, err := models.GetAllCarousel()
	if err != nil {
		return helpers.ServerErrorResponse(data, err, c.Controller)
	}

	c.ViewArgs["carousels"] = carousels
	return c.RenderTemplate("Carousels/index.html")
}

func (c CarouselView) Form() revel.Result {
	return c.RenderTemplate("Carousels/form.html")
}
