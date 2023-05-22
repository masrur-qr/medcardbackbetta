package ehrcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medcard-new/begening/controllers/handlefile"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/controllers/velidation"
	"net/http"
	"strconv"
	"strings"

	"medcard-new/begening/evtvariables"
	"medcard-new/begening/structures"
	"os"
	"time"

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
	clientOptions := options.Client().ApplyURI(evtvariables.DBUrl)
	// clientOptions := options.Client().ApplyURI("mongodb://mas:mas@34.148.119.65:27017")

	clientG, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: ", err)
	}
	ctxG, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	// collection := client.Database("MedCard").Collection("users")
	ctx = ctxG
	client = clientG
}
func DoctorClientForView(c *gin.Context) {
	var (
		ViewStruct       structures.Views
		ViewStructDecode structures.Views
	)
	c.ShouldBindJSON(&ViewStruct)
	CookieData := jwtgen.Velidation(c)
	fmt.Printf("ViewStruct: %v\n", ViewStruct)

	stringJSON, err := json.Marshal(ViewStruct)
	if err != nil {
		log.Printf("Marshel err %v\n", err)
	}

	if CookieData.Permissions == "client" {
		isPassedFields, _ := velidation.TestTheStruct(c, "clientFLSname:doctorid:sickness:phone", string(stringJSON), "FieldsCheck:true,DBCheck:false", "", "")

		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		Authenticationservice()
		collection := client.Database("MedCard").Collection("views")
		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		err := collection.FindOne(ctx, bson.M{"clientid": CookieData.Id, "doctorid": ViewStruct.DoctorId}).Decode(&ViewStructDecode)
		if err != nil {
			log.Printf("Find ERR views%v\n", err)
		}
		if isPassedFields == true && ViewStructDecode.Sickness == "" {
			premetivid := primitive.NewObjectID().Hex()
			ViewStruct.Id = premetivid
			ViewStruct.ClientId = CookieData.Id
			ViewStruct.Date = ""
			_, err := collection.InsertOne(ctx, ViewStruct)
			if err != nil {
				log.Printf("Insert ERR views%v\n", err)
			}
		} else {
			c.JSON(400, gin.H{
				"Code": "You Just pusted such request",
			})
		}
	} else if CookieData.Permissions == "doctor" {
		isPassedFields, _ := velidation.TestTheStruct(c, "clientid:doctorid:date", string(stringJSON), "FieldsCheck:true,DBCheck:false", "", "")
		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		Authenticationservice()
		collection := client.Database("MedCard").Collection("views")
		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		err := collection.FindOne(ctx, bson.M{"clientid": ViewStruct.ClientId, "doctorid": ViewStruct.DoctorId}).Decode(&ViewStructDecode)
		if err != nil {
			log.Printf("Find ERR views%v\n", err)
		}
		// log.Printf("Insert ERR views%v\n", isPassedFields)

		if ViewStructDecode.Sickness != "" && isPassedFields == true {
			collection.DeleteOne(ctx, bson.M{"clientid": ViewStruct.ClientId, "doctorid": ViewStruct.DoctorId})

			// ?Create time zone from date and time
			splitedDate := strings.Split(ViewStruct.Date, "-")
			splitedTime := strings.Split(strings.Split(ViewStruct.Date, ";")[1], ":")
			Year, _ := strconv.Atoi(splitedDate[0])
			Month, _ := strconv.Atoi(splitedDate[1])
			Date, _ := strconv.Atoi(strings.Split(splitedDate[2], ";")[0])
			Hour, _ := strconv.Atoi(splitedTime[0])
			Minute, _ := strconv.Atoi(splitedTime[1])
			dateZoneFormat := fmt.Sprintf("%v-0%v-%vT%v:%v:%v.%vZ",Year,Month,Date,Hour,Minute,"00","050")
			ViewStructDecode.Date = dateZoneFormat
			// ? insert into db
			_, err = collection.InsertOne(ctx, ViewStructDecode)
			if err != nil {
				log.Printf("Insert || delete Error%v\n", err)
				return
			}
		} else {
			c.JSON(400, gin.H{
				"Code": "Cannot Find the user",
			})
		}
	}
}
func AddFilesToEhr(c *gin.Context) {
	var (
		FilesStruct structures.File
	)
	stringJSON := c.Request.FormValue("json")
	files, handler, errIMG := c.Request.FormFile("img")
	// """""""""""""""""""""""check The file on existense"""""""""""""""""""""""
	if errIMG != nil {
		c.JSON(409, gin.H{
			"sttus": "NOIMGFILEEXIST",
		})
	}

	files.Seek(23, 23)
	log.Printf("File Name %s%v\n", handler.Filename, stringJSON)
	json.Unmarshal([]byte(stringJSON), &FilesStruct)
	jsStr, err := json.Marshal(FilesStruct)
	if err != nil {
		log.Printf("%v", err)
	}
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("ehrfiles")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	isPassedFields, _ := velidation.TestTheStruct(c, "clientFLSname:clientid:doctorid:description:doctorFLSname:title", string(jsStr), "FieldsCheck:true,DBCheck:false", "", "")
	if isPassedFields == true {
		premetivid := primitive.NewObjectID().Hex()
		FilesStruct.Id = premetivid
		FilesStruct.ImgUrl = handlefile.Handlefile(c, "./static/upload")
		collection.InsertOne(ctx, FilesStruct)
		c.JSON(200, gin.H{
			"Code": "Request Seccessfully Handleed",
		})
	}
}
func ExpiredLinks(c *gin.Context) {
	//! http://127.0.0.1:5500/link?client=6468f42e1b2b6c995ac8dfc8&id=345464489.jpg&type=client
	typeOfFile := c.Request.URL.Query().Get("type")
	if typeOfFile == "doctor" {
		staticFiles(c, "./static/doctors/upload-", "")
	} else if typeOfFile == "client" {
		cookieData := jwtgen.Velidation(c)
		if cookieData.Permissions == "client" || cookieData.Permissions == "doctor" {
			staticFiles(c, "./static/upload/upload-", cookieData.Id)
		}
	}
}
func staticFiles(c *gin.Context, path string, id string) {
	fmt.Println(time.Now())
	imgId := c.Request.URL.Query().Get("id")
	clientId := c.Request.URL.Query().Get("client")
	// parsing time "2023-5-22T18:20:00Z05:00" as "2006-01-02T15:04:05Z07:00": cannot parse "5-22T18:20:00Z05:00" as "01"
	var (
		ehrfiles  structures.File
		viewsList structures.Views
	)
	if id != "" {
		Authenticationservice()
		conn := client.Database("MedCard").Collection("ehrfiles")
		conn.FindOne(ctx, bson.M{"clientid": id, "imgurl": "-" + imgId}).Decode(&ehrfiles)
		if ehrfiles.ImgUrl == "-"+imgId {
			http.ServeFile(c.Writer, c.Request, path+imgId)
		} else {
			connView := client.Database("MedCard").Collection("views")
			connView.FindOne(ctx, bson.M{"clientid": clientId, "doctorid": id}).Decode(&viewsList)
			if viewsList.ClientId == id || viewsList.DoctorId == id {
				fmt.Println("2.1")
				Curenttime := time.Now().UTC()
				NewTimeZone := time.FixedZone("Tajikistan", 5*3600)
				tajikistanTimeZone := Curenttime.In(NewTimeZone)
				utcTime, err := time.Parse(time.RFC3339, viewsList.Date)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("tajikistanTimeZone: %v\n", tajikistanTimeZone)
				fmt.Printf("expireTime: %v\n", utcTime)
				if tajikistanTimeZone.After(utcTime) {
					c.JSON(404, gin.H{
						"Code": "Your session to this file expired",
					})
				} else {
					http.ServeFile(c.Writer, c.Request, path+imgId)
				}
			} else {
				c.JSON(404, gin.H{
					"Code": "Your session to this file expired || You have no access to it",
				})
			}
		}

	} else {
		fmt.Println("1.2")
		http.ServeFile(c.Writer, c.Request, path+imgId)
	}
}
