package authenticationservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medcard-new/begening/evtvariables"
	"medcard-new/begening/structures"
	"net/http"
	"os"
	"strings"
	"time"

	"medcard-new/begening/controllers/bycrypt"
	"medcard-new/begening/controllers/handlefile"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/controllers/velidation"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx    context.Context
	client *mongo.Client
)
var DB_Url string = os.Getenv("DBURL")

func Authenticationservice() {
	// clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// clientOptions := options.Client().ApplyURI("mongodb://mas:mas@34.148.119.65:27017")
	clientOptions := options.Client().ApplyURI(evtvariables.DBUrl)
	clientG, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: ", err)
	}
	ctxG, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	// collection := client.Database("MedCard").Collection("users")
	ctx = ctxG
	client = clientG
}
func Signin(c *gin.Context) {
	var SignupStruct structures.Signin
	c.ShouldBindJSON(&SignupStruct)
	log.Printf("str %v\n", SignupStruct)

	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""

	valueStruct, err := json.Marshal(SignupStruct)
	if err != nil {
		log.Printf("Marshel Eror %v\n", err)
	}
	var DecodedSigninStruct structures.Signin
	checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "phone:password", string(valueStruct), "FieldsCheck:true,DBCheck:false", "", "")
	collection.FindOne(ctx, bson.M{"phone": SignupStruct.Phone}).Decode(&DecodedSigninStruct)
	fmt.Printf("DecodedSigninStruct: %v\n", DecodedSigninStruct)
	passwordCHeck := bycrypt.CompareHashPasswords(DecodedSigninStruct.Password, SignupStruct.Password)
	if checkPointOne != false && checkPointTwo != true && DecodedSigninStruct.Password != "" && passwordCHeck != false {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    jwtgen.GenerateToken(c, SignupStruct.Phone),
			Expires:  time.Now().Add(30 * time.Hour),
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			MaxAge: 0,
			Domain: "",
		})
		c.JSON(200, gin.H{
			"Code":       "Authorised",
			"Id":         DecodedSigninStruct.Userid,
			"Permission": DecodedSigninStruct.Permissions,
		})
		collection.FindOne(ctx, bson.M{})
	} else {
		c.JSON(400, gin.H{
			"Code": "Cannot_Authorised",
		})
	}
}
func Signup(c *gin.Context) {
	var (
		//   DecodedSigninStruct structures.Signup
		SigninStruct structures.Signup
		checkPoint   bool
	)
	c.ShouldBindJSON(&SigninStruct)
	log.Printf("Marshel Eror %v\n", SigninStruct)
	valueStruct, err := json.Marshal(SigninStruct)
	if err != nil {
		log.Printf("Marshel Eror %v\n", err)
	}
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "phone:blood:password:email:name:surname:lastname:birth:gender:disabilaties:adress:workplace", string(valueStruct), "FieldsCheck:true,DBCheck:true", "client", "")
	log.Println(checkPoint)

	if checkPointOne != false && strings.Split(SigninStruct.Password, ":")[len(strings.Split(SigninStruct.Password, ":"))-1] == "Create" {
		log.Println("adminpassed")
		primitiveid := primitive.NewObjectID().Hex()
		hashedPass, err := bycrypt.HashPassword(strings.Split(SigninStruct.Password, ":")[0])
		log.Println(strings.Split(SigninStruct.Password, ":")[0])
		if err != nil {
			log.Printf("Err Hash%v", err)
		}
		SigninStruct.Password = hashedPass
		SigninStruct.Userid = primitiveid
		SigninStruct.Permissions = "admin"
		// SigninStruct.ImgUrl = handlefile.Handlefile(c,"./static/uploadUser")
		collection.InsertOne(ctx, SigninStruct)
		//------------------------------------ Send success-----------------------------------
		c.JSON(200, gin.H{
			"Code": "Succeded",
		})
	} else if checkPointOne != false && checkPointTwo != false {
		primitiveid := primitive.NewObjectID().Hex()
		hashedPass, err := bycrypt.HashPassword(SigninStruct.Password)
		if err != nil {
			log.Printf("Err Hash%v", err)
		}
		SigninStruct.Password = hashedPass
		SigninStruct.Userid = primitiveid
		SigninStruct.Permissions = "client"
		// SigninStruct.ImgUrl = handlefile.Handlefile(c,"./static/uploadUser")
		collection.InsertOne(ctx, SigninStruct)
		//------------------------------------ Send success-----------------------------------
		c.JSON(200, gin.H{
			"Code": "Succeded",
		})
	}else if checkPointTwo == false {
		c.JSON(304, gin.H{
			"Code": "Error User already exist",
		})
	}
}
func Signout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "nul",
		Expires:  time.Now().Add(-20 * time.Hour),
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge: 0,
		Path:     "/",
	})
	c.JSON(200, gin.H{
		"Code": "Succeded",
	})
}
func LoginCheck(c *gin.Context) {
	var User structures.Signup
	c.ShouldBindJSON(&User)
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""

	CookieData := jwtgen.Velidation(c)
	type SendData struct {
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Lastname    string `json:"lastname"`
		Userid      string `bson:"_id"`
		ImgUrl      string `json:"imgurl"`
		Permissions string `json:"permissions"`
	}

	var DecodedSigninStruct SendData
	collection.FindOne(ctx, bson.M{"_id": CookieData.Id}).Decode(&DecodedSigninStruct)

	if DecodedSigninStruct.Name != "" {

		c.JSON(200, gin.H{
			"Code": "Succeded",
			"Json": DecodedSigninStruct,
		})
	} else {
		c.JSON(505, gin.H{
			"Code": "You_are_not_authorized",
		})
	}
}
func SignupDoctor(c *gin.Context) {
	var (
		// DecodedSigninStruct structures.Signup
		SignupDoctor structures.SignupDoctor
		checkPoint   bool
	)
	// """""""""get he json request from client """""""""
	jsonFM := c.Request.FormValue("json")
	files, handler, errIMG := c.Request.FormFile("img")
	// """""""""""""""""""""""check The file on existense"""""""""""""""""""""""
	if errIMG != nil {
		c.JSON(409, gin.H{
			"sttus": "NOIMGFILEEXIST",
		})
	}

	files.Seek(23, 23)
	log.Printf("File Name %s\n", handler.Filename)
	// """""""""""""""""""""bind the request data into structure"""""""""""""""""""""
	json.Unmarshal([]byte(jsonFM), &SignupDoctor)

	log.Printf("str %v\n", SignupDoctor)
	log.Printf("Marshel Eror %v\n", SignupDoctor)
	valueStruct, err := json.Marshal(SignupDoctor)
	if err != nil {
		log.Printf("Marshel Eror %v\n", err)
	}
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "phone:password:email:name:surname:lastname:position", string(valueStruct), "FieldsCheck:true,DBCheck:true", "doctor", "")
	log.Println(checkPoint)

	if checkPointOne != false && checkPointTwo != false {
		hashedPass, err := bycrypt.HashPassword(SignupDoctor.Password)
		if err != nil {
			log.Printf("Err Hash%v", err)
		}
		primitiveid := primitive.NewObjectID().Hex()
		SignupDoctor.Password = hashedPass
		SignupDoctor.Userid = primitiveid
		SignupDoctor.Permissions = "doctor"
		SignupDoctor.ImgUrl = handlefile.Handlefile(c, "./static/upload")
		//   SignupDoctor.History = append(SignupDoctor.History, structures.History{
		// 	Year: "2022-12",
		// 	Position: "jfdfdd",
		//   })
		collection.InsertOne(ctx, SignupDoctor)
	}else if checkPointTwo == false {
		c.JSON(302, gin.H{
			"Code": "Error User already exist",
		})	
	}
}
