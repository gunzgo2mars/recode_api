package read

import (
	"html"
	"main/Module/crypto"
	"main/Module/database"
	"main/Module/response"
	"net/http"

	"github.com/labstack/echo"
)

type (
	Barcode struct {
		code string `json:"barcode"`
	}

	DashboardAuth struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AppAuth struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

// # Admin API # //

func ReadAuthAdmin(context echo.Context) (err error) {

	// Ignition Start !
	db := database.IgnitionStart()

	dashboardAuth := new(DashboardAuth)

	if err = context.Bind(dashboardAuth); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Server didnt receive data from client.", true, http.StatusOK))

	}

	email := html.EscapeString(dashboardAuth.Email)
	password := html.EscapeString(dashboardAuth.Password)

	var id int
	var userEmail string
	var userPassword string
	var firstname string
	var lastname string
	var typeUser int

	err = db.QueryRow("SELECT id , email , password ,  firstname , lastname , type FROM users WHERE email=? AND type=1", email).Scan(&id, &userEmail, &userPassword, &firstname, &lastname, &typeUser)

	if err != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", err.Error(), true, http.StatusOK))
	}

	if crypto.CheckPasswordHash(password, userPassword) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Incorrect password.", true, http.StatusOK))

	} else {

		accessTokenAPI := crypto.RandStringBytes(150)

		stmt, stmtErr := db.Prepare("INSERT INTO token (token , user_id , user_type) VALUES(? , ? , ?)")

		if stmtErr != nil {
			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Statement", stmtErr.Error(), true, http.StatusOK))
		}

		_, errResponse := stmt.Exec(html.EscapeString(string(accessTokenAPI)), id, typeUser)

		if errResponse != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Exec", errResponse.Error(), true, http.StatusOK))

		}

		return context.JSON(http.StatusOK, response.ResponsePayloadData("Success", "User have been access to system.", false, http.StatusOK, map[string]interface{}{
			"uid":       id,
			"email":     userEmail,
			"firstname": firstname,
			"lastname":  lastname,
			"typeUser":  typeUser,
			"token":     accessTokenAPI,
		}))

	}

}

// # Static API # //

func ReadAuthUser(context echo.Context) (err error) {

	// Ignition Start !

	db := database.IgnitionStart()

	appAuth := new(AppAuth)

	if err = context.Bind(appAuth); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

	email := html.EscapeString(appAuth.Email)
	password := html.EscapeString(appAuth.Password)

	var id int
	var userEmail string
	var userPassword string
	var firstname string
	var lastname string
	var typeUser int

	err = db.QueryRow("SELECT id , email , password ,  firstname , lastname , type FROM users WHERE email=? AND type=0", email).Scan(&id, &userEmail, &userPassword, &firstname, &lastname, &typeUser)

	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Query", err.Error(), true, http.StatusOK))

	}

	if crypto.CheckPasswordHash(password, userPassword) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Incorrect password.", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Read Auth user", "Test Read auth user router.", true, http.StatusOK))

}

func ReadContent(context echo.Context) (err error) {

	// Ignition Start!

	return context.JSON(http.StatusOK, response.ResponseMessage("Read Content", "Test Read Content Router", true, http.StatusOK))

}

func ReadBarcode(context echo.Context) (err error) {

	// Ignition Start!
	barcode := new(Barcode)

	if err := context.Bind(barcode); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Serverside didnt receive data[Barcode] from client.", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Server side have receive data[Barcode] from client", false, http.StatusOK))
}

func ReadReview(context echo.Context) (err error) {

	// Ignition Start!

	return nil

}
