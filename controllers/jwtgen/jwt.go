package jwtgen

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medcard-new/begening/evtvariables"
	"medcard-new/begening/structures"
	"net/http"

	// "strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClaimsTok struct {
	Phone string `json:"phone"`
	Id string `json:"id"`
	Permissions string `json:"permissions"`
	jwt.StandardClaims
}

var myKey = []byte("sekKey")

func GenerateToken(c *gin.Context,phone string) string {
	// explore he db tofind user id
	// clientOptions := options.Client().ApplyURI("mongodb://mas:mas@34.148.119.65:27017")
	// clientOptions := options.Client().ApplyURI(os.Getenv("DB_URL"))
	clientOptions := options.Client().ApplyURI(evtvariables.DBUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: ", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	// collection := client.Database("MedCard").Collection("questions")
	collection := client.Database("MedCard").Collection("users")
	var DbgetUser structures.Signup
	collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&DbgetUser)
	fmt.Printf("DbgetUser: %v\n", DbgetUser)
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(30 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &ClaimsTok{
		Phone: phone,
		Id: DbgetUser.Userid,
		Permissions: DbgetUser.Permissions,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	return tokenString
}

func Velidation(c *gin.Context) ClaimsTok {
	// We can obtain the session token from the requests cookies, which come with every request
	cookie, err := c.Request.Cookie("tokenForUthenticateUser")
	if err != nil {
		c.JSON(400,gin.H{
			"status":"COOKIE_DOES_NOT_EXIST",
		})
	}
	// Get the JWT string from the cookie
	tknStr := cookie.Value
	// Initialize a new instance of `Claims`
	claims := &ClaimsTok{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		fmt.Fprintf(c.Writer, "error 2")
	}
	if !token.Valid {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(c.Writer, "error 3")
	}

	payloadLogin := GetAccessDetails(tknStr)
	jsonStr := string(payloadLogin)
	var jwtData ClaimsTok
	json.Unmarshal([]byte(jsonStr),&jwtData)

	var ClaimsObj ClaimsTok
	ClaimsObj.Id = jwtData.Id
	ClaimsObj.Phone = jwtData.Phone
	ClaimsObj.Permissions = jwtData.Permissions

	//"""""""""""""" Online check if the user is online or not""""""""""""""

	return ClaimsObj
}
func GetAccessDetails(tokenStr string) ([]byte) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		 // check token signing method etc
		 return []byte("sekKey"), nil
	})

	if err != nil {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid{
		fmt.Println(claims)
		ClaimsTok, err := json.Marshal(token.Claims.(jwt.MapClaims))
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
		return ClaimsTok
	}else{
		return nil
	}
}