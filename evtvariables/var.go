package evtvariables

import (
	"fmt"
	"os"
)

var IpUrl = os.Getenv("IpUrl")
var DBUrl = os.Getenv("DBUrl")
var Port = os.Getenv("PORT")

func New()  {
	fmt.Printf("Port: %v;IpUrl %v;DBUrl: %v;\n", Port,IpUrl,DBUrl)
}