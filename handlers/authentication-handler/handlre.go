package authenticationhandler

import (
	"medcard-new/begening/controllers/authenticationservice"

	"github.com/gin-gonic/gin"
)

func AuthHandler(){
	r := gin.Default()

	r.POST("/signup",authenticationservice.Signup)
	r.POST("/signin",authenticationservice.Signin)
	r.POST("/signupdoctor",authenticationservice.SignupDoctor)

	r.Run(":5501")
}