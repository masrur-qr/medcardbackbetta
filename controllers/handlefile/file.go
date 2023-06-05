package handlefile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Handlefile(c *gin.Context,directory string)string{
	// """"""""""""""""""get the img""""""""""""""""""
	c.Request.ParseMultipartForm(10 * 1024 * 1024)
	// formFiles haeders
	files, handler, err := c.Request.FormFile("img")
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer files.Close()
	format := strings.Split(handler.Filename, ".")
	var tempFiles *os.File
	if format[len(format)-1] == "jpg" || format[len(format)-1] == "JPG" || format[len(format)-1] == "Jpg" {
		tempFiles, err = ioutil.TempFile(directory, "upload-*.jpg")
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
	}else if format[len(format)-1] == "pdf" || format[len(format)-1] == "PDF" || format[len(format)-1] == "Pdf" {
		tempFiles, err = ioutil.TempFile(directory, "upload-*.pdf")
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
	}else {
		tempFiles, err = ioutil.TempFile(directory, "upload-*.jpg")
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
	}
	// create temporary files within the folder
	defer tempFiles.Close()
	// read all files to upload
	fileByte, err := ioutil.ReadAll(files)
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	// write  the byte arrey into temp files
	tempFiles.Write((fileByte))
	spliString := strings.Split(tempFiles.Name(), "upload")
	idString := spliString[len(spliString)-1]
	return idString
}