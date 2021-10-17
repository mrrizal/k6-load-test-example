package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

type UserPayload struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func generateUser() {
	var users []UserPayload
	for i := 0; i < 1000; i++ {
		var user UserPayload
		user.Username = fmt.Sprintf("username_%d", i)
		user.Password = fmt.Sprintf("password%d", i)
		user.FirstName = fmt.Sprintf("first_name_%d", i)
		user.LastName = fmt.Sprintf("last_name_%d", i)
		users = append(users, user)
	}

	usersJSON, err := json.MarshalIndent(&users, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	_ = ioutil.WriteFile("users.json", usersJSON, 0644)
}

func loadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Info("Failed load .env file")
	}
	log.Info("Successful load .env file")
}

type User struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Username      string         `json:"username" gorm:"unique;not null;default:null"`
	Password      string         `json:"password" gorm:"not null;default:null"`
	FirstName     string         `json:"first_name" gorm:"not null;default:null"`
	LastName      string         `json:"last_name" gorm:"not null;default:null"`
	ProfilePicURL string         `json:"profile_pic_url"`
}

var (
	DBConn *gorm.DB
)

func InitDatabase() {
	u, err := pq.ParseURL(Configs.DatabaseURL)

	if err != nil {
		log.Panic(err.Error())
	}

	DBConn, err = gorm.Open(postgres.Open(u))
	if err != nil {
		log.Panic(err.Error())
	}
	log.Info("Successful connected to database.")
}

type configs struct {
	DatabaseURL string
	SecretKey   string
}

var Configs configs

func GetSettings() {
	Configs.DatabaseURL = os.Getenv("DATABASE_URL")
	Configs.SecretKey = os.Getenv("SECRET_KEY")
}

func getSingedUpUsers() {
	re := regexp.MustCompile("[a-zA-Z]+\\_(\\d+)$")
	var users []User
	var result []User
	DBConn.Model(&User{}).Scan(&users)
	for _, user := range users {
		tempUser := user
		match := re.FindStringSubmatch(user.Username)
		tempUser.Password = fmt.Sprintf("password%s", match[1])
		result = append(result, tempUser)
	}
	usersJSON, err := json.MarshalIndent(&result, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	_ = ioutil.WriteFile("signed-up-users.json", usersJSON, 0644)
}

func generateJWTToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(Configs.SecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

func getJWTUserToken() {
	var users []User
	var tokens []string
	DBConn.Model(&User{}).Scan(&users)
	for _, user := range users {
		token, _ := generateJWTToken(user)
		tokens = append(tokens, token)
	}
	tokensJSON, err := json.MarshalIndent(&tokens, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	_ = ioutil.WriteFile("tokens.json", tokensJSON, 0644)
}

func main() {
	loadEnvFile()
	GetSettings()
	InitDatabase()
	genUser := flag.Bool("generate-user", false, "a bool")
	singedUpUser := flag.Bool("signed-up-user", false, "a bool")
	genJWTToken := flag.Bool("gen-jwt-token", false, "a bool")
	flag.Parse()

	if *genUser {
		generateUser()
	}

	if *singedUpUser {
		getSingedUpUsers()
	}

	if *genJWTToken {
		getJWTUserToken()
	}
}
