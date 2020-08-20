package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/truongxuan1998/netflix/database"
	"github.com/truongxuan1998/netflix/models"
)

func Homepage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello")
}
func Watch(c echo.Context) error {
	return c.String(http.StatusOK, "Watch movie")
}
func Browser(c echo.Context) error {
	var movies []models.Movie
	database.Sess.Select("*").From(database.Movietable).Load(&movies)
	resp := new(models.ResponseMovie)
	resp.Movies = movies
	//////////////jwt//////////////////
	//	user := c.Get("user").(*jwt.Token)
	//	claims := user.Claims.(jwt.MapClaims)
	//	email := claims["email"].(string)
	//	authorized := claims["authorized"].(bool)
	//	message := fmt.Sprintf("Hello %s %v", email, authorized)
	//	return c.JSON(http.StatusOK, map[string]string{
	//		"Text": message,
	//	})
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("Email: ", claims["email"], "Authorized: ", claims["authorized"])
	return c.JSON(http.StatusOK, resp)
}
func TvSerial(c echo.Context) error {
	return c.String(http.StatusOK, "phim nhieu tap")
}

func Movies(c echo.Context) error {
	return c.String(http.StatusOK, "phim ")
}
func Latest(c echo.Context) error {
	return c.String(http.StatusOK, "phim moi nhat")
}
func CreateEndpoint(c echo.Context) error {
	req := new(models.Movie)
	if err := c.Bind(req); err != nil {
		return err
	}
	database.Sess.InsertInto(database.Movietable).Columns("id", "name", "genre").Values(req.ID, req.Name, req.Genre).Exec()
	return c.NoContent(http.StatusOK)
}
func PostList(c echo.Context) error {
	req := new(models.List)
	if err := c.Bind(req); err != nil {
		return err
	}
	database.Sess.InsertInto(database.Listtable).Columns("name", "id", "genre").Values(req.Name, req.ID, req.Genre).Exec()
	return c.NoContent(http.StatusOK)

}
func DeleteList(c echo.Context) error {
	id := c.Param("id")
	database.Sess.DeleteFrom(database.Listtable).Where("id = ?", id).Exec()
	return c.NoContent(http.StatusOK)
}
func ShowList(c echo.Context) error {
	var lists []models.List
	database.Sess.Select("*").From(database.Listtable).Load(&lists)
	resp := new(models.ResponseList)
	resp.Lists = lists
	return c.JSON(http.StatusOK, resp)
}
func SelectFilm(c echo.Context) error {
	var result models.Movie
	id := c.Param("id")
	database.Sess.Select("*").From(database.Movietable).Where("id= ?", id).Load(&result)
	return c.JSON(http.StatusOK, result)
}
func SearchHandler(c echo.Context) error {
	var result models.Movie
	Name := c.Param("name")
	database.Sess.Select("*").From(database.Movietable).Where("name= ?", Name).Load(&result)

	u, err := url.Parse(c.Request().URL.String())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	params := u.Query()
	searchkey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}
	fmt.Println("Search Query is:", searchkey)
	fmt.Println("Result page is:", page)
	return c.JSON(http.StatusOK, result)
}
func PostPerson(c echo.Context) error {
	req := new(models.Person)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	database.Sess.InsertInto(database.Persontable).Columns("name", "id").Values(req.Name, req.ID).Exec()
	return c.NoContent(http.StatusOK)
}
func PersonList(c echo.Context) error {
	var persons []models.Person
	database.Sess.Select("*").From(database.Persontable).Load(&persons)
	resp := new(models.ResponsePerson)
	resp.Persons = persons
	return c.JSON(http.StatusOK, resp)
}
func PersonInfo(c echo.Context) error {
	var result models.Person
	id := c.Param("id")
	database.Sess.Select("*").From(database.Persontable).Where("id=?", id).Load(&result)
	return c.JSON(http.StatusOK, result)
}
func DeletePerson(c echo.Context) error {
	id := c.Param("id")
	database.Sess.DeleteFrom(database.Persontable).Where("id = ?", id).Exec()
	return c.NoContent(http.StatusOK)
}
func SameContent(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	var result models.Movie
	Genre := c.Param("genre")
	database.Sess.Select("*").From(database.Movietable).Where("genre= ?", Genre).Load(&result)
	return json.NewEncoder(c.Response()).Encode(result)
}
func UpFile(c echo.Context) error {
	c.Request().ParseMultipartForm(100 << 20)
	// Retrieve the file from form data
	file, handler, err := c.Request().FormFile("upfile")
	if err != nil {
		return c.String(http.StatusBadRequest, "upfile error")
	}
	// Close the file when finish
	defer file.Close()
	// This is path which we want to store the file
	f, err := os.OpenFile("./public/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return c.String(http.StatusBadRequest, "Save error file")
	}
	defer f.Close()
	// Copy the file to the destination path
	io.Copy(f, file)
	fmt.Println("success: %v", handler.Header)
	return c.String(http.StatusOK, "Success")
}
