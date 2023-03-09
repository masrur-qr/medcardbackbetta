package controllers

import (
	"log"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/structures"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuestions(c *gin.Context) {
	var (
		QuestionsDb structures.Questions
	)
	Authenticationservice()
	collection := client.Database("MedCard").Collection("questions")
	cur, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Printf("Err find Question%v\n", err)
	}
	defer cur.Close(ctx)

	var QuestionsDbArr []structures.Questions
	for cur.Next(ctx) {
		cur.Decode(&QuestionsDb)
		QuestionsDbArr = append(QuestionsDbArr, QuestionsDb)
		log.Println(QuestionsDbArr)
	}

	c.JSON(200, gin.H{
		"Code": "Request Handeled",
		"Json": QuestionsDbArr,
	})
}
func GetDoctors(c *gin.Context) {
	var (
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
		log.Println(DoctorDbArr)
	}

	c.JSON(200, gin.H{
		"Code": "Request Handeled",
		"Json": DoctorDbArr,
	})
}
func Statistics(c *gin.Context) {
	var (
		Statistics structures.GlobeStruct
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "admin" {
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		// ==================== List all users =========================
		cur, err := collection.Find(ctx, bson.M{})

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer cur.Close(ctx)

		var StatisticsArr []structures.GlobeStruct
		for cur.Next(ctx) {
			cur.Decode(&Statistics)
			StatisticsArr = append(StatisticsArr, Statistics)
			log.Println(StatisticsArr)
		}
		// ==================== List all doctors =========================
		curDoc, err := collection.Find(ctx, bson.M{"permissions":"doctor"})

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer curDoc.Close(ctx)

		var StatisticsArrDoc []structures.GlobeStruct
		for curDoc.Next(ctx) {
			curDoc.Decode(&Statistics)
			StatisticsArrDoc = append(StatisticsArrDoc, Statistics)
			log.Println(StatisticsArrDoc)
		}
		// ==================== List all doctors =========================
		curCl, err := collection.Find(ctx, bson.M{"permissions":"client"})

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer curCl.Close(ctx)

		var StatisticsArrCl []structures.GlobeStruct
		for curCl.Next(ctx) {
			curCl.Decode(&Statistics)
			StatisticsArrCl = append(StatisticsArrCl, Statistics)
			log.Println(StatisticsArrCl)
		}

		// ======================================= All Users =====================================
		c.JSON(200, gin.H{
			"Code":  "Request Handeled",
			"Users": len(StatisticsArr) - 1,
			"Doctors": len(StatisticsArrDoc),
			"Clients": len(StatisticsArrCl),
		})
	} else {
		c.JSON(400, gin.H{
			"Code": "No Permissions",
		})
	}
}
func GetClients(c *gin.Context) {
	var (
		ClientsDB structures.Signup
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "admin" {
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		cur, err := collection.Find(ctx, bson.M{"permissions": "client"})

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer cur.Close(ctx)

		var ClientsDBArr []structures.Signup
		for cur.Next(ctx) {
			cur.Decode(&ClientsDB)
			ClientsDBArr = append(ClientsDBArr, ClientsDB)
			log.Println(ClientsDBArr)
		}

		c.JSON(200, gin.H{
			"Code": "Request Handeled",
			"Json": ClientsDBArr,
		})
	} else {
		c.JSON(400, gin.H{
			"Code": "No Permissions",
		})
	}
}
func GetViews(c *gin.Context) {
	var (
		EHRFileDB structures.Views
		cur       *mongo.Cursor
		err       error
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "client" || CookieData.Permissions == "doctor" {
		Authenticationservice()
		collection := client.Database("MedCard").Collection("views")
		if CookieData.Permissions == "client" {
			cur, err = collection.Find(ctx, bson.M{"clientid": CookieData.Id, "doctorid": c.Request.URL.RawQuery})
		} else if CookieData.Permissions == "doctor" {
			cur, err = collection.Find(ctx, bson.M{"doctorid": CookieData.Id, "clientid": c.Request.URL.RawQuery})
		}

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer cur.Close(ctx)

		var EHRFileDBArr []structures.Views
		for cur.Next(ctx) {
			cur.Decode(&EHRFileDB)
			EHRFileDBArr = append(EHRFileDBArr, EHRFileDB)
			log.Println(EHRFileDBArr)
		}

		c.JSON(200, gin.H{
			"Code": "Request Handeled",
			"Json": EHRFileDBArr,
		})
	} else {
		c.JSON(400, gin.H{
			"Code": "No Permissions",
		})
	}
}
func GetClient(c *gin.Context) {
	var (
		ClientsDB structures.Signup
		DoctorDB structures.SignupDoctor
		Files     structures.File
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	log.Println(CookieData.Permissions)

	if CookieData.Permissions == "client"{
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		err := collection.FindOne(ctx, bson.M{"_id": c.Request.URL.RawQuery}).Decode(&ClientsDB)

		collectionView := client.Database("MedCard").Collection("ehrfiles")
		cur, errTwo := collectionView.Find(ctx, bson.M{"clientid": c.Request.URL.RawQuery})

		defer cur.Close(ctx)

		var ViewsArr []structures.File
		for cur.Next(ctx) {
			cur.Decode(&Files)
			ViewsArr = append(ViewsArr, Files)
		}

		if err != nil || errTwo != nil {
			log.Printf("Err find user %v\n", err)
		}
		if ClientsDB.Name != "" {
			log.Println(c.Request.URL.RawQuery)
			ClientsDB.Password = "null"
			c.JSON(200, gin.H{
				"Code":  "Request Handeled",
				"Json":  ClientsDB,
				"Files": ViewsArr,
			})
		} else {
			c.JSON(400, gin.H{
				"Code": "User NotFound",
			})
		}
	}else {
		// ? ========================================== get Client data ============================
		Authenticationservice()
		collectionCli := client.Database("MedCard").Collection("users")
		err := collectionCli.FindOne(ctx, bson.M{"_id": c.Request.URL.RawQuery}).Decode(&ClientsDB)

		collectionView := client.Database("MedCard").Collection("ehrfiles")
		cur, errTwo := collectionView.Find(ctx, bson.M{"clientid": c.Request.URL.RawQuery})

		defer cur.Close(ctx)

		var ViewsArr []structures.File
		for cur.Next(ctx) {
			cur.Decode(&Files)
			ViewsArr = append(ViewsArr, Files)
		}

		if err != nil || errTwo != nil {
			log.Printf("Err find user %v\n", err)
		}
		// ? ================================== Get doctor data ==========================
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		err = collection.FindOne(ctx, bson.M{"_id": CookieData.Id}).Decode(&DoctorDB)

		if err != nil{
			log.Printf("Err find user %v\n", err)
		}
		if DoctorDB.Name != "" {
			log.Println(c.Request.URL.RawQuery)
			DoctorDB.Password = "null"
			c.JSON(200, gin.H{
				"Code":  "Request Handeled",
				"UserJson":  ClientsDB,
				"Files": ViewsArr,
				"Json":  DoctorDB,
			})
		} else {
			c.JSON(400, gin.H{
				"Code": "User NotFound",
			})
		}
	}
}
