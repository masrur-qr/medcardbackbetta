package apihandler

import (
	"medcard-new/begening/controllers/authenticationservice"
	"medcard-new/begening/controllers/ehrcontroller"
	controllers "medcard-new/begening/controllers/handlersservice"

	"github.com/gin-gonic/gin"
)

func Handlers(){
	r := gin.Default()

	r.POST("/insertquestion",controllers.InsertQuestions)
	r.POST("/profilechange",controllers.ProfileChange)

	// ================================== New Route ==============================

	r.POST("/signup",authenticationservice.Signup)
	r.POST("/signin",authenticationservice.Signin)
	r.POST("/signupdoctor",authenticationservice.SignupDoctor)
	r.POST("/handleviews",ehrcontroller.DoctorClientForView)
	r.POST("/filesadd",ehrcontroller.AddFilesToEhr)
	
	r.GET("/getclient",controllers.GetClient)
	r.GET("/getquestion",controllers.GetQuestions)
	r.GET("/getdoctors",controllers.GetDoctors)
	r.GET("/statistics",controllers.Statistics)
	r.GET("/getclients",controllers.GetClients)
	r.GET("/getviews",controllers.GetViews)

	r.Run(":5500")
}