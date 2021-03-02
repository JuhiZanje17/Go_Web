package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = getCredentials("HOST")
	password = getCredentials("PASSWORD")
	dbname   = getCredentials("DB")
)
var db *gorm.DB

type User struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `json:"name"`
	PWD  string `json:"pwd"`
	Role string `json:"role"`
}

var secKey = []byte("my_secret_key")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func getCredentials(key string) string {

	// load .env file
	err := godotenv.Load("../../jwtCred.env")

	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}

	return os.Getenv(key)
}

func createTable() {
	var err error
	psqlInfo := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", password, host, dbname)

	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	sqldb, err := db.DB()
	err = sqldb.Ping()
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&User{})

}

func login(w http.ResponseWriter, r *http.Request) {

	var cred Credentials
	var u User
	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("name = ? ", cred.Username).First(&u)

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(u.PWD), []byte(cred.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		fmt.Println("matched")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Role:     u.Role,
		Username: u.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(secKey)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.Write([]byte("completed"))

}

func Welcome(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("claims=", claims)

	if claims.Role == "admin" {
		http.Redirect(w, r, "/admin", 301)
		return
	} else {
		http.Redirect(w, r, "/user", 301)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func adminPage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Welcome Admin")))

}
func userPage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Welcome User")))

}

func signupPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PWD), 8)
	user.PWD = string(hashedPassword)
	fmt.Println(user)
	db.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func main() {

	boolCreate := flag.Bool("create", false, "boolean")
	flag.Parse()

	if *boolCreate {
		createTable()
		fmt.Println("done")
	} else {
		fmt.Println("not inserted")
	}
	http.HandleFunc("/login", login)
	http.HandleFunc("/admin", adminPage)
	http.HandleFunc("/user", userPage)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/signup", signupPage)
	fmt.Println("listening on 8000")

	http.ListenAndServe(":8000", nil)
}
