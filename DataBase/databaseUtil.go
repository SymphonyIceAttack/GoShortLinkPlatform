package databaseUtil

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LinkObject struct {
	ID        uint   `gorm:"primaryKey"`
	WholeUrl  string `gorm:"uniqueIndex"`
	ShortUrl  string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func LoadDataBase() (*gorm.DB, error) {
	godotenv.Load("../.env.local")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := "5432"
	dbname := os.Getenv("POSTGRES_DATABASE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s slmode=require TimeZone=Asia/Shanghai", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {

		return nil, err
	}

	return db, nil

}
