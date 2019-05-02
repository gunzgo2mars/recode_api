package delete

import (
	"html"
	"main/Module/database"
	"main/Module/response"
	"net/http"

	"github.com/labstack/echo"
)

type (
	DestroyPayload struct {
		AccessTokenAPI string `json:"ac_token_api"`
	}
)

// # Admin API # //

func DestroyAdminPayload(context echo.Context) (err error) {

	// Ignition Start!

	db := database.IgnitionStart()

	destroyPayload := new(DestroyPayload)

	if err = context.Bind(destroyPayload); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Bind", err.Error(), true, http.StatusOK))

	}

	if destroyPayload.AccessTokenAPI == "" {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Validation", "validation : access token is empty.", true, http.StatusOK))
	}

	AccessToken := html.EscapeString(destroyPayload.AccessTokenAPI)

	stmt, stmtErr := db.Prepare("DELETE FROM `token` WHERE `token` = ?;")

	if stmtErr != nil {

		defer db.Close()

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Statement", stmtErr.Error(), true, http.StatusOK))

	}

	_, errResponse := stmt.Exec(AccessToken)

	if errResponse != nil {

		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Exec", errResponse.Error(), true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Destroy", "Admin token payload has been destroyed.", false, http.StatusOK))

}

func DeleteUserByAdmin(c echo.Context) (err error) {

	// Ignition start!
	return nil

}

func DeleteProductByAdmin(c echo.Context) (err error) {

	// Ignition start!
	return nil

}

func DeleteReviewByAdmin(c echo.Context) (err error) {

	// Ignition start!
	return nil
}

func DeletePlaceByAdmin(c echo.Context) (err error) {

	// Ignition start!
	return nil

}

func DeleteCategoryByAdmin(c echo.Context) (err error) {

	// Ignition start!
	return nil

}

func DeleteBrandByAdmin(c echo.Context) (err error) {

	// Ignition start!
	return nil

}

// # Static API # //

func DeleteReview(c echo.Context) (err error) {

	// Ignition start!
	return nil

}
