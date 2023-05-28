package controllers

import (
	"context"
	"fmt"
	"log"
	"medcard-new/begening/controllers/jwtgen"
	"medcard-new/begening/structures"
	"time"

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
		// log.Println(QuestionsDbArr)
	}
	if len(QuestionsDbArr) == 0 {
		c.JSON(200, gin.H{
			"Code": "Request Handeled",
			"Json": []string{},
		})
	} else {
		c.JSON(200, gin.H{
			"Code": "Request Handeled",
			"Json": QuestionsDbArr,
		})
	}
}
func GetDoctors(c *gin.Context) {
	var (
		DoctorDb structures.SignupDoctor
	)
	Authenticationservice()
	collection := client.Database("MedCard").Collection("users")
	cur, err := collection.Find(ctx, bson.M{"permissions": "doctor"})

	if err != nil {
		log.Printf("Err find Question%v\n", err)
	}
	defer cur.Close(ctx)

	var DoctorDbArr []structures.SignupDoctor
	for cur.Next(ctx) {
		cur.Decode(&DoctorDb)
		DoctorDbArr = append(DoctorDbArr, DoctorDb)
		// log.Println(DoctorDbArr)
	}

	if len(DoctorDbArr) == 0 {
		c.JSON(200, gin.H{
			"Code": "Request Handeled",
			"Json": []string{},
		})
	} else {
		c.JSON(200, gin.H{
			"Code": "Request Handeled",
			"Json": DoctorDbArr,
		})
	}
}
func Statistics(c *gin.Context) {
	var (
		Statistics      structures.GlobeStruct
		StatisticsUsers structures.Signup
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
			// log.Println(StatisticsArr)
		}
		// ==================== List all doctors =========================
		curDoc, err := collection.Find(ctx, bson.M{"permissions": "doctor"})

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer curDoc.Close(ctx)

		var StatisticsArrDoc []structures.GlobeStruct
		for curDoc.Next(ctx) {
			curDoc.Decode(&Statistics)
			StatisticsArrDoc = append(StatisticsArrDoc, Statistics)
			// log.Println(StatisticsArrDoc)
		}
		// ==================== List all doctors =========================
		curCl, err := collection.Find(ctx, bson.M{"permissions": "client"})

		if err != nil {
			log.Printf("Err find Question%v\n", err)
		}
		defer curCl.Close(ctx)

		var StatisticsArrCl []structures.GlobeStruct
		for curCl.Next(ctx) {
			curCl.Decode(&Statistics)
			StatisticsArrCl = append(StatisticsArrCl, Statistics)
			// log.Println(StatisticsArrCl)
		}
		// ======================================= Filter Users By Blood =====================================
		// ? Blood 1
		var StatisticUserArrOne []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "1"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrOne = append(StatisticUserArrOne, StatisticsUsers)
		}
		// ? Blood 2
		var StatisticUserArrTwo []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "2"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrTwo = append(StatisticUserArrTwo, StatisticsUsers)
		}
		// ? Blood 3
		var StatisticUserArrThree []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "3"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrThree = append(StatisticUserArrThree, StatisticsUsers)
		}
		// ? Blood 4
		var StatisticUserArrFour []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "4"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrFour = append(StatisticUserArrFour, StatisticsUsers)
		}
		// ======================================= Filter Users By Gender =====================================
		var StatisticUserArrGenderMale []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"gender": "male"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrGenderMale = append(StatisticUserArrGenderMale, StatisticsUsers)
		}
		var StatisticUserArrGenderFemale []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"gender": "female"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrGenderFemale = append(StatisticUserArrGenderFemale, StatisticsUsers)
		}
		// ======================================= Disabilaties =====================================
		// ? Disabliaties 1
		var StatisticUserArrAbilatiesOne []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "1"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrAbilatiesOne = append(StatisticUserArrAbilatiesOne, StatisticsUsers)
		}
		// ? Disabliaties 2
		var StatisticUserArrAbilatiesTwo []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "2"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrAbilatiesTwo = append(StatisticUserArrAbilatiesTwo, StatisticsUsers)
		}
		// ? Disabliaties 3
		var StatisticUserArrAbilatiesThree []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "3"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrAbilatiesThree = append(StatisticUserArrAbilatiesThree, StatisticsUsers)
		}
		// ? Disabliaties 4
		var StatisticUserArrAbilatiesFour []structures.Signup
		cur, _ = collection.Find(ctx, bson.M{"blood": "4"})
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			cur.Decode(&StatisticsUsers)
			StatisticUserArrAbilatiesFour = append(StatisticUserArrAbilatiesFour, StatisticsUsers)
		}
		// ======================================= All Users =====================================
		c.JSON(200, gin.H{
			"Code":              "Request Handeled",
			"Users":             len(StatisticsArr) - 1,
			"Doctors":           len(StatisticsArrDoc),
			"Clients":           len(StatisticsArrCl),
			"Blood-1":           len(StatisticUserArrOne),
			"Blood-2":           len(StatisticUserArrTwo),
			"Blood-3":           len(StatisticUserArrThree),
			"Blood-4":           len(StatisticUserArrFour),
			"Gender-Male":       len(StatisticUserArrGenderMale),
			"Gender-Female":     len(StatisticUserArrGenderFemale),
			"Disabliaties-1":    len(StatisticUserArrAbilatiesOne),
			"Disabliaties-2":    len(StatisticUserArrAbilatiesTwo),
			"Disabliaties-3":    len(StatisticUserArrAbilatiesThree),
			"Disabliaties-None": len(StatisticUserArrAbilatiesFour),
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
			// log.Println(ClientsDBArr)
		}

		if len(ClientsDBArr) == 0 {
			c.JSON(200, gin.H{
				"Code": "Request Handeled",
				"Json": []string{},
			})
		} else {
			c.JSON(200, gin.H{
				"Code": "Request Handeled",
				"Json": ClientsDBArr,
			})
		}
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
			cur, err = collection.Find(ctx, bson.M{"clientid": CookieData.Id})
		} else if CookieData.Permissions == "doctor" {
			cur, err = collection.Find(ctx, bson.M{"doctorid": CookieData.Id})
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
		if len(EHRFileDBArr) == 0 {
			c.JSON(200, gin.H{
				"Code": "Request Handeled",
				"Json": []string{},
			})
		} else {
			c.JSON(200, gin.H{
				"Code": "Request Handeled",
				"Json": EHRFileDBArr,
			})
		}

	} else {
		c.JSON(400, gin.H{
			"Code": "No Permissions",
		})
	}
}
func ListViewsAdmin(c *gin.Context) {
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "admin" {
		var (
			DecodeViews     structures.Views
			DecodeViewsArch structures.Views
			DecodeViewsArr  []structures.Views
		)
		Authenticationservice()
		connArch := client.Database("MedCard").Collection("viewsarchive")
		client := client.Database("MedCard").Collection("views")
		// ! Get Unarchived views
		cur, err := client.Find(ctx, bson.M{})
		if err != nil {
			fmt.Printf("Error No user found%v", err)
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			cur.Decode(&DecodeViews)
			DecodeViewsArr = append(DecodeViewsArr, DecodeViews)
		}
		// !  Get Archived views
		curArch, err := connArch.Find(ctx, bson.M{})
		if err != nil {
			fmt.Printf("Error Arch %v", err)
		}
		defer curArch.Close(ctx)
		for curArch.Next(ctx) {
			curArch.Decode(&DecodeViewsArch)
			DecodeViewsArr = append(DecodeViewsArr, DecodeViewsArch)
		}

		if len(DecodeViewsArr) == 0 {
			c.JSON(200, gin.H{
				"Views": []string{},
			})
		} else {
			c.JSON(200, gin.H{
				"Views": DecodeViewsArr,
			})
		}
	} else {
		c.JSON(400, gin.H{
			"Code": "No Permissions",
		})
	}
}
func GetClient(c *gin.Context) {
	var (
		ClientsDB structures.Signup
		DoctorDB  structures.SignupDoctor
		Files     structures.File
		AdminDecode structures.Admin
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	log.Println(CookieData.Permissions)

	if CookieData.Permissions == "client" {
		fmt.Println("test client")

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
			if len(ViewsArr) == 0 {
				c.JSON(200, gin.H{
					"Code":  "Request Handeled",
					"Json":  ClientsDB,
					"Files": []string{},
				})
			} else {
				c.JSON(200, gin.H{
					"Code":  "Request Handeled",
					"Json":  ClientsDB,
					"Files": ViewsArr,
				})
			}
		} else {
			c.JSON(400, gin.H{
				"Code": "User NotFound",
			})
		}
	}else if CookieData.Permissions == "admin" {
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		collection.FindOne(ctx, bson.M{"_id": c.Request.URL.RawQuery }).Decode(&ClientsDB)
		collection.FindOne(ctx, bson.M{"_id": CookieData.Id}).Decode(&AdminDecode)
		fmt.Println("test admin")
		fmt.Printf("AdminDecode: %v\n", AdminDecode)
		log.Println(c.Request.URL.RawQuery)

		if AdminDecode.Permissions == "admin" && ClientsDB.Lastname != ""{
			fmt.Printf("ClientsDB: %v\n", ClientsDB)
			if ClientsDB.Name != "" {
				log.Println(c.Request.URL.RawQuery)
				ClientsDB.Password = "null"
				c.JSON(200, gin.H{
					"Code": "Request Handeled",
					"Json": ClientsDB,
				})
			} else {
				c.JSON(400, gin.H{
					"Code": "User NotFound",
				})
			}
		}else if AdminDecode.Permissions == "admin" && c.Request.URL.RawQuery == "" {
			fmt.Println("dede")
			if AdminDecode.Name != "" {
				log.Println(c.Request.URL.RawQuery)
				AdminDecode.Password = "null"
				c.JSON(200, gin.H{
					"Code": "Request Handeled",
					"Json": AdminDecode,
				})
			} else {
				c.JSON(400, gin.H{
					"Code": "User NotFound",
				})
			}
		}
		
	} else if CookieData.Permissions == "doctor" {
		fmt.Println(c.Request.URL.RawQuery)

		var (
			DecodeViews structures.Views
		)
		collectionViews := client.Database("MedCard").Collection("views")
		ctxForAccess , _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := collectionViews.FindOne(ctxForAccess, bson.M{"doctorid": CookieData.Id, "clientid": c.Request.URL.RawQuery}).Decode(&DecodeViews)
		if err != nil {
			fmt.Printf("err doc %v", err)
		}
		collectionCli := client.Database("MedCard").Collection("users")
		err = collectionCli.FindOne(ctxForAccess, bson.M{"_id": c.Request.URL.RawQuery}).Decode(&ClientsDB)
		// ? ================================== Get doctor data ==========================
		err = collectionCli.FindOne(ctx, bson.M{"_id": CookieData.Id}).Decode(&DoctorDB)
		if err != nil {
			log.Printf("Err find user %v\n", err)
		}
		// ? ========== Validate if it this doctor has beed set data and has access to views ===========
		if DecodeViews.DoctorId == CookieData.Id && ClientsDB.Permissions != "doctor" && DecodeViews.Date != "" {
			// ? ===============get Client data ClientsDB taken priveously is also wil be sended  ===========
			Authenticationservice()

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

			
			if DoctorDB.Name != "" {
				log.Println(c.Request.URL.RawQuery)
				DoctorDB.Password = "null"
				c.JSON(200, gin.H{
					"Code":     "Request Handeled",
					"UserJson": ClientsDB,
					"Files":    ViewsArr,
					"Json":     DoctorDB,
				})
			} else {
				c.JSON(400, gin.H{
					"Code": "User NotFound",
				})
			}
		} else {
			DoctorDB.Password = "null"
			c.JSON(200, gin.H{
				"UserJson": DoctorDB ,
			})
		}
	}
}
