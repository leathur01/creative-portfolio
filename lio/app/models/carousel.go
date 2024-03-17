package models

// TODO: Refactor to manage carousels in the service class

import (
	"creative-portfolio/lio/app"
	"database/sql"
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
	ContentType string // Used only for validation
	UploadedAt  time.Time
}

const (
	_      = iota             // iota = 0, value is discarded
	KB int = 1 << (10 * iota) // iota = 1, KB = 1024
	MB                        // iota = 2, MB = 1 << (10 * 2) = 1.048.576
	GB                        // iota = 3, GB = 1 << (10 * 3) = 1073741824
)

// TODO: Allow the admin to set this limitation
var CarouselLimitOnUI = 5
var fileTypeRegex = regexp.MustCompile("^(jpeg|png)$")
var contentTypeRegex = regexp.MustCompile("^image$")

var CurrentCarousels = struct {
	Initialized bool
	Carousels   map[int]int
}{
	Initialized: false,
	Carousels:   make(map[int]int),
}

func NewCarousel() Carousel {
	return Carousel{}
}

func (carousel *Carousel) Validate(v *revel.Validation) {
	// For some unknown reasons, revel doesn't set the key value for this validation
	// So I set the key manually
	validationResult := v.Range(carousel.Order, 0, CarouselLimitOnUI)
	validationResult.Key("carousel order")
	validationResult.Message(`value Out of Range: The number you've entered is out of range. Please enter a value between 0 and %d.`, CarouselLimitOnUI)

	v.Required(carousel.FileSize)
	validationResult = v.Range(carousel.FileSize, 2*KB, 1*MB)
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

func GetAllCarousels() ([]*Carousel, error) {
	query := `
		SELECT *
		FROM carousel
		ORDER by id ASC
	`
	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carousels := []*Carousel{}
	for rows.Next() {
		var carousel Carousel
		err := rows.Scan(
			&carousel.Id,
			&carousel.Order,
			&carousel.FilePath,
			&carousel.FileSize,
			&carousel.FileType,
			&carousel.UploadedAt,
		)

		if err != nil {
			return nil, err
		}

		carousels = append(carousels, &carousel)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return carousels, nil
}

func UpdateCarousel(id, order int) error {
	query := `
		UPDATE carousel 
		SET "order" = $1
		WHERE id = $2
	`

	args := []interface{}{order, id}
	_, err := app.DB.Exec(query, args...)
	return err
}

func GetCarousel(id int) (*Carousel, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		select id, "order", file_path, file_size, file_type, uploaded_at
		from carousel
		where id = $1;
	`

	carousel := NewCarousel()
	err := app.DB.QueryRow(query, id).Scan(
		&carousel.Id,
		&carousel.Order,
		&carousel.FilePath,
		&carousel.FileSize,
		&carousel.FileType,
		&carousel.UploadedAt,
	)

	if err != nil {
		return nil, err
	}

	return &carousel, nil
}
