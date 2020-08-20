package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/truongxuan1998/netflix/routes"
)

func main() {
	fmt.Println("Welcome")
	api := echo.New()

	////////////////middlewares////////////////////
	api.Use(middleware.Logger())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}))
	api.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "static",
	}))
	api.Use(middleware.Recover())
	//	log := middleware.Logger()
	isLoggedIn := middleware.JWT([]byte("SecretJWT"))
	/////////////////// routes////////////////////
	api.GET("/browser", routes.Browser, isLoggedIn)
	api.GET("/browser/genre/83", routes.TvSerial, isLoggedIn)
	api.GET("/browser/genre/3499", routes.Movies, isLoggedIn)
	api.GET("/browser/mylist", routes.ShowList, isLoggedIn)
	api.GET("/browser/person/:id", routes.PersonInfo, isLoggedIn)
	api.GET("/browser/person", routes.PersonList, isLoggedIn)
	api.GET("/", routes.Homepage)
	api.GET("/watch", routes.Watch)
	api.GET("/latest", routes.Latest)
	api.GET("/browser/:id", routes.SelectFilm, isLoggedIn)
	api.GET("/browser/genre/:genre", routes.SameContent, isLoggedIn)
	api.GET("/search/:name", routes.SearchHandler)

	api.POST("/login", routes.Login)
	api.POST("/signup", routes.Signup)
	api.POST("/movie", routes.CreateEndpoint)
	api.POST("/browser/mylist", routes.PostList, isLoggedIn)
	api.POST("/browser/person", routes.PostPerson)
	api.POST("/upfile", routes.UpFile, isLoggedIn)
	api.DELETE("/browser/mylist/:id", routes.DeleteList, isLoggedIn)
	api.DELETE("/person/:id", routes.DeletePerson, isLoggedIn)
	api.Logger.Fatal(api.Start(":8080"))
}
