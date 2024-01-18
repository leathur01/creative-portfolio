package views

import (
	"net/http"

	"github.com/revel/revel"
)

func errorResponse(responseData map[string]interface{}, status int, message interface{}, c *revel.Controller) revel.Result {
	c.Response.Status = status
	responseData["error"] = message
	return c.RenderJSON(responseData)
}

func serverErrorResponse(responseData map[string]interface{}, err error, c *revel.Controller) revel.Result {
	revel.AppLog.Error(err.Error())
	message := "the server encounterd a problem"

	// The method is used in many place and I don't have the time for reactoring but I need to get the error for debugging at the moment
	// So I use vararg err... as a temporary way to get error
	// TODO: Reafactor the function to be more semantic by getting the error in a usual way
	return errorResponse(responseData, http.StatusInternalServerError, message, c)
}

func notFoundResponse(responseData map[string]interface{}, c *revel.Controller) revel.Result {
	message := "the requested resource could not be found"
	return errorResponse(responseData, http.StatusNotFound, message, c)
}

func badRequestResponse(responseData map[string]interface{}, err string, c *revel.Controller) revel.Result {
	return errorResponse(responseData, http.StatusNotFound, err, c)
}

func failedValidationResponse(responseData map[string]interface{}, errors []*revel.ValidationError, c *revel.Controller) revel.Result {
	return errorResponse(responseData, http.StatusUnprocessableEntity, errors, c)
}
