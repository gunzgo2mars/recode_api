package server

import (
	"database/sql"
	"encoding/json"
	"main/API/create"
	"main/API/read"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	User struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Type      int    `json:"type"`
	}

	Response struct {
		title   string
		message string
		err     bool
	}
)

func Signup(c echo.Context) (err error) {

	c.Response().Header().Set(echo.HeaderContentType, "application/json")

	db, dbErr := sql.Open("mysql", "dev:123456@tcp(127.0.0.1:3306)/recode")

	if dbErr != nil {

		errResponse := &Response{title: "Error", message: "Error Message", err: true}
		resTxt, err := json.Marshal(errResponse)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, resTxt)
	}

	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}

	stmt, stmtErr := db.Prepare("INSERT INTO users (email , password , firstname , lastname, type) VALUES ( ? , ? , ? , ? , ?);")

	if stmtErr != nil {
		return c.JSON(http.StatusOK, stmtErr.Error())
	}

	_, insertErr := stmt.Exec(user.Email, user.Password, user.Firstname, user.Lastname, user.Type)

	if insertErr != nil {

		errResponse := &Response{title: "Error", message: "Insert Error", err: true}
		resTxt, err := json.Marshal(errResponse)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resTxt)
	}

	defer stmt.Close()

	success := &Response{title: "Success", message: "Success Message", err: false}
	resTxt, err := json.Marshal(success)

	return c.JSON(http.StatusOK, resTxt)
}

func Authentication(c echo.Context) error {

	// Plan To Code.
	return c.JSON(http.StatusOK, "Authentication")
}

func InitServer() {

	server := echo.New()

	// -- Config Server -- //

	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	// -- Config Server -- //

	// --- Admin Router ---
	admin := server.Group("/admin/API/v1")

	// -> Post Methods
	// # Create
	admin.POST("/createProductByAdmin", create.CreateProductByAdmin)
	admin.POST("/createPlaceByAdmin", create.CreatePlaceByAdmin)
	admin.POST("/createCategoryByAdmin", create.CreateCategoryByAdmin)
	admin.POST("/createBrandByAdmin", create.CreateBrandByAdmin)
	// #Create -> // Firebase [Clound Firestore]
	admin.POST("/createAvailableItems", create.CreateAvailableItems)
	// #Create -> // Extension
	admin.POST("/createExtensionThumbnail", create.CreateProductThumbnail)
	admin.POST("/createExtensionPlaceGallery", create.CreatePlaceGallery)

	// # Read
	admin.POST("/readBarcode", read.ReadBarcode)

	// # Delete

	admin.GET("/", func(c echo.Context) error {

		return c.JSON(http.StatusOK, "Admin API Document")

	})

	admin.POST("/api/v1/signup", Signup)

	// --- Admin Router ---

	// -- Static Router ---

	static := server.Group("/API")
	// -> Post Methods
	// # create
	static.POST("/createUser", create.CreateUser)
	static.POST("/createProduct", create.CreateProduct)
	static.POST("/createReview", create.CreateReview)

	// -> Get Methods
	// # Read

	static.GET("/welcome", func(c echo.Context) error {

		return c.JSON(http.StatusOK, "Welcome to recode/api/#")

	})

	static.GET("*", func(c echo.Context) error {

		return c.JSON(http.StatusOK, "404 Page not found.")

	})

	server.GET("/", func(c echo.Context) error {

		return c.JSON(http.StatusOK, "init server")

	})

	// -> POST Methods
	// -- Static Router ---

	// Int server
	server.Start(":1234")

}
