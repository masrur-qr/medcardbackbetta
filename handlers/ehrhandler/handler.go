package ehrhandler

import (
	"medcard-new/begening/controllers/ehrcontroller"

	"github.com/gin-gonic/gin"
)

func Handler(){
	r := gin.Default()

	r.POST("/handleviews",ehrcontroller.DoctorClientForView)
	r.POST("/filesadd",ehrcontroller.AddFilesToEhr)

	r.Run(":5502")
}