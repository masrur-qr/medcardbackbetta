package evtvariables

import (
	"fmt"
	"os"
)

// var IpUrl = os.Getenv("IpUrl")
var IpUrl = "http://127.0.0.1:5173"
// var DBUrl = os.Getenv("DBUrl")
var DBUrl = "mongodb://127.0.0.1:27017"
var Port = os.Getenv("PORT")

func New()  {
	fmt.Printf("Port: %v;IpUrl %v;DBUrl: %v;\n", Port,IpUrl,DBUrl)
}