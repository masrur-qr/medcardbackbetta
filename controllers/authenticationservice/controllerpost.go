package authenticationservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medcard-new/begening/evtvariables"
	"medcard-new/begening/structures"
	"net/http"
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

func Authenticationservice() {
	clientOptions := options.Client().ApplyURI(evtvariables.DBUrl)

	clientG, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: 1", err)
	}
	ctxG, _ := context.WithTimeout(context.Background(), 15*time.Minute)
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
			Name:     "tokenForUthenticateUser",
			Value:    jwtgen.GenerateToken(c, SignupStruct.Phone),
			Expires:  time.Now().Add(30 * time.Hour),
			HttpOnly: false,
			Secure:   false,
			// Path:     "/",
			// SameSite: http.SameSiteNoneMode,
			MaxAge: 0,
			Domain: ".console.academy",
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

func SignupAdmin() {
	
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	primitiveid := primitive.NewObjectID().Hex()
	hashedPass, err := bycrypt.HashPassword(strings.Split("admin123", ":")[0])
	if err != nil {
		log.Printf("Err Hash%v", err)
	}
	var SigninStruct structures.Signup = structures.Signup{
		Userid: primitiveid,
		Phone: "+992934009886",
		Password:hashedPass ,
		Email: "masrur.gdsc@gmail.com",
		Name: "Masrur",
		Surname: "Qurbonmamadov",
		Lastname: "",
		Birth: "31.03.2003",
		Gender: "Male",
		Disabilaties: "No",
		Blood: "",
		Adress: "",
		Workplace: "",
		ImgUrl: "",
		Permissions: "admin",
	}
	collection.InsertOne(ctx, SigninStruct)
}
func Signup(c *gin.Context) {
	var (
		SigninStruct structures.Signup
		allowedHosts = c.GetHeader("Origin")
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
	checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "phone:blood:password:name:surname:lastname:birth:gender:disabilaties:adress:workplace", string(valueStruct), "FieldsCheck:true,DBCheck:true", "client", "")

	if allowedHosts == evtvariables.IpUrl {
		if checkPointOne != false && strings.Split(SigninStruct.Password, ":")[len(strings.Split(SigninStruct.Password, ":"))-1] == "Create" {
			primitiveid := primitive.NewObjectID().Hex()
			hashedPass, err := bycrypt.HashPassword(strings.Split(SigninStruct.Password, ":")[0])
			log.Println(strings.Split(SigninStruct.Password, ":")[0])
			if err != nil {
				log.Printf("Err Hash%v", err)
			}
			SigninStruct.Password = hashedPass
			SigninStruct.Userid = primitiveid
			SigninStruct.Permissions = "admin"
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
			collection.InsertOne(ctx, SigninStruct)
			//------------------------------------ Send success-----------------------------------
			c.JSON(200, gin.H{
				"Code": "Succeded",
			})
		} else if checkPointTwo == false {
			c.JSON(404, gin.H{
				"Code": "Error User already exist OR this phone numbers are already taken",
			})
		}
	} else {
		c.JSON(505, gin.H{
			"Message": "Unknown Host",
		})
	}

}
func Signout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "expire",
		Value:    "nul",
		Expires:  time.Now().Add(-20 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		// SameSite: http.SameSiteNoneMode,
		MaxAge: 0,
		Path:   "/",
		Domain: ".medcard.space",
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
		allowedHosts = c.GetHeader("Origin")
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

	valueStruct, err := json.Marshal(SignupDoctor)
	if err != nil {
		log.Printf("Marshel Eror %v\n", err)
	}
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "phone:password:name:surname:lastname:position", string(valueStruct), "FieldsCheck:true,DBCheck:true", "doctor", "")
	log.Println(checkPoint)

	if allowedHosts == evtvariables.IpUrl {
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

			collection.InsertOne(ctx, SignupDoctor)
		} else if checkPointTwo == false {
			c.JSON(304, gin.H{
				"Code": "Error User already exist OR this phone numbers are already taken",
			})
		}
	} else {
		c.JSON(505, gin.H{
			"Message": "Unknown Host",
		})
	}
}

func ResetPassword(c *gin.Context) {
	var (
		allowedHosts = c.GetHeader("Origin")
		ResetStruct  structures.Reset
		UserDecode   structures.Reset
	)

	c.ShouldBindJSON(&ResetStruct)
	if allowedHosts == evtvariables.IpUrl {
		if ResetStruct.Phone != "" && ResetStruct.NewPassword == ResetStruct.Password {
			Authenticationservice()
			connection := client.Database("MedCard").Collection("users")
			connection.FindOne(ctx, bson.M{"phone": ResetStruct.Phone}).Decode(&UserDecode)
			// ? Updating the field
			HashedPass, _ := bycrypt.HashPassword(ResetStruct.NewPassword)
			_, err := connection.UpdateMany(ctx, bson.M{
				"phone": ResetStruct.Phone,
			},
				bson.D{
					{"$set", bson.M{"password": HashedPass}},
				},
			)
			c.JSON(200, gin.H{
				"Message": "Request successfully handled",
			})
			if err != nil {
				c.JSON(400, gin.H{
					"Message": "User noy found",
				})
			}
		}
	} else {
		c.JSON(505, gin.H{
			"Message": "Unknown Host",
		})
	}
}
