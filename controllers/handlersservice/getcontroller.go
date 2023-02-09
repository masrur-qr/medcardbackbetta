package controllers

import (
	"log"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/structures"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuestions(c *gin.Context){
	var(
		QuestionsDb structures.Questions
	)
	Authenticationservice()
	collection := client.Database("MedCard").Collection("questions")
	cur, err := collection.Find(ctx, bson.M{})

	if err != nil{
		log.Printf("Err find Question%v\n",err)
	}
	defer cur.Close(ctx)

	var QuestionsDbArr []structures.Questions
	for cur.Next(ctx){
		cur.Decode(&QuestionsDb)
		QuestionsDbArr = append(QuestionsDbArr, QuestionsDb)
		log.Panicln(QuestionsDbArr)
	}

	c.JSON(200,gin.H{
		"Code":"Request Handeled",
		"Json": QuestionsDbArr,
	})
}
func GetDoctors(c *gin.Context){
	var(
		DoctorDb structures.SignupDoctor
	)
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	cur, err := collection.Find(ctx, bson.M{"permissions":"doctor"})

	if err != nil{
		log.Printf("Err find Question%v\n",err)
	}
	defer cur.Close(ctx)

	var DoctorDbArr []structures.SignupDoctor
	for cur.Next(ctx){
		cur.Decode(&DoctorDb)
		DoctorDbArr = append(DoctorDbArr, DoctorDb)
		log.Panicln(DoctorDbArr)
	}

	c.JSON(200,gin.H{
		"Code":"Request Handeled",
		"Json": DoctorDbArr,
	})
}
func Statistics(c *gin.Context){
	var(
		Statistics structures.GlobeStruct
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "admin"{
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		cur, err := collection.Find(ctx, bson.M{"sta":"doctor"})

		if err != nil{
			log.Printf("Err find Question%v\n",err)
		}
		defer cur.Close(ctx)

		var StatisticsArr []structures.GlobeStruct
		for cur.Next(ctx){
			cur.Decode(&Statistics)
			StatisticsArr = append(StatisticsArr, Statistics)
			log.Panicln(StatisticsArr)
		}

		// ======================================= All Users =====================================
		c.JSON(200,gin.H{
			"Code":"Request Handeled",
			"Users": len(StatisticsArr)-1,
		})
	}else{
		c.JSON(400,gin.H{
			"Code":"No Permissions",
		})
	}
}
func GetClients(c *gin.Context){
	var(
		ClientsDB structures.Signup
	)
    // ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "admin"{
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		cur, err := collection.Find(ctx, bson.M{"permissions":"client"})

		if err != nil{
			log.Printf("Err find Question%v\n",err)
		}
		defer cur.Close(ctx)

		var ClientsDBArr []structures.Signup
		for cur.Next(ctx){
			cur.Decode(&ClientsDB)
			ClientsDBArr = append(ClientsDBArr, ClientsDB)
			log.Panicln(ClientsDBArr)
		}

		c.JSON(200,gin.H{
			"Code":"Request Handeled",
			"Json": ClientsDBArr,
		})
	}else{
		c.JSON(400,gin.H{
			"Code":"No Permissions",
		})
	}
}
func GetViews(c *gin.Context){
	var(
		EHRFileDB structures.File
		cur *mongo.Cursor
		err error
	)
    // ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "client" || CookieData.Permissions == "doctor"{
		Authenticationservice()
		collection := client.Database("MedCard").Collection("views")
		if CookieData.Permissions == "client"{
			cur, err = collection.Find(ctx, bson.M{"clientid": CookieData.Id})
		}else if CookieData.Permissions == "doctor"{
			cur, err = collection.Find(ctx, bson.M{"doctorid": CookieData.Id})
		}

		if err != nil{
			log.Printf("Err find Question%v\n",err)
		}
		defer cur.Close(ctx)

		var EHRFileDBArr []structures.File
		for cur.Next(ctx){
			cur.Decode(&EHRFileDB)
			EHRFileDBArr = append(EHRFileDBArr, EHRFileDB)
			log.Panicln(EHRFileDBArr)
		}

		c.JSON(200,gin.H{
			"Code":"Request Handeled",
			"Json": EHRFileDBArr,
		})
	}else{
		c.JSON(400,gin.H{
			"Code":"No Permissions",
		})
	}
}
func GetClient(c *gin.Context){
	var(
		ClientsDB structures.Signup
		Files structures.File
	)
    // ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	log.Println(CookieData.Permissions)

	if CookieData.Permissions == "client" || CookieData.Permissions == "doctor"{
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		err := collection.FindOne(ctx, bson.M{"_id": c.Request.URL.RawQuery}).Decode(&ClientsDB)

		collectionView := client.Database("MedCard").Collection("ehrfiles")
		cur , errTwo := collectionView.Find(ctx, bson.M{"clientid": c.Request.URL.RawQuery})

		defer cur.Close(ctx)

		var ViewsArr []structures.File
		for cur.Next(ctx){
			cur.Decode(&Files)
			ViewsArr = append(ViewsArr, Files)
		}

		if err != nil || errTwo != nil{
			log.Printf("Err find user %v\n",err)
		}
		if ClientsDB.Name != ""{
			log.Println(c.Request.URL.RawQuery)
	
			c.JSON(200,gin.H{
				"Code" : "Request Handeled",
				"Json" : ClientsDB,
				"Files" : ViewsArr,
			})
		}else{
			c.JSON(400,gin.H{
				"Code":"User NotFound",
			})
		}
	}else{
		c.JSON(400,gin.H{
			"Code":"No Permissions",
		})
	}
}
