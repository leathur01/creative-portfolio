package models

// TODO: Refactor to manage carousels in the service class

import (
	"creative-portfolio/lio/app"
	"regexp"
	"time"

	"github.com/revel/revel"
)

type Carousel struct {
	Id          int
	Order       int    `json:"order"`
	FilePath    string `json:"filePath"`
	FileSize    int    `json:"fileSize"`
	FileType    string `json:"fileType"`
	ContentType string
	UploadedAt  time.Time
}

const (
	_      = iota             // iota = 0, value is discarded
	KB int = 1 << (10 * iota) // iota = 1, KB = 1024
	MB                        // iota = 2, MB = 1 << (10 * 2) = 1.048.576 1,042,592
	GB                        // iota = 3, GB = 1 << (10 * 3) = 1073741824
)

// TODO: Allow the admin to set this limitation
var carouselLimitOnUI = 5
var currentCarousels = make(map[int]int)
var fileNameRegex = regexp.MustCompile("^[a-zA-Z0-9-._'() ]+$")
var fileTypeRegex = regexp.MustCompile("^(jpeg|png)$")
var contentTypeRegex = regexp.MustCompile("^image$")

func NewCarousel() Carousel {
	return Carousel{}
}

func (carousel *Carousel) Validate(v *revel.Validation) {
	// For some unknown reasons, revel doesn't set the key value for this validation
	// So I set the key manually
	v.Range(carousel.Order, 0, carouselLimitOnUI).Key("carousel order")

	v.Required(carousel.FileSize)
	validationResult := v.Range(carousel.FileSize, 2*KB, 1*MB)
	validationResult.Key("carousel image file size")
	validationResult.Message("the size of the image has to be between 2KB and 1MB")

	v.Required(carousel.FileType).Key("carousel file type")
	v.Match(carousel.FileType, fileTypeRegex).Message("the file type can only be jpeg or png")

	v.Required(carousel.ContentType).Key("carousel content type")
	v.Match(carousel.ContentType, contentTypeRegex).Message("the file has to be an image")
}

func InsertCarousel(c Carousel) (int, error) {
	query := `
		INSERT INTO carousel("order", file_path, file_size, file_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var uploadedImageId int
	args := []interface{}{c.Order, c.FilePath, c.FileSize, c.FileType}
	err := app.DB.QueryRow(query, args...).Scan(
		&uploadedImageId,
	)

	return uploadedImageId, err
}
