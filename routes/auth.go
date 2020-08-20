package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/truongxuan1998/netflix/database"
	"github.com/truongxuan1998/netflix/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}
	//	log.Printf("req data = %v", req)
	/////////////Check email and password////////////
	emailcheck := models.IsEmpty(req.Email)
	passwordcheck := models.IsEmpty(req.Password)
	if emailcheck || passwordcheck {
		return c.String(http.StatusBadRequest, "Data can't empty")
	}
	if !govalidator.IsEmail(req.Email) {
		return c.String(http.StatusBadRequest, "Email is invalid")
	}
	req.Email = models.Santize(req.Email)
	req.Password = models.Santize(req.Password)
	///////////////////Signup, hashedPassword user/////////////////////
	var user string
	err := database.Sess.QueryRow("SELECT email FROM user where email =?", req.Email).Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.String(http.StatusInternalServerError, "unable to create your account")
		}
		// post email and hashedpassword
		_, err = database.Sess.Exec("Insert into user(email,password) VALUES(?,?)", req.Email, hashedPassword)
		if err != nil {
			return c.String(http.StatusInternalServerError, "unable to create your account")
		}
		// Send email and hashedpassword OK
		return c.String(http.StatusOK, "Signup successful")
	default:
		return c.String(http.StatusConflict, "Email does exist")
	}
}
func Login(c echo.Context) error {
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}
	log.Printf("req data = %v", req)
	/////////////Check email and password////////////
	emailcheck := models.IsEmpty(req.Email)
	passwordcheck := models.IsEmpty(req.Password)
	if emailcheck || passwordcheck {
		return c.String(http.StatusBadRequest, "Data can't empty")
	}
	if !govalidator.IsEmail(req.Email) {
		return c.String(http.StatusBadRequest, "Email is invalid")
	}
	var dbemail string
	var dbpassword string
	err := database.Sess.QueryRow("SELECT email, password from user where email=?", req.Email).Scan(&dbemail, &dbpassword)
	if err != nil {
		return c.String(http.StatusBadRequest, "Email not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(req.Password))
	if err != nil {
		return c.String(http.StatusBadRequest, "Password were wrong")
	}
	////////////createJWT///////////////
	token, err := models.CreateJWT(req.Email)
	if err != nil {
		log.Println("Error creating JWT", err)
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"Message": "You were logged in",
		"Token":   token,
	})

	//	email := c.Request().FormValue("email")
	//	password := c.Request().FormValue("password")
	//	email := c.QueryParam("email")
	//	password := c.QueryParam("password")
	//	emailcheck := IsEmpty(email)
	//	passwordcheck := IsEmpty(password)
	//	if emailcheck || passwordcheck {
	//		return c.String(http.StatusBadRequest, "Data can't empty")
	//	}
	//	if !govalidator.IsEmail(email) {
	//		return c.String(http.StatusBadRequest, "Email is invalid")
	//	}
	//	if email == "ntcn14089@gmail.com" && password == "netflix123" {
	////////////createJWT///////////////
	//		token, err := createJWT(email)
	//		if err != nil {
	//			log.Println("Error creating JWT", err)
	//			return c.String(http.StatusInternalServerError, "Something went wrong")
	//		}
	//		return c.JSON(http.StatusOK, map[string]string{
	//			"Message": "You were logged in",
	//			"Token":   token,
	//		})
	//	}
	//	return c.String(http.StatusUnauthorized, "Your email or password is incorrect")
}
