package main

import (
	"context"
	"fmt"
	"log"
	"server/config"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// "gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
)

type SqlLogger struct {
	gormLogger.Interface
}

func (l SqlLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64), err error,
) {
	sql, _ := fc()
	fmt.Printf("%v\n==============================\n", sql)
}

var db *gorm.DB

func main() {
	// var err error

	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
    AllowOrigins: os.Getenv("FRONTEND_ENTPOINT"),
    AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}))
	app.Use(func (c *fiber.Ctx) error  {
		idToken := c.Get("Authorization")
		firebase, err := config.InitializeFirebase();
		if err != nil {
			panic(err)
		}
	
		ctx := context.Background()
		client, err := firebase.Auth(ctx)
	
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "unexpected error",
			})
		}
	
		token, err := client.VerifyIDToken(ctx, idToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized error",
			})
		}
	
		data := token.Claims;
		c.Locals("data", data)
		c.Next()

		return nil
	})

	// dial := mysql.Open("root@tcp(127.0.0.1:3306)/firebase-auth?parseTime=true")
	// db, err = gorm.Open(dial, &gorm.Config{
	// 	Logger: &SqlLogger{},
	// 	DryRun: false,
	// })
	// if err != nil {
	// 	panic(err)
	// }


	app.Get("/", func(c *fiber.Ctx) error {
		data := c.Locals("data").(map[string]interface{})
		fmt.Println(data["name"])

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "done",
		})
	})

	port := fmt.Sprintf(":%v", os.Getenv("API_PORT"))
	app.Listen(port)

	// user := User{
	// 	Name: data["name"].(string),
	// 	Email: data["email"].(string),
	// 	ImageURL: data["picture"].(string),
	// 	UserId: data["user_id"].(string),
	// }
	// CreateUser(user)

	// log.Printf("Verified ID token: %v\n", data["email"])

}

func InitializeDatabase() (*gorm.DB, error) {
	database_uri := fmt.Sprintf(
		"%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	dial := mysql.Open(database_uri)
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

type User struct {
  ID       uint    `gorm:"primaryKey"`
  Name     string  `gorm:"column:name"`
  Email    string  `gorm:"column:email"`
  ImageURL string  `gorm:"column:imageURL"`
  UserId   string  `gorm:"column:userId"` 
}


func CreateUser(user User) {
	tx := db.Create(&user)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	fmt.Println(user)
}