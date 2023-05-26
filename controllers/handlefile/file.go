package handlefile

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func Handlefile(c *gin.Context,directory string)string{
	// """"""""""""""""""get the img""""""""""""""""""
	// upload of 10MB files
	c.Request.ParseMultipartForm(1000 << 50)
	// formFiles haeders
	files, handler, err := c.Request.FormFile("img")
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer files.Close()
	fmt.Printf("File Name %s\n", handler.Filename)
	// create temporary files within the folder
	tempFiles, err := ioutil.TempFile(directory, "upload-*.jpg")
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer tempFiles.Close()
	// read all files to upload
	fileByte, err := ioutil.ReadAll(files)
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	// write  the byte arrey into temp files
	tempFiles.Write((fileByte))
	log.Printf("file1%v",strings.Split(tempFiles.Name(), "upload"))
	spliString := strings.Split(tempFiles.Name(), "upload")
	idString := spliString[len(spliString)-1]
	return idString
}