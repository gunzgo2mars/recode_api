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

	// # Admin Struct # //

	AdminPayload struct {
		Email     string `json:"admin_email"`
		Password  string `json:"admin_password"`
		Firstname string `json:"admin_firstname"`
		Lastname  string `json:"admin_lastname"`
	}

	// #-- Product Struct --#

	ProductPayload struct {
		Barcode    string           `json:"barcode"`
		CategoryID string           `json:"category"`
		BrandID    string           `json:"brand"`
		Thumbnail  string           `json:"thumbnail"`
		Price      string           `json:"price"`
		Name       string           `json:"name"`
		Detail     string           `json:"detail"`
		Gallery    []ProductGallery `json:"gallery"`
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

	ProductGallery struct {
		FirstImage  string `json:"first_image"`
		SecondImage string `json:"second_image"`
		ThirdImage  string `json:"third_image"`
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
		Name     string         `json:"name"`
		Lat      float32        `json:"lat"`
		Lon      float32        `json:"lon"`
		Des      string         `json:"des"`
		Gallery  []PlaceGallery `json:"gallery"`
		APIToken string         `json:"ac_token_api"`
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

	// # -- News --#

	News struct {
		Title       string        `json:"news_title"`
		Description string        `json:"news_des"`
		Thumbnail   string        `json:"news_thumbnail"`
		Type        string        `json:"news_type"`
		Ref         string        `json:"news_ref"`
		Gallery     []NewsGallery `json:"news_gallery"`
		APIToken    string        `json:"ac_token_api"`
	}

	NewsGallery struct {
		FirstImage  string `json:"news_firstimage"`
		SecondImage string `json:"news_secondimage"`
		ThirdImage  string `json:"news_thirdimage"`
	}

	// # -- News -- #

	CreateError struct {
		Title   string
		Message string
		Error   bool
	}
)

// # Admin API #//

func CreatAdminByAdmin(context echo.Context) (err error) {

	db := database.IgnitionStart()
	admin := new(AdminPayload)

	if err = context.Bind(admin); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Data", "Error Data type from client", true, http.StatusOK))

	}

	admin_email := html.EscapeString(admin.Email)
	admin_firstname := html.EscapeString(admin.Firstname)
	admin_lastname := html.EscapeString(admin.Lastname)

	if validation.ValidateEmail(admin_email) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Email", "Error Email Data Type.", true, http.StatusOK))

	}

	if validation.ValidateLetter(admin_firstname) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Firstname", "Error Firstname Data Type.", true, http.StatusOK))

	}

	if validation.ValidateLetter(admin_lastname) != true {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Lastname", "Error Lastname Data Type.", true, http.StatusOK))

	}

	hashPassword, err := crypto.HashPassword(admin.Password)

	checkSQL := `SELECT email FROM users WHERE email=?;`

	rowCheck := db.QueryRow(checkSQL, admin_email)

	image := "init_image.jpg"

	switch err := rowCheck.Scan(&admin_email); err {
	case sql.ErrNoRows:
		stmt, stmtErr := db.Prepare("INSERT INTO users (email , password , firstname , lastname, type , profile_image) VALUES ( ? , ? , ? , ? , ? , ?);")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Statement[SQL]", stmtErr.Error(), true, http.StatusOK))
		}

		_, errRes := stmt.Exec(admin_email, hashPassword, admin_firstname, admin_lastname, 0, image)

		if errRes != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Erro Execute[SQL]", errRes.Error(), true, http.StatusOK))

		}

		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Success to insert data to database.", false, http.StatusOK))
	case nil:
		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Data is collision", true, http.StatusOK))
	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))
	}

}

func CreateNewsByAdmin(context echo.Context) (err error) {

	newsStruct := new(News)

	db := database.IgnitionStart()

	if err := context.Bind(newsStruct); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

	verifyTokenAPI := `SELECT token FROM token WHERE token=?`

	checkToken := db.QueryRow(verifyTokenAPI, newsStruct.APIToken)

	switch err := checkToken.Scan(&verifyTokenAPI); err {

	case sql.ErrNoRows:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", sql.ErrNoRows.Error(), true, http.StatusOK))

	case nil:

		newsGalleryJson, _ := json.Marshal(newsStruct.Gallery)

		stmt, stmtErr := db.Prepare("INSERT INTO news (title , des , thumbnail , type , ref , gallery) VALUES(? , ? , ? , ? , ? , ?);")

		if stmtErr != nil {
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Statement", stmtErr.Error(), true, http.StatusOK))
		}

		_, errExec := stmt.Exec(html.EscapeString(newsStruct.Title), html.EscapeString(newsStruct.Description), html.EscapeString(newsStruct.Thumbnail), newsStruct.Type, html.EscapeString(newsStruct.Ref), newsGalleryJson)

		if errExec != nil {
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Execute", errExec.Error(), true, http.StatusOK))
		}

		defer db.Close()

		return context.JSON(http.StatusOK, response.ResponseMessage("Success", "News payload has been insert to database.", false, http.StatusOK))

	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

}

func CreateProductByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	productStruct := new(ProductPayload)

	db := database.IgnitionStart()

	if err := context.Bind(productStruct); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

	checkExistBarcodeSQL := `SELECT barcode FROM products WHERE barcode=?;`

	rowCheckBarCode := db.QueryRow(checkExistBarcodeSQL, productStruct.Barcode)

	productNameJson, _ := json.Marshal(productStruct.Name)
	productPriceJson, _ := json.Marshal(productStruct.Price)
	productDetailJson, _ := json.Marshal(productStruct.Detail)

	productGalleryJson, _ := json.Marshal(productStruct.Gallery)

	category, err := strconv.Atoi(productStruct.CategoryID)
	brand, err := strconv.Atoi(productStruct.BrandID)

	t := time.Now()

	thumbnail := t.Format("20060102150405") + productStruct.Thumbnail

	switch err := rowCheckBarCode.Scan(&checkExistBarcodeSQL); err {

	case sql.ErrNoRows:
		stmt, stmtErr := db.Prepare("INSERT INTO products (barcode , name , detail , price , thumbnail , gallery , categories_id , brands_id) VALUES(? , ? , ? , ? , ? , ?, ? , ?)")

		if stmtErr != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Statement", stmtErr.Error(), true, http.StatusOK))
		}

		_, errRes := stmt.Exec(productStruct.Barcode, productNameJson, productDetailJson, productPriceJson, thumbnail, productGalleryJson, category, brand)

		if errRes != nil {

			defer db.Close()
			return context.JSON(http.StatusOK, response.ResponseMessage("Error Execute", errRes.Error(), true, http.StatusOK))

		}

		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Success to insert data to database.", false, http.StatusOK))
	case nil:
		defer db.Close()
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", "Data is collision", true, http.StatusOK))
	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))
	}

}

func CreatePlaceByAdmin(context echo.Context) (err error) {

	// Ignition Start!

	placeStruct := new(Place)

	db := database.IgnitionStart()

	if err := context.Bind(placeStruct); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Param", err.Error(), true, http.StatusOK))

	}

	name := html.EscapeString(placeStruct.Name)
	des := html.EscapeString(placeStruct.Des)

	placeGalleryJson, _ := json.Marshal(placeStruct.Gallery)

	verifyToken := `SELECT token FROM token WHERE token=?`

	checkToken := db.QueryRow(verifyToken, placeStruct.APIToken)

	switch err := checkToken.Scan(&verifyToken); err {

	case sql.ErrNoRows:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error SQL", sql.ErrNoRows.Error(), true, http.StatusOK))

	case nil:
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
			return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

		}

	default:
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

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

func CreateNewsThumbnail(context echo.Context) (err error) {

	file, errFile := context.FormFile("news_thumbnailfile")

	if errFile != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error File", errFile.Error()+"[News API Thumbnail]", true, http.StatusOK))

	}

	src, errSrc := file.Open()

	if errSrc != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Source", errSrc.Error()+"[News API Thumbnail]", true, http.StatusOK))

	}

	defer src.Close()

	dst, errDst := os.Create("./payload/news_images/thumbnail/" + file.Filename)

	if errDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errDst.Error()+"[News API Thumbnail]", true, http.StatusOK))

	}

	defer dst.Close()

	if _, errIO := io.Copy(dst, src); errIO != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error I/O", errIO.Error()+"[News API Thumbnail]", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "News thumbnail has been uploaded.", false, http.StatusOK))

}

func CreatePlaceThumbnail(context echo.Context) (err error) {

	file, errFile := context.FormFile("place_thumbnail")

	if errFile != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error File", errFile.Error(), true, http.StatusOK))

	}

	src, errSrc := file.Open()

	if errSrc != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Source", errSrc.Error(), true, http.StatusOK))

	}

	defer src.Close()

	dst, errDst := os.Create("./payload/place_image/" + file.Filename)

	if errDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errDst.Error(), true, http.StatusOK))

	}

	if _, errIO := io.Copy(dst, src); errIO != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Errro I/O", errIO.Error(), true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Place thumbnail has been uploaded.", false, http.StatusOK))

}

func CreateProductThumbnail(context echo.Context) (err error) {

	file, err := context.FormFile("thumbnail_file")

	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

	src, err := file.Open()

	if err != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))
	}

	defer src.Close()

	t := time.Now()

	dst, err := os.Create("./payload/product_images/thumbnail/" + t.Format("20060102150405") + file.Filename)
	if err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", err.Error(), true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Uploaded!", false, http.StatusOK))

}

func CreateNewsGallery(context echo.Context) (err error) {

	firstFile, firstErr := context.FormFile("news_firstfile")
	secondFile, secondErr := context.FormFile("news_secondfile")
	thirdFile, thirdErr := context.FormFile("news_thirdfile")

	if firstErr != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", firstErr.Error()+"[News API Gallery]", true, http.StatusOK))

	}

	if secondErr != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", firstErr.Error()+"[News API Gallery]", true, http.StatusOK))

	}

	if thirdErr != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error", firstErr.Error()+"[News API Gallery]", true, http.StatusOK))

	}

	firstSrc, firstSrcErr := firstFile.Open()
	secondSrc, secondSrcErr := secondFile.Open()
	thirdSrc, thirdSrcErr := thirdFile.Open()

	if firstSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error File", firstSrcErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}
	if secondSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error File", secondSrcErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}
	if thirdSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error File", thirdSrcErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}

	defer firstSrc.Close()
	defer secondSrc.Close()
	defer thirdSrc.Close()

	firstDst, firstDstErr := os.Create("./payload/news_images/gallery/" + firstFile.Filename)
	secondDst, secondDstErr := os.Create("./payload/news_images/gallery/" + secondFile.Filename)
	thirdDst, thirdDstErr := os.Create("./payload/news_images/gallery/" + thirdFile.Filename)

	if firstDstErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", firstDstErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}
	if secondDstErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", secondDstErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}
	if thirdDstErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", thirdDstErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}

	defer firstDst.Close()
	defer secondDst.Close()
	defer thirdDst.Close()

	if _, firstIoErr := io.Copy(firstDst, firstSrc); firstIoErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", firstIoErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}

	if _, secondIoErr := io.Copy(secondDst, secondSrc); secondIoErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", secondIoErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}

	if _, thirdIoErr := io.Copy(thirdDst, thirdSrc); thirdIoErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", thirdIoErr.Error()+"[News API Gallery]", true, http.StatusOK))
	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "News gallery has been uploaded.", false, http.StatusOK))

}

func CreatePlaceGallery(context echo.Context) (err error) {

	firstFile, firstErr := context.FormFile("first_image")
	secondFile, secondErr := context.FormFile("second_image")
	thirdfile, thirdErr := context.FormFile("third_image")

	if firstErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", firstErr.Error()+"[First]", true, http.StatusOK))
	}

	if secondErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", secondErr.Error()+"[Second]", true, http.StatusOK))
	}

	if thirdErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Open file", thirdErr.Error()+"[Third]", true, http.StatusOK))
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

	firstFileDst, errFirstDst := os.Create("./payload/place_images/gallery/" + firstFile.Filename)
	secondFileDst, errSecondDst := os.Create("./payload/place_images/gallery/" + secondFile.Filename)
	thirdFileDst, errThirdDst := os.Create("./payload/place_images/gallery/" + thirdfile.Filename)

	if errFirstDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errFirstDst.Error(), true, http.StatusOK))

	}

	defer firstFileDst.Close()

	if errSecondDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errSecondDst.Error(), true, http.StatusOK))

	}

	defer secondFileDst.Close()

	if errThirdDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errThirdDst.Error(), true, http.StatusOK))

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

	firstFile, firstErr := context.FormFile("first_image")
	secondFile, secondErr := context.FormFile("second_image")
	thirdfile, thirdErr := context.FormFile("third_image")

	if firstErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", firstErr.Error()+"[First]", true, http.StatusOK))
	}

	if secondErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", secondErr.Error()+"[Second]", true, http.StatusOK))
	}

	if thirdErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error Open file", thirdErr.Error()+"[Third]", true, http.StatusOK))
	}

	firstSrc, firstSrcErr := firstFile.Open()
	secondSrc, secondSrcErr := secondFile.Open()
	thirdSrc, thirdSrcErr := thirdfile.Open()

	if firstSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", firstSrcErr.Error()+"[First]", true, http.StatusOK))
	}

	defer firstSrc.Close()

	if secondSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", secondSrcErr.Error()+"[Second]", true, http.StatusOK))
	}

	defer secondSrc.Close()

	if thirdSrcErr != nil {
		return context.JSON(http.StatusOK, response.ResponseMessage("Error open file", thirdSrcErr.Error()+"[Third]", true, http.StatusOK))
	}

	defer thirdSrc.Close()

	firstFileDst, errFirstDst := os.Create("./payload/product_images/gallery/" + firstFile.Filename)
	secondFileDst, errSecondDst := os.Create("./payload/product_images/gallery/" + secondFile.Filename)
	thirdFileDst, errThirdDst := os.Create("./payload/product_images/gallery/" + thirdfile.Filename)

	if errFirstDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errFirstDst.Error(), true, http.StatusOK))

	}

	defer firstFileDst.Close()

	if errSecondDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errSecondDst.Error(), true, http.StatusOK))

	}

	defer secondFileDst.Close()

	if errThirdDst != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error Destination", errThirdDst.Error(), true, http.StatusOK))

	}

	defer thirdFileDst.Close()

	if _, err = io.Copy(firstFileDst, firstSrc); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", err.Error()+"[First]", true, http.StatusOK))

	}

	if _, err = io.Copy(secondFileDst, secondSrc); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", err.Error()+"[Second]", true, http.StatusOK))

	}

	if _, err = io.Copy(thirdFileDst, thirdSrc); err != nil {

		return context.JSON(http.StatusOK, response.ResponseMessage("Error IO", err.Error()+"[Third]", true, http.StatusOK))

	}

	return context.JSON(http.StatusOK, response.ResponseMessage("Success", "Uploaded!", false, http.StatusOK))

}
