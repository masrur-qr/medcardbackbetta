package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medcard-new/begening/controllers/bycrypt"
	"medcard-new/begening/controllers/handlefile"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/controllers/velidation"
	"medcard-new/begening/evtvariables"
	"medcard-new/begening/structures"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GlobeStruct struct {
	QuestionId          string `json:"questionid" bson:"_id"`
	QuestionsText       string `json:"questiontext"`
	QuestionsTitle      string `json:"questiontitle"`
	QuestionsAuthorName string `json:"questionauthorname"`
}

var (
	ctx    context.Context
	client *mongo.Client
)


func Authenticationservice() {
	clientOptions := options.Client().ApplyURI(evtvariables.DBUrl)
	clientG, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: ", err)
	}
	ctxG, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	ctx = ctxG
	client = clientG
}
func InsertQuestions(c *gin.Context) {
	var (
		Questions   GlobeStruct
		QuestionsDb GlobeStruct
	)
	c.ShouldBindJSON(&Questions)
	fmt.Printf("Questions: %v\n", Questions)
	// """"""""""""""""""""""JWT VALIDATION""""""""""""""""""""""""""
	ClaimsData := jwtgen.Velidation(c)
	log.Println(ClaimsData)
	// """"""""""""""""""""""JWT VALIDATION""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("questions")
	collection.FindOne(ctx, bson.M{"questionstitle": Questions.QuestionsTitle}).Decode(&QuestionsDb)

	valueStruct, err := json.Marshal(Questions)
	if err != nil {
		log.Printf("Marshel Error: %v\n", err)
	}

	Exist, Empty := velidation.TestTheQuestion(string(valueStruct))
	if Exist == true  && Empty == false  {
		var primitiveId = primitive.NewObjectID().Hex()
		Questions.QuestionId = primitiveId

		collection.InsertOne(ctx, Questions)
	} else {
		c.JSON(400, gin.H{
			"Code": "The Question Already exist",
		})
	}
}
func ProfileChange(c *gin.Context) {
	var (
		CheckpointPassed    bool = false
		ChangeStruct        structures.GlobeStruct
		DecodedSigninStruct structures.SignupDoctor
	)
	jsonFM := c.Request.FormValue("json")
	_, _, errIMG := c.Request.FormFile("img")
	// """""""""""""""""""""bind the request data into structure"""""""""""""""""""""
	json.Unmarshal([]byte(jsonFM), &ChangeStruct)
	fmt.Printf("ChangeStruct: %v\n", ChangeStruct)
	// """""""""""""""""""""bind the request data into structure"""""""""""""""""""""
	json.Unmarshal([]byte(jsonFM), &ChangeStruct)
	valueStruct, err := json.Marshal(ChangeStruct)
	if err != nil {
		log.Printf("err%v", err)
	}

	CookieData := jwtgen.Velidation(c)
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")

	if CookieData.Permissions == "admin" {
		var ChangeStruct structures.Admin
		json.Unmarshal([]byte(jsonFM), &ChangeStruct)
		checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "phone:email:name:surname:lastname", string(valueStruct), "FieldsCheck:true,DBCheck:true", "admin", CookieData.Id)
		log.Println(CheckpointPassed)
		collection.FindOne(ctx, bson.M{"name": ChangeStruct.Name, "surname": ChangeStruct.Surname, "permissions": CookieData.Permissions}).Decode(&DecodedSigninStruct)
		if checkPointOne != false && checkPointTwo == false {
			log.Printf("ds1%v\n", ChangeStruct)
			ChangeStruct.Userid = CookieData.Id
			ChangeStruct.Permissions = CookieData.Permissions
			if ChangeStruct.Password == "" {
				ChangeStruct.Password = DecodedSigninStruct.Password
			} else {
				fmt.Println("elsePAssdone:")
				hashedPass, err := bycrypt.HashPassword(ChangeStruct.Password)
				if err != nil {
					log.Printf("Err Hash%v", err)
				}
				ChangeStruct.Password = hashedPass
			}
			if errIMG != nil {
				ChangeStruct.ImgUrl = DecodedSigninStruct.ImgUrl
			} else {
				imgid := handlefile.Handlefile(c, "./static/upload")
				ChangeStruct.ImgUrl = imgid
			}
			_, err = collection.ReplaceOne(ctx, bson.M{"_id": CookieData.Id}, ChangeStruct)
			if err != nil {
				log.Printf("Err insert", err)
			}
			collection.InsertOne(ctx, ChangeStruct)
			c.JSON(200, gin.H{
				"Code": "Your Request Successfully Handeled",
			})
		}else if checkPointTwo == false {
			c.JSON(404, gin.H{
				"Code": "User not Found1.2",
			})
		}
	} else if CookieData.Permissions == "doctor" {
		var ChangeStruct structures.SignupDoctor
		json.Unmarshal([]byte(jsonFM), &ChangeStruct)
		log.Printf("ds2%v\n", ChangeStruct)
		checkPointOne, checkPointTwo := velidation.TestTheProfileChange(string(valueStruct),CookieData.Id)
		collection.FindOne(ctx, bson.M{"name": ChangeStruct.Name, "surname": ChangeStruct.Surname, "permissions": CookieData.Permissions}).Decode(&DecodedSigninStruct)

		if checkPointOne  && checkPointTwo  {
			log.Printf("ds1%v\n", ChangeStruct)
			ChangeStruct.Userid = CookieData.Id
			ChangeStruct.Permissions = CookieData.Permissions
			if ChangeStruct.Password == "" {
				ChangeStruct.Password = DecodedSigninStruct.Password
			} else {
				fmt.Println("elsePAssdone:")
				hashedPass, err := bycrypt.HashPassword(ChangeStruct.Password)
				if err != nil {
					log.Printf("Err Hash%v", err)
				}
				ChangeStruct.Password = hashedPass
			}
			if errIMG != nil {
				log.Printf("123%v\n", ChangeStruct)
				ChangeStruct.ImgUrl = DecodedSigninStruct.ImgUrl
			} else {
				log.Printf("456%v\n", ChangeStruct)
				imgid := handlefile.Handlefile(c, "./static/upload")
				ChangeStruct.ImgUrl = imgid
			}
			_, err = collection.ReplaceOne(ctx, bson.M{"_id": CookieData.Id}, ChangeStruct)
			if err != nil {
				log.Printf("Err insert", err)
			}
			collection.InsertOne(ctx, ChangeStruct)
			c.JSON(200, gin.H{
				"Code": "Your Request Successfully Handeled",
			})
		}else if checkPointTwo == false {
			c.JSON(404, gin.H{
				"Code": "User not Found1.2",
			})
		}
	} else if CookieData.Permissions == "client" {
		var ChangeStruct structures.Signup
		var DecodedSigninStruct structures.Signup
		json.Unmarshal([]byte(jsonFM), &ChangeStruct)
		checkPointOne, checkPointTwo := velidation.TestTheStruct(c, "email:phone", string(valueStruct), "FieldsCheck:true,DBCheck:true", "client", CookieData.Id)
		fmt.Printf("checkPointOne: %v\n", checkPointOne)
		fmt.Printf("checkPointTwo: %v\n", checkPointTwo)
		collection.FindOne(ctx, bson.M{"name": ChangeStruct.Name, "surname": ChangeStruct.Surname, "permissions": CookieData.Permissions}).Decode(&DecodedSigninStruct)
		if checkPointOne != false && checkPointTwo == false {
			log.Printf("ds1%v\n", ChangeStruct)
			DecodedSigninStruct.Userid = CookieData.Id
			DecodedSigninStruct.Permissions = CookieData.Permissions
			DecodedSigninStruct.Email = ChangeStruct.Email
			DecodedSigninStruct.Phone = ChangeStruct.Phone
			if ChangeStruct.Password != "" {
				hashedPass, err := bycrypt.HashPassword(ChangeStruct.Password)
				if err != nil {
					log.Printf("Err Hash%v", err)
				}
				DecodedSigninStruct.Password = hashedPass
			}
			if errIMG != nil {
				log.Printf("123%v\n", ChangeStruct)
				// ChangeStruct.ImgUrl = DecodedSigninStruct.ImgUrl
			} else {
				log.Printf("456%v\n", ChangeStruct)
				imgid := handlefile.Handlefile(c, "./static/upload")
				DecodedSigninStruct.ImgUrl = imgid
			}
			fmt.Printf("DecodedSigninStruct: %v\n", DecodedSigninStruct)
			_, err = collection.ReplaceOne(ctx, bson.M{"_id": CookieData.Id}, DecodedSigninStruct)
			if err != nil {
				log.Printf("Err insert", err)
			}
			// log.Printf("cur",cur)
			c.JSON(200, gin.H{
				"Code": "Your Request Successfully Handeled",
			})
		}else if checkPointTwo == false {
			c.JSON(404, gin.H{
				"Code": "User not Found1.2",
			})
		}
	}
}


func Cors(c *gin.Context) {
	fmt.Printf("Ip: %v\n", evtvariables.IpUrl)
	c.Writer.Header().Set("Access-Control-Allow-Origin", evtvariables.IpUrl)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}

	c.Next()
}
