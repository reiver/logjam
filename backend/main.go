package main

import (
	"github.com/sparkscience/logjam/backend/srv/log"
	
	"fmt"
)

func main() {
	log := logsrv.Begin()
	defer log.End()
	
	fmt.Println("Hello world! — I am logjam! ✨")
}
