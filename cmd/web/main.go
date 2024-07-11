package main

import (
	"log"
	"os"

	rank "github.com/eugene/iizi_errand"
	"github.com/eugene/iizi_errand/pkg/models/psql"
	"github.com/gofiber/fiber/v2"
)

var (
	infoLogger = log.New(os.Stdout, "[iizi]	[INFO]: \t", log.Ldate | log.Ltime | log.Lshortfile)
	errorLogger = log.New(os.Stdout, "[iizi] [ERROR]: \t", log.Ldate | log.Ltime | log.Lshortfile)
)


func main() {
	cat := []string{
		"name", "category", "casley", "boge", "cook", "cleaner", "laundry", "dog walking", "grocery shopping",
	}
	cat1 := []string{
		"name", "casley",
	}
	rank.RankRunner(cat, cat1)
	app := fiber.New( // initialize the new app instance
		fiber.Config{
			CaseSensitive: true,
			// StrictRouting: true,
			ServerHeader:  "Fiber",
			AppName: "iiziErand v1.0.1",
		},
	)
	// run DB migrations
	engine, err:= psql.DBConnection()
	if err != nil {
		errorLogger.Println(err)
	}

	r := Repository{
		DBConn: engine,
	}
	r.Routes(app)
	
	// serve on port 3000
	infoLogger.Println("Server is running on port 3000")
	err = app.Listen(":3000")
	if err != nil {
		errorLogger.Println(err)
	}
}

