package apihandler

import (
	"crypto/tls"
	"fmt"
	"medcard-new/begening/controllers/authenticationservice"
	"medcard-new/begening/controllers/ehrcontroller"
	controllers "medcard-new/begening/controllers/handlersservice"
	gomail "gopkg.in/mail.v2"
	"github.com/gin-gonic/gin"
)

func Handlers(){	
	r := gin.Default()

	// r.StaticFS("/static", gin.Dir("./static", true))
	r.StaticFS("/static", gin.Dir("./static", false))
	r.Use(controllers.Cors)

	r.POST("/insertquestion",controllers.InsertQuestions)
	r.POST("/profilechange",controllers.ProfileChange)
	r.GET("/link",ehrcontroller.ExpiredLinks)
	// ================================== New Route ==============================

	r.POST("/signup",authenticationservice.Signup)
	r.POST("/reset",authenticationservice.ResetPassword)
	r.POST("/signin",authenticationservice.Signin)
	r.POST("/signout",authenticationservice.Signout)
	r.POST("/logincheck",authenticationservice.LoginCheck)
	r.POST("/signupdoctor",authenticationservice.SignupDoctor)
	r.POST("/handleviews",ehrcontroller.DoctorClientForView)
	r.POST("/filesadd",ehrcontroller.AddFilesToEhr)
	
	r.GET("/getclient",controllers.GetClient)
	r.GET("/getquestion",controllers.GetQuestions)
	r.GET("/getdoctors",controllers.GetDoctors)
	r.GET("/statistics",controllers.Statistics)
	r.GET("/getclients",controllers.GetClients)
	r.GET("/getviews",controllers.GetViews)
	r.GET("/listviews",controllers.ListViewsAdmin)

	// ? send mail if route is out of range

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		m := gomail.NewMessage()

		// Set E-Mail sender
		m.SetHeader("From", "notification@medcard.space")
	  
		// Set E-Mail receivers
		m.SetHeader("To", "munosibshoev@khorog.dev")
	  
		// Set E-Mail subject
		m.SetHeader("Subject", "Worrning !! SomeOne trying to do something")
	  
		// Set E-Mail body. You can set plain text or html with text/html\
		text := fmt.Sprintf("This mail was sent becouse someone  sending different request to us check your logs and block it IP: %v and PATH: %v and HEADERS: %v ",c.RemoteIP(),c.Request.URL.Path,c.Request.Header)
		m.SetBody("text/plain", text)
	  
		// Settings for SMTP server
		d := gomail.NewDialer("smtp.beget.com", 2525, "notification@medcard.space", "Yit&Iak0")
	  
		// This is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	  
		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
		  fmt.Println(err)
		  panic(err)
		}
	  
		return
	})

	
	r.Run(":5500")
}
