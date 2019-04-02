package create

import (
	"database/sql"
	"encoding/json"
	"html"
	"io"
	"main/Module/crypto"
	"main/Module/database"
	"main/Module/response"
	"main/Module/validation"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type (
	User struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Type      int    `json:"type"`
	}

	// #-- Product Struct --#

	ProductPayload struct {
		Barcode    string          `json:"barcode"`
		CategoryID string          `json:"category"`
		BrandID    string          `json:"brand"`
		Thumbnail  string          `json:"thumbnail"`
		Price      []ProductPrice  `json:"price"`
		Name       []ProductName   `json:"name"`
		Detail     []ProductDetail `json:"detail"`
	}

	ProductName struct {
		TH string `json:"th"`
		EN string `json:"en"`
		CN string `json:"cn"`
	}

	ProductPrice struct {
		TH string `json:"th"`
		EN string `json:"en"`
		CN string `json:"cn"`
	}

	ProductDetail struct {
		TH string `json:"th"`
		EN string `json:"en"`
		CN string `json:"cn"`
	}

	// #-- Product --#

	Category struct {
		CategoryName string `json:"categoryname"`
	}

	// #-- Review --#

	ReviewPayload struct {
	}

	// #-- Review --#
	CreateError struct {
		Title   string
		Message string
		Error   bool
	}
)

// # Admin API #//

func CreateProductByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	productStruct := new(ProductPayload)

	db := database.IgnitionStart()

	if err := context.Bind(productStruct); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "[Params data error : POST]", true, http.StatusOK))

	}

	checkExistBarcodeSQL := `SELECT barcode FROM products WHERE barcode=?;`

	rowCheckBarCode := db.QueryRow(checkExistBarcodeSQL, productStruct.Barcode)

	productNameJson, _ := json.Marshal(productStruct.Name)
	productPriceJson, _ := json.Marshal(productStruct.Price)
	productDetailJson, _ := json.Marshal(productStruct.Detail)

	category, err := strconv.Atoi(productStruct.CategoryID)
	brand, err := strconv.Atoi(productStruct.BrandID)

	switch err := rowCheckBarCode.Scan(&checkExistBarcodeSQL); err {

	case sql.ErrNoRows:
		stmt, stmtErr := db.Prepare("INSERT INTO products (barcode , name , detail , price , thumbnail , categories_id , brands_id) VALUES(? , ? , ? , ? , ? , ? , ?)")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Failed to insert data.", true, http.StatusOK))
		}

		_, errRes := stmt.Exec(productStruct.Barcode, productNameJson, productDetailJson, productPriceJson, productStruct.Thumbnail, category, brand)

		if errRes != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", errRes.Error(), true, http.StatusOK))

		}

		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Success to insert data to database.", false, http.StatusOK))
	case nil:
		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Data is collision", true, http.StatusOK))
	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Switch", "Test", true, http.StatusOK))
	}

}

func CreateProductThumbnailByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	return nil

}

func CreatePlaceByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	return nil

}

func CreateCategoryByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	categoryStruct := new(Category)

	db := database.IgnitionStart()

	if err := context.Bind(categoryStruct); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Param", "[POST ERROR PARAMS]", true, http.StatusOK))

	}

	name := html.EscapeString(categoryStruct.CategoryName)

	if validation.ValidateLetter(name) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error data string", "", true, http.StatusOK))

	}

	checkSQL := `SELECT name FROM categories WHERE name=?;`

	rowCheck := db.QueryRow(checkSQL, name)

	switch err := rowCheck.Scan(&name); err {

	case sql.ErrNoRows:

		stmt, stmtErr := db.Prepare("INSERT INTO categories (name) VALUES (?);")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Failed to insert data.", true, http.StatusOK))

		}

		_, errRes := stmt.Exec(name)

		if errRes != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Erro SQL", "Failed to Execute SQL", true, http.StatusOK))

		}

		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Success to insert data to database.", false, http.StatusOK))

	case nil:
		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Data is collision", true, http.StatusOK))
	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Switch", "Test", true, http.StatusOK))

	}

}

func CreateBrandByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	return nil

}

// # Static API # //

func CreateUser(context echo.Context) (err error) {

	db := database.IgnitionStart()
	user := new(User)

	if err = context.Bind(user); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Data", "Error Data type from client", true, http.StatusOK))

	}

	email := html.EscapeString(user.Email)
	firstname := html.EscapeString(user.Firstname)
	lastname := html.EscapeString(user.Lastname)

	if validation.ValidateEmail(email) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Email", "Error Email Data Type.", true, http.StatusOK))

	}

	if validation.ValidateLetter(firstname) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Firstname", "Error Firstname Data Type.", true, http.StatusOK))

	}

	if validation.ValidateLetter(lastname) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Lastname", "Error Lastname Data Type.", true, http.StatusOK))

	}

	hashPassword, err := crypto.HashPassword(user.Password)

	checkSQL := `SELECT email FROM users WHERE email=?;`

	rowCheck := db.QueryRow(checkSQL, email)

	switch err := rowCheck.Scan(&email); err {
	case sql.ErrNoRows:
		stmt, stmtErr := db.Prepare("INSERT INTO users (email , password , firstname , lastname, type) VALUES ( ? , ? , ? , ? , ?);")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Failed to insert data.", true, http.StatusOK))
		}

		_, errRes := stmt.Exec(user.Email, hashPassword, user.Firstname, user.Lastname, user.Type)

		if errRes != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Erro SQL", "Failed to Execute SQL", true, http.StatusOK))

		}

		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Success to insert data to database.", false, http.StatusOK))
	case nil:
		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Data is collision", true, http.StatusOK))
	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Switch", "Test", true, http.StatusOK))
	}

}

func CreateProduct(context echo.Context) (err error) {

	// Ignition Start!

	product := new(ProductPayload)

	if err = context.Bind(product); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error bind data param", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, product)

}

func CreateReview(context echo.Context) error {

	// Ignition Start!
	return nil

}

func CreateAnalysisData(context echo.Context) error {

	// Ignition Start!

	return nil

}

// # Extension Create API # //

func CreateProductThumbnail(context echo.Context) (err error) {

	file, err := context.FormFile("file")

	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 1", true, http.StatusOK))

	}

	src, err := file.Open()

	if err != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 2", true, http.StatusOK))
	}

	defer src.Close()

	t := time.Now()

	dst, err := os.Create(t.Format("20060102150405") + file.Filename)
	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 3", true, http.StatusOK))

	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 4", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Uploaded!", false, http.StatusOK))

}

func CreateProductGallery(context echo.Context) (err error) {

	return nil

}

func TestCreateImage(context echo.Context) (err error) {

	file, err := context.FormFile("file")

	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 1", true, http.StatusOK))

	}

	src, err := file.Open()

	if err != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 2", true, http.StatusOK))
	}

	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 3", true, http.StatusOK))

	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", "Error Receive formfile image 4", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Uploaded!", false, http.StatusOK))

}
