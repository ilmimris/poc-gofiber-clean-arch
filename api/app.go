package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_authorRepoMysql "github.com/ilmimris/poc-gofiber-clean-arch/pkg/author/repository/mysql"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"

	_authorRepoPsql "github.com/ilmimris/poc-gofiber-clean-arch/pkg/author/repository/psql"
	_postDelivery "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/delivery/rest"
	_postRepoMysql "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/repository/mysql"

	_postRepoPsql "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/repository/psql"
	_postUsecase "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/usecase"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
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

func databaseConnectionPsql() (*sql.DB, error) {
	// Read db configuration
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	val := url.Values{}

	// Config db postgres
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
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

func databaseConnectionMysql() (*sql.DB, error) {
	// Read db configuration
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	val := url.Values{}

	// Config db mysql
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	// Connect with database
	dbConn, err := sql.Open(`mysql`, dsn)

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
	dbKind := viper.GetString(`database.kind`)

	var db *sql.DB
	var err error

	switch dbKind {
	case "mysql":
		db, err = databaseConnectionMysql()
	case "postgres":
		db, err = databaseConnectionPsql()
	}
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

	var postRepo domain.PostRepository
	var authorRepo domain.AuthorRepository

	switch dbKind {
	case "mysql":
		postRepo = _postRepoMysql.NewMysqlPostRepository(db)
		authorRepo = _authorRepoMysql.NewMysqlAuthorRepository(db)
	case "postgres":
		postRepo = _postRepoPsql.NewPsqlPostRepository(db)
		authorRepo = _authorRepoPsql.NewPsqlAuthorRepository(db)
	}

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	postUcase := _postUsecase.NewPostUsecase(postRepo, authorRepo, timeoutContext)

	// Create a Fiber app
	app := fiber.New()
	app.Use(cors.New())

	// Use loggoer middleware
	app.Use(logger.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to the clean-architecture!"))
	})

	_postDelivery.NewPostHandler(app, postUcase)

	log.Fatal(app.Listen(viper.GetString(`server.address`)))

}
