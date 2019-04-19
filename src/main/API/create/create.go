package create

import (
	"context"
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

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"

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

	// #-- Place --#

	Place struct {
		Name    string         `json:"name"`
		Lat     float32        `json:"lat"`
		Lon     float32        `json:"lon"`
		Des     string         `json:"des"`
		Gallery []PlaceGallery `json:"gallery"`
	}

	PlaceGallery struct {
		FirstImage  string `json:"first_image"`
		SecondImage string `json:"second_image"`
		ThirdImage  string `json:"third_image"`
	}

	AvailableItems struct {
		PlaceID   string `json:"place_id"`
		ProductID string `json:"product_id"`
	}

	// #-- Place --#

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

	t := time.Now()

	thumbnail := t.Format("20060102150405") + productStruct.Thumbnail

	switch err := rowCheckBarCode.Scan(&checkExistBarcodeSQL); err {

	case sql.ErrNoRows:
		stmt, stmtErr := db.Prepare("INSERT INTO products (barcode , name , detail , price , thumbnail , categories_id , brands_id) VALUES(? , ? , ? , ? , ? , ? , ?)")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Failed to insert data.", true, http.StatusOK))
		}

		_, errRes := stmt.Exec(productStruct.Barcode, productNameJson, productDetailJson, productPriceJson, thumbnail, category, brand)

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

	placeStruct := new(Place)

	db := database.IgnitionStart()

	if err := context.Bind(placeStruct); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Param", "[POST ERROR PARAMS]", true, http.StatusOK))

	}

	name := html.EscapeString(placeStruct.Name)
	des := html.EscapeString(placeStruct.Des)

	if validation.ValidateLetter(name) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Data", "Error Validate text is not a letter", true, http.StatusOK))

	}

	placeGalleryJson, _ := json.Marshal(placeStruct.Gallery)

	checkSQL := `SELECT name FROM places WHERE name=?;`

	rowCheck := db.QueryRow(checkSQL, name)

	switch err := rowCheck.Scan(&name); err {

	case sql.ErrNoRows:
		stmt, stmtErr := db.Prepare("INSERT INTO places ( name , lat , lon , des , gallery ) VALUES ( ? , ? , ? , ? , ? )")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Failed to insert data.", true, http.StatusOK))

		}

		_, errRes := stmt.Exec(name, placeStruct.Lat, placeStruct.Lon, des, placeGalleryJson)

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

func CreateAvailableItems(c echo.Context) (err error) {

	availableItems := new(AvailableItems)

	if err := c.Bind(availableItems); err != nil {

		return c.JSON(http.StatusOK, response.ResponseMessage("Error Param", "[POST ERROR PARAMS]", true, http.StatusOK))

	}

	sa := option.WithCredentialsFile("./gcloud_key.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)

	client, err := app.Firestore(context.Background())

	if err != nil {

		return c.JSON(http.StatusOK, response.ResponseMessage("Error Firebase", "Failed to connect to firebase with credentials file.", true, http.StatusOK))

	}

	defer client.Close()

	return c.JSON(http.StatusOK, response.ResponseMessage("Success", "Successfuly connection [Firebase/Cloud Firestore]", false, http.StatusOK))

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

func CreatePlaceGallery(context echo.Context) (err error) {

	firstFile, firstErr := context.FormFile("first_image")
	secondFile, secondErr := context.FormFile("secode_image")
	thirdfile, thirdErr := context.FormFile("third_image")

	if firstErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", "Failed to receive file [Image : 1].", true, http.StatusOK))
	}

	if secondErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", "Failed to receive file [Image : 2]", true, http.StatusOK))
	}

	if thirdErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Open file", "Failed to receive file [Image : 3]", true, http.StatusOK))
	}

	firstSrc, firstSrcErr := firstFile.Open()
	secondSrc, secondSrcErr := secondFile.Open()
	thirdSrc, thirdSrcErr := thirdfile.Open()

	if firstSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", "Failed to open file [Image : 1]", true, http.StatusOK))
	}

	defer firstSrc.Close()

	if secondSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", "Failed to open file [Image : 2]", true, http.StatusOK))
	}

	defer secondSrc.Close()

	if thirdSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", "Failed to open file [Image : 3]", true, http.StatusOK))
	}

	defer thirdSrc.Close()

	firstFileDst, errFirstDst := os.Create("payload/place_image/" + firstFile.Filename)
	secondFileDst, errSecondDst := os.Create("payload/place_image/" + secondFile.Filename)
	thirdFileDst, errThirdDst := os.Create("payload/place_image/" + thirdfile.Filename)

	if errFirstDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("", "", true, http.StatusOK))

	}

	defer firstFileDst.Close()

	if errSecondDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("", "", true, http.StatusOK))

	}

	defer secondFileDst.Close()

	if errThirdDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("", "", true, http.StatusOK))

	}

	defer thirdFileDst.Close()

	if _, err = io.Copy(firstFileDst, firstSrc); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", "Failed to copy src file to dst [Image : 1]", true, http.StatusOK))

	}

	if _, err = io.Copy(secondFileDst, secondSrc); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", "Failed to copy src file to dst [Image : 2]", true, http.StatusOK))

	}

	if _, err = io.Copy(thirdFileDst, thirdSrc); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", "Failed to copy src file to dst [Image : 3]", true, http.StatusOK))

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
