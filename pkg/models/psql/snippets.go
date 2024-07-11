package psql

import (
	"fmt"
	"os"

	"github.com/eugene/iizi_errand/pkg/models"
	_ "github.com/lib/pq"
	"xorm.io/xorm"

	"github.com/joho/godotenv"
)


func DBConnection() (*xorm.Engine, error) {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// connect to database
	dsn := fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, os.Getenv("DBUSERNAME"), os.Getenv("PASSWORD"), "iizidb")
	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil{
		return nil, err
	}
	if err := engine.Ping(); err != nil{
		return nil, err
	}

	if err := engine.Sync( // migrate tables to DB
			new(models.UserModel),
			new(models.TaskModel),
			new(models.RatingModel),
			new(models.UserProfile),
			new(models.ErrandRunnerProfile),
			new(models.LoginData),
			new(models.ErrandApplication),
		); err != nil{
		return nil, err
	}
	return engine, err
}
