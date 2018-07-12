/*
	Parse JSON from a request
	The methods addCat, addDog and addHamster accomplish the same parsing,
	but Cat is fastest performance, Hamster is easiest to write.
 */

package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"io/ioutil"
	"github.com/labstack/gommon/log"
	"encoding/json"
	"github.com/labstack/echo/middleware"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

type Cat struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
}

type Dog struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
}

type Hamster struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
}

type JwtClaims struct {
	Name 		string 	  `json:"name"`
	jwt.StandardClaims
}


func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello from the web side!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nand his type is: %s\n", catName, catType))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "did you want JSON or String data?",
	})
}

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()
	//Read All method. b = body, err is the error if it gets one
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("failed reading the request body for addCat(): %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("failed unmarshaling in addCat(): %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("This is your cat: %#v\n", cat)
	return c.String(http.StatusOK, "we got your cat!")
}

func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("failed reading the request body for addDog(): %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("This is your dog: %#v\n", dog)
	return c.String(http.StatusOK, "we got your dog!")
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("failed reading the request body for addHamster(): %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("This is your dog: %#v\n", hamster)
	return c.String(http.StatusOK, "we got your hamster!")
}


func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "SECRET ADMIN PAGE ACCESSED")
}

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "you are on the secret cookie page!")
}

func mainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return c.String(http.StatusInternalServerError, "Oops, something went wrong")
	} else {
		log.Print("User Name: ", claims["name"]," ", " User ID: ", claims["jti"])
		return c.String(http.StatusOK, "you are on the secret JWT page!")
	}
	// claims is in interface, so if you want to use it as a string, you'll
	// need to typecast to string ( claims["name"].(string) ), but here we don't
	// have to because Print accepts interfaces.
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	//check username and password against database after hashing the password
	if username == "jack" && password == "1234" {
		cookie := &http.Cookie{}

		// this is the same
		//cookie := new(http.Cookie)

		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)

		c.SetCookie(cookie)

		// TODO: create JWT token creation
		token, err := createJwtToken()
		if err != nil {
			log.Print("Error Creating JWT token", err)
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "you were logged in",
			"token": token,
		})
	}
	return c.String(http.StatusUnauthorized, "Incorrect username or password")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"jack",
		jwt.StandardClaims{
			Id: "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}
	return token, nil
}


////////////////// CUSTOM MIDDLEWARES ///////////////////
//Adds Server name to response
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//use set to create our own Header name
		c.Response().Header().Set(echo.HeaderServer, "BlueBot/1.0")
		c.Response().Header().Set("notReallyHeader", "thisHasNoMeaning")

		return next(c)
	}
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")
		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present"){
				return c.String(http.StatusUnauthorized, "you don't have any cookie!")
			}
			//log.Print(err)
			return err
		}
		if cookie.Value == "some_string" {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "you don't have the right cookie")
	}
}


func main() {
	fmt.Println("Welcome to the Server")

	e:= echo.New()

	e.Use(ServerHeader)

	//can add middleware here: g := e.Group("/admin", middleware.Logger(), middleware...) instead of USE method used below
	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")
	jwtGroup := e.Group("/jwt")

	e.Use(middleware.Static("../static"))

	// this logs the server interaction
	//g.Use(middleware.Logger())
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:`[${time_rfc3339}]  ${status}  ${method}  ${host}${path}  ${latency_human}` + "\n",
	}))

	adminGroup.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		if username == "jack" && password == "1234" {
			return true, nil
		}
		return false, nil
	}))

	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey: []byte("mySecret"),
	}))

	cookieGroup.Use(checkCookie)

	cookieGroup.GET("/main", mainCookie)

	adminGroup.GET("/main", mainAdmin) // can also add middleware directly to a method like GET or POST after the handler (mainAdmin)

	jwtGroup.GET("/main", mainJwt)

	e.GET("/login", login)
	e.GET("/hello", hello)
	e.GET("/cats/:id", getCats) //the colon defines the name of the end parameter

	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)

	e.Start(":8000")
}


