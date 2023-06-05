package evtvariables

import (
	"os"
)

var IpUrl = os.Getenv("IpUrl")
var DBUrl = os.Getenv("DBUrl")
var Port = os.Getenv("PORT")
