package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	port, success := os.LookupEnv("EXPOSE_PORT")
	if !success {
		log.Println("Unable to read EXPOSE_PORT environment variable. Defaulting to port 8080")
		port = "8080"
	}

	app := App{router: SetUpRouter()}
	app.Run(fmt.Sprintf(":%s", port))
}
