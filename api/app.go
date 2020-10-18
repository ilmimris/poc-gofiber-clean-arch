package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_authorRepo "github.com/ilmimris/poc-gofiber-clean-arch/pkg/author/repository/psql"
	_postDelivery "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/delivery/rest"
	_postRepo "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/repository/psql"
	_postUsecase "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/usecase"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func databaseConnection() (*sql.DB, error) {
	// Read db configuration
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	// Connect with database
	dbConn, err := sql.Open(`postgres`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return dbConn, err
}

func main() {
	db, err := databaseConnection()
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatalf("Database Connection error: %s", err)
	}
	fmt.Println("Database connection success!")

	postRepo := _postRepo.NewPsqlPostRepository(db)
	authorRepo := _authorRepo.NewPsqlAuthorRepository(db)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	postUcase := _postUsecase.NewPostUsecase(postRepo, authorRepo, timeoutContext)

	// Create a Fiber app
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to the clean-architecture!"))
	})

	_postDelivery.NewPostHandler(app, postUcase)

	_ = app.Listen(":8080")

}
