package velidation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medcard-new/begening/structures"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx context.Context
	client *mongo.Client
)
func Authenticationservice(){
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// clientOptions := options.Client().ApplyURI(os.Getenv("DB_URL"))
	clientG, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Mongo.connect() ERROR: ", err)
	}
	ctxG, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	// collection := client.Database("MedCard").Collection("users")
	ctx = ctxG
	client = clientG
}

func TestTheStruct(c *gin.Context,Required string,valueStruct string,options string,permission string,ID string) (bool,bool){
	var (
		countFields int = 0
		ControlCheck bool = false
		IsPassedFields bool = false
		IsPassedDB bool = false
		optionsArr = strings.Split(options, ",")
	)
	// """"""""""""""""""""""FIELDS CHECKER""""""""""""""""""""""""""
	if strings.Split(optionsArr[0], ":")[1] == "true"{
		splitReqField := strings.Split(strings.Join(strings.Split(strings.Join(strings.Split(string(valueStruct), "{"), ""), "}"), ""), ",")
		splitTheRequired := strings.Split(Required, ":")
		for i := 0; i < len(splitReqField); i++{
			if strings.Split(splitReqField[i], ":")[1] != `""`{
				for u := 0; u < len(splitTheRequired); u++{
					if strings.Join(strings.Split(strings.Split(splitReqField[i], ":")[0], `"`), "") == splitTheRequired[u]{
						countFields += 1
					}
				}
			}
		}
	fmt.Printf("strings.Split(optionsArr[1], \":\")[1]: %v\n", countFields)
		if countFields == len(splitTheRequired){
			IsPassedFields = true
		}else{
			c.JSON(400,gin.H{
				"Code":"Error there is empty field || fields",
			})
		}
	}
	// """"""""""""""""""""""FIELDS CHECKER""""""""""""""""""""""""""
	// """"""""""""""""""""""DB FOR EXISTENCE CHECKER""""""""""""""""""""""""""
	fmt.Printf("strings.Split(optionsArr[1], \":\")[1]: %v\n", strings.Split(optionsArr[1], ":")[1])
	if strings.Split(optionsArr[1], ":")[1] == "true"{
		var(
			SigninStruct structures.GlobeStruct
			DecodedSigninStruct structures.GlobeStruct
		)
		json.Unmarshal([]byte(valueStruct),&SigninStruct)
		log.Println(SigninStruct)
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		if ID != ""{
			
			if permission == "client"{
				err := collection.FindOne(ctx,bson.M{"_id":ID,"permissions":permission}).Decode(&DecodedSigninStruct)
				if err != nil{
					c.JSON(200, gin.H{
						"Code":"User not Found",
					})
				}
			}else{
				err := collection.FindOne(ctx,bson.M{"_id":ID,"permissions":permission}).Decode(&DecodedSigninStruct)
				if err != nil{
					c.JSON(200, gin.H{
						"Code":"User not Found",
					})
				}
			}
			log.Println("ss")
		}else{
			
			if permission == "client"{
				err := collection.FindOne(ctx,bson.M{"_id":ID,"permissions":permission}).Decode(&DecodedSigninStruct)
				if err != nil{
					c.JSON(200, gin.H{
						"Code":"User not Found",
					})
				}
			}else{
				err := collection.FindOne(ctx,bson.M{"name":SigninStruct.Name,"surname":SigninStruct.Surname,"permissions":permission}).Decode(&DecodedSigninStruct)
				if err != nil{
					c.JSON(200, gin.H{
						"Code":"User not Found",
					})
				}
			}
			
			log.Println("err")
		}
		if DecodedSigninStruct.Name == ""{
			collection.FindOne(ctx,bson.M{"phone":SigninStruct.Phone}).Decode(&DecodedSigninStruct)
			if DecodedSigninStruct.Name == ""{
				IsPassedDB = true
			}else{
				c.JSON(400,gin.H{
					"Code":"Error this phone numbers are already taken",
				})
			}
		}else{
			c.JSON(400,gin.H{
				"Code":"Error User already exist",
			})
		}
	}
	// """"""""""""""""""""""DB FOR EXISTENCE CHECKER""""""""""""""""""""""""""
	
// log.Printf("t1%v",IsPassedFields == true || strings.Split(optionsArr[0], ":")[1] == "false" )
// log.Printf("t1%v",!IsPassedFields || strings.Split(optionsArr[0], ":")[1] == "false" )
// log.Printf("t1%v",IsPassedDB == true || strings.Split(optionsArr[1], ":")[1] == "false")
	// if IsPassedFields == true || strings.Split(optionsArr[0], ":")[1] == "false" && IsPassedDB == true || strings.Split(optionsArr[1], ":")[1] == "false"{
	// 	ControlCheck = true
	// }

	log.Printf("check%v\n",ControlCheck)
	return IsPassedFields,IsPassedDB
}