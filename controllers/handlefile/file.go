package handlefile

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Handlefile(c *gin.Context,directory string)string{
	// """"""""""""""""""get the img""""""""""""""""""
	// upload of 10MB files
	// c.Request.ParseMultipartForm(1000 << 50)
	c.Request.ParseMultipartForm(5 * 1024 * 1024)
	// formFiles haeders
	files, handler, err := c.Request.FormFile("img")
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer files.Close()
	fmt.Printf("File Name %s\n", strings.Split(handler.Filename, "."))
	format := strings.Split(handler.Filename, ".")
	fmt.Printf("format[len(format)-1]: %v\n", format[len(format)-1])
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
	log.Printf("file1%v",strings.Split(tempFiles.Name(), "upload"))
	spliString := strings.Split(tempFiles.Name(), "upload")
	idString := spliString[len(spliString)-1]
	return idString
}