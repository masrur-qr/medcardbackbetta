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
		if len(DoctorDb.History) == 0 {
			DoctorDb.History = []structures.History{}
		}
		DoctorDbArr = append(DoctorDbArr, DoctorDb)
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
		StatisticsUsers structures.Signup
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	if CookieData.Permissions == "admin" {
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")

		var usersAmount = []string{
			"doctor",
			"client",
			"allUsers",
		}

		var countUsers = make(map[string]int, 1)

		for i := 0; i < len(usersAmount); i++ {
			// ? Bloodes
			if usersAmount[i] == "allUsers" {
				countUsers[usersAmount[i]] = countUsers["client"] + countUsers["doctor"]
			} else {
				var StatisticUser []structures.Signup
				cur, _ := collection.Find(ctx, bson.M{"permissions": usersAmount[i]})
				defer cur.Close(ctx)
				for cur.Next(ctx) {
					cur.Decode(&StatisticsUsers)
					StatisticUser = append(StatisticUser, StatisticsUsers)
				}
				countUsers[usersAmount[i]] = len(StatisticUser)
			}
		}
		// ======================================= Filter Users By Blood =====================================

		var BloodTypes = []string{
			"I(0)+",
			"II(A)+",
			"III(B)+",
			"IV(AB)+",
			"I(0)-",
			"II(A)-",
			"III(B)-",
			"IV(AB)-",
			"Unknown",
		}
		var BloodTypesForFront = []string{
			"fierstPosetive",
			"secondPosetive",
			"thirdPosetive",
			"fourthPosetive",
			"fierstNegative",
			"secondNegative",
			"thirdNegative",
			"fourthNegative",
			"Unknown",
		}

		var countBloodTypes = make(map[string]int, 1)

		for i := 0; i < len(BloodTypes); i++ {
			// ? Bloodes
			var StatisticUser []structures.Signup
			cur, _ := collection.Find(ctx, bson.M{"blood": BloodTypes[i]})
			defer cur.Close(ctx)
			for cur.Next(ctx) {
				cur.Decode(&StatisticsUsers)
				StatisticUser = append(StatisticUser, StatisticsUsers)
			}
			countBloodTypes[BloodTypesForFront[i]] = len(StatisticUser)
		}
		// ======================================= Filter Users By Gender =====================================
		var genders = []string{
			"Мужской",
			"Женский",
		}
		var countGender = make(map[string]int, 1)
		for i := 0; i < len(genders); i++ {
			// ? Disabliaties
			var StatisticUser []structures.Signup
			cur, _ := collection.Find(ctx, bson.M{"gender": genders[i]})
			defer cur.Close(ctx)
			for cur.Next(ctx) {
				cur.Decode(&StatisticsUsers)
				StatisticUser = append(StatisticUser, StatisticsUsers)
			}
			countGender[genders[i]] = len(StatisticUser)
		}
		// ======================================= Disabilaties =====================================
		var disabilaties = []string{
			"1-й-степени",
			"2-й-степени",
			"3-й-степени",
			"Не имеется",
		}
		var disabilatiesForFront = []string{
			"first",
			"second",
			"third",
			"None",
		}

		var counDisabilaties = make(map[string]int, 1)

		for i := 0; i < len(disabilaties); i++ {
			// ? Disabliaties
			var StatisticUser []structures.Signup
			cur, _ := collection.Find(ctx, bson.M{"disabilaties": disabilaties[i]})
			defer cur.Close(ctx)
			for cur.Next(ctx) {
				cur.Decode(&StatisticsUsers)
				StatisticUser = append(StatisticUser, StatisticsUsers)
			}
			counDisabilaties[disabilatiesForFront[i]] = len(StatisticUser)
		}
		// ======================================= All Users =====================================
		c.JSON(200, gin.H{
			"Code":         "Request Handeled",
			"UsersAmount":  countUsers,
			"GenderMale":   countGender,
			"Bloodes":      countBloodTypes,
			"Disabliaties": counDisabilaties,
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
		ClientsDB   structures.Signup
		DoctorDB    structures.SignupDoctor
		Files       structures.File
		AdminDecode structures.Admin
	)
	// ==================== Cookie validation =========================
	CookieData := jwtgen.Velidation(c)

	log.Println(CookieData.Permissions)

	if CookieData.Permissions == "client" {
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
	} else if CookieData.Permissions == "admin" {
		Authenticationservice()
		collection := client.Database("MedCard").Collection("users")
		collectionFiles := client.Database("MedCard").Collection("ehrfiles")
		collection.FindOne(ctx, bson.M{"_id": c.Request.URL.RawQuery}).Decode(&ClientsDB)
		collection.FindOne(ctx, bson.M{"_id": CookieData.Id}).Decode(&AdminDecode)
		cur, err := collectionFiles.Find(ctx, bson.M{"clientid": c.Request.URL.RawQuery, "doctorid": CookieData.Id})
		if err != nil {
			fmt.Printf("error get fies%v", err)
		}
		defer cur.Close(ctx)
		var fiesArr []structures.File
		for cur.Next(ctx) {
			cur.Decode(&Files)
			fiesArr = append(fiesArr, Files)
		}

		if AdminDecode.Permissions == "admin" && ClientsDB.Lastname != "" {
			fmt.Printf("ClientsDB: %v\n", ClientsDB)
			if ClientsDB.Name != "" {
				ClientsDB.Password = "null"
				if len(fiesArr) != 0 {
					c.JSON(200, gin.H{
						"Code":     "Request Handeled",
						"UserJson": ClientsDB,
						"Files":    fiesArr,
					})
				} else {
					c.JSON(200, gin.H{
						"Code":     "Request Handeled",
						"UserJson": ClientsDB,
						"Files":    []string{},
					})
				}
			} else {
				c.JSON(400, gin.H{
					"Code": "User NotFound",
				})
			}
		} else if AdminDecode.Permissions == "admin" && c.Request.URL.RawQuery == "" {
			if AdminDecode.Name != "" {
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
		ctxForAccess, _ := context.WithTimeout(context.Background(), 5*time.Second)
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
				"UserJson": DoctorDB,
			})
		}
	}
}
