package read

import (
	"main/Module/response"
	"net/http"

	"github.com/labstack/echo"
)

type (
	Barcode struct {
		code string `json:"barcode"`
	}
)

// # Admin API # //

func ReadAuthAdmin(context echo.Context) (err error) {

	// Ignition Start !

	return nil

}

// # Static API # //

func ReadAuthUser(context echo.Context) (err error) {

	// Ignition Start !

	return context.JSON(http.StatusOK, response.ResponseMessage("Read Auth user", "Test Read auth user router.", true, http.StatusOK))

}

func ReadContent(context echo.Context) (err error) {

	// Ignition Start!

	return context.JSON(http.StatusOK, response.ResponseMessage("Read Content", "Test Read Content Router", true, http.StatusOK))

}

func ReadBarcode(context echo.Context) (err error) {

	// Ignition Start!

	return nil
}

func ReadReview(context echo.Context) (err error) {

	// Ignition Start!

	return nil

}
