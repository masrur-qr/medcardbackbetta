package ehrcontroller

import (
	"context"
	"encoding/json"
	"log"
	"medcard-new/begening/controllers/handlefile"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/controllers/velidation"
	// "medcard-new/begening/evtvariables"
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
	ctx context.Context
	client *mongo.Client
)
var DB_Url string = os.Getenv("DBURL")

func Authenticationservice(){
	// clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// clientOptions := options.Client().ApplyURI(evtvariables.DBUrl)
	clientOptions := options.Client().ApplyURI("mongodb://mas:mas@34.148.119.65:27017")

	clientG, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: ", err)
	}
	ctxG, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	// collection := client.Database("MedCard").Collection("users")
	ctx = ctxG
	client = clientG
}
func DoctorClientForView(c *gin.Context){
	var(
		ViewStruct structures.Views
		ViewStructDecode structures.Views
	)
	c.ShouldBindJSON(&ViewStruct)
	CookieData := jwtgen.Velidation(c)

	stringJSON , err:= json.Marshal(ViewStruct)
	if err != nil{
		log.Printf("Marshel err %v\n",err)
	}

	if CookieData.Permissions == "client"{
		isPassedFields , _ := velidation.TestTheStruct(c,"clientFLSname:clientid:doctorid:sickness:phone",string(stringJSON),"FieldsCheck:true,DBCheck:false","","")

		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		Authenticationservice()
		collection := client.Database("MedCard").Collection("views")
		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		err := collection.FindOne(ctx,bson.M{"clientid":ViewStruct.ClientId,"doctorid":ViewStruct.DoctorId}).Decode(&ViewStructDecode)
		if err != nil{
			log.Printf("Find ERR views%v\n", err)
		}
		if isPassedFields == true && ViewStructDecode.Sickness == ""{
			premetivid := primitive.NewObjectID().Hex()
			ViewStruct.Id = premetivid
			ViewStruct.Date = ""
			_ , err := collection.InsertOne(ctx,ViewStruct)
			if err != nil{
				log.Printf("Insert ERR views%v\n", err)
			}
		}else{
			c.JSON(400, gin.H{
				"Code":"You Just pusted such request",
			})
		}
	}else if CookieData.Permissions == "doctor"{
		isPassedFields , _ := velidation.TestTheStruct(c,"clientid:doctorid:date",string(stringJSON),"FieldsCheck:true,DBCheck:false","","")
		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		Authenticationservice()
		collection := client.Database("MedCard").Collection("views")
		// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
		err := collection.FindOne(ctx,bson.M{"clientid":ViewStruct.ClientId,"doctorid":ViewStruct.DoctorId}).Decode(&ViewStructDecode)
		if err != nil{
			log.Printf("Find ERR views%v\n", err)
		}
		log.Printf("Insert ERR views%v\n", isPassedFields)
	
		if ViewStructDecode.Sickness != "" && isPassedFields == true{
			collection.DeleteOne(ctx,bson.M{"clientid":ViewStruct.ClientId,"doctorid":ViewStruct.DoctorId})

			ViewStructDecode.Date = ViewStruct.Date
			_ ,err := collection.InsertOne(ctx,ViewStructDecode)
			if err != nil{
				log.Printf("Insert || delete Error%v\n", err)
				return
			}
		}else{
			c.JSON(400, gin.H{
				"Code":"Cannot Find the user",
			})
		}
	}
}
func AddFilesToEhr(c *gin.Context){
	var(
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
	log.Printf("File Name %s%v\n", handler.Filename,stringJSON)
	json.Unmarshal([]byte(stringJSON),&FilesStruct)
	jsStr , err := json.Marshal(FilesStruct)
	if err != nil{
		log.Printf("%v",err)
	}
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	Authenticationservice()
	collection := client.Database("MedCard").Collection("ehrfiles")
	// """"""""""""""""""""""""""""""""""DB CONNECTION""""""""""""""""""""""""""""""""""""""""""""""""""""
	isPassedFields , _ := velidation.TestTheStruct(c,"clientFLSname:clientid:doctorid:description:doctorFLSname:title",string(jsStr),"FieldsCheck:true,DBCheck:false","","")
	if isPassedFields == true{
		premetivid := primitive.NewObjectID().Hex()
		FilesStruct.Id = premetivid
		FilesStruct.ImgUrl = handlefile.Handlefile(c,"./static/uploadfille")
		collection.InsertOne(ctx,FilesStruct)
		c.JSON(200,gin.H{
			"Code":"Request Seccessfully Handleed",
		})
	}
}