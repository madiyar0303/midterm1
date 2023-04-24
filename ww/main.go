package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

// Define product struct
type Product struct {
    ID       int
    Name     string
    Price    float64
    Rating   float64
    Comments []Comment
}

// Define comment struct
type Comment struct {
    ID       int
    Username string
    Text     string
}

// Define function to filter products by price and rating
func filterProducts(products []Product, minPrice float64, maxPrice float64, minRating float64) []Product {
    var filteredProducts []Product
    for _, product := range products {
        if product.Price >= minPrice && product.Price <= maxPrice && product.Rating >= minRating {
            filteredProducts = append(filteredProducts, product)
        }
    }
    return filteredProducts
}

// Define function to rate a product
func rateProduct(productID int, rating float64) error {
    // Implement logic to rate a product in the database
    return nil
}

// Define function to comment on an article
func commentOnArticle(articleID int, username string, commentText string) error {
    // Implement logic to comment on an article in the database
    return nil
}

// Define HTTP handler function for filtering products by price and rating
func filterProductsHandler(w http.ResponseWriter, r *http.Request) {
    // Get products from the database
    products := getProductsFromDatabase()

    // Get query parameters from the request
    minPrice := r.FormValue("min_price")
    maxPrice := r.FormValue("max_price")
    minRating := r.FormValue("min_rating")

    // Filter products by price and rating
    filteredProducts := filterProducts(products, minPrice, maxPrice, minRating)

    // Return filtered products as JSON response
    json.NewEncoder(w).Encode(filteredProducts)
}

// Define HTTP handler function for rating a product
func rateProductHandler(w http.ResponseWriter, r *http.Request) {
    // Get product ID and rating from the request
    productID := r.FormValue("product_id")
    rating := r.FormValue("rating")

    // Convert rating to float64
    ratingFloat, err := strconv.ParseFloat(rating, 64)
    if err != nil {
        http.Error(w, "Invalid rating value", http.StatusBadRequest)
        return
    }

    // Rate the product
    err = rateProduct(productID, ratingFloat)
    if err != nil {
        http.Error(w, "Failed to rate product", http.StatusInternalServerError)
        return
    }

    // Return success response
    w.WriteHeader(http.StatusNoContent)
}

// Define HTTP handler function for commenting on an article
func commentOnArticleHandler(w http.ResponseWriter, r *http.Request) {
    // Get article ID, username, and comment text from the request
    articleID := r.FormValue("article_id")
    username := r.FormValue("username")
    commentText := r.FormValue("comment_text")

    // Comment on the article
    err := commentOnArticle(articleID, username, commentText)
    if err != nil {
        http.Error(w, "Failed to comment on article", http.StatusInternalServerError)
        return
    }

    // Return success response
    w.WriteHeader(http.StatusNoContent)
}

// Define function to get products from the database
func getProductsFromDatabase() []Product {
    // Implement logic to get products from the database
    var products []Product
    // ...
    return products
}

var tpl *template.Template
var db *sql.DB
type User struct {
	Name, Surname string
	Age           uint16
	Happiness     float64
	Hobbies       []string
}

func home_page(w http.ResponseWriter, r *http.Request) {
	madi := User{Name: "Madi", Surname: "Tetsuya", Age: 20, Happiness: 0.2, Hobbies: []string{"Games", "Movies", "Music"}}
	//fmt.Fprintf(w, madi.getAllInfo())
	tmpl, _ := template.ParseFiles("templates/home_page.html")
	tmpl.Execute(w, madi)
}

func main() {
	
	tpl, _ = template.ParseGlob("templates/*.html")
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/", home_page)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.ListenAndServe("localhost:3333", nil)
}
// loginHandler serves form for users to login with
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// loginAuthHandler authenticates user login
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password:", password)
	// retrieve password from db to compare (hash) with user supplied password's hash
	var hash string
	stmt := "SELECT Hash FROM bcrypt WHERE Username = ?"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&hash)
	fmt.Println("hash from db:", hash)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// returns nill on succcess
	if err == nil {
		fmt.Fprint(w, "You have successfully logged in :)")
		return
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
}

// registerHandler serves form for registring new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// registerAuthHandler creates new user in database
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	/*
		1. check username criteria
		2. check password criteria
		3. check if username is already exists in database
		4. create bcrypt hash from password
		5. insert username and password hash in database
		(email validation will be in another video)
	*/
	fmt.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	// check username for only alphaNumeric characters
	var nameAlphaNumeric = true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			nameAlphaNumeric = false
		}
	}
	// check username pswdLength
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	// check password criteria
	password := r.FormValue("password")
	fmt.Println("password:", password, "\npswdLength:", len(password))
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}
	// check if username already exists for availability
	stmt := "SELECT UserID FROM bcrypt WHERE username = ?"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}
	// create hash from password
	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))
	// func (db *DB) Prepare(query string) (*Stmt, error)
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO bcrypt (Username, Hash) VALUES (?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()
	var result sql.Result
	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err = insertStmt.Exec(username, hash)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Fprint(w, "congrats, your account has been successfully created")
}
