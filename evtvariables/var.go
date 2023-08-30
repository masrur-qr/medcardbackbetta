package evtvariables

import (
	"os"
)

var IpUrl = os.Getenv("IpUrl")
var DBUrl = os.Getenv("DBUrl")
var Port = os.Getenv("PORT")

// var IpUrl = "http://192.168.11.41:3000"
// var DBUrl = "mongodb://root:14vDuB2YdS@192.168.1.101:31342"
// var Port = "5500"
