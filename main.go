package main

import (
	"demo-echo/driver"
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_articleHttpDeliveryMiddleware "demo-echo/article/delivery/http/middleware"
	_articleRepo "demo-echo/article/repository/mysql"
	_articleUcase "demo-echo/article/usecase"
	_authorRepo "demo-echo/author/repository/mysql"

	_articleHttpDelivery "demo-echo/article/delivery"
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

func main() {
	db, err := driver.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = db.SQL.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.SQL.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	artMiddL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(artMiddL.CORS)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(db.SQL.DB())
	articleRepo := _articleRepo.NewMysqlArticleRepository(db.SQL.DB())

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	articleUsecase := _articleUcase.NewArticleUsecase(articleRepo, authorRepo, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, articleUsecase)
	log.Fatal(e.Start(viper.GetString("server.address")))
}
