package main

import (
	"log"
	"os"
	"strconv"
	"time"

	rank "github.com/eugene/iizi_errand"
	"github.com/eugene/iizi_errand/pkg/models/psql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/joho/godotenv"
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
			EnableTrustedProxyCheck: true,
		},
	)

	// cache middleware
	app.Use(cache.New(cache.Config{
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			newCacheTime, _ := strconv.Atoi(c.GetRespHeader("Cache-Time", "600"))
			return time.Second * time.Duration(newCacheTime)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
	}))

	// recover from panic error
	// Initialize default config
	app.Use(recover.New())

	// rate limit
	app.Use(limiter.New(limiter.Config{
		Max:        2,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Use IP address as the key
		},
	}))

	
	// load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		return
	}

	// run DB migrations
	engine, err:= psql.DBConnection()
	if err != nil {
		errorLogger.Println(err)
	}

	r := Repository{
		DBConn: engine,
	}
	r.Routes(app)
	
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "11750" // Default port if not specified in the environment
	}
	// serve on port 3000
	err = app.Listen("0.0.0.0:"+PORT)
	if err != nil {
		errorLogger.Println(err)
	}
	infoLogger.Println("Server is running on port %s", PORT)
}


