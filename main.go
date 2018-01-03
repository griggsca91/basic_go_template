package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func CreateSession(c *gin.Context) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(c.Request, "session-name")
	// Set some session values.

	log.Println(session)

	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	session.Save(c.Request, c.Writer)
}

func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	db := DB()
	log.Print(db)

	err := CreateSchema(db)
	if err != nil {
		panic(err)
	}
	game1 := &Game{
		Name: "ff7",
	}
	err = db.Insert(game1)
	if err != nil {
		panic(err)
	}

	user1 := &User{
		Username: "cgriggs",
		Email:    "cgriggs@gmail.com",
	}
	err = db.Insert(user1)
	if err != nil {
		panic(err)
	}

	user1.TrackGame(game1)

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Static("/public", "./public")
	r.LoadHTMLGlob("templates/*")

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired())
	{
		authorized.GET("/", homepageEndpoint)
	}

	r.POST("/signup", postSignupEndpoint)

	r.POST("/login", postLoginEndpoint)
	r.GET("/login", getLoginEndpoint)
	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}

func postLoginEndpoint(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	db := DB()
	var users []User
	err := db.Model(&users).
		Where("username = ?", username).
		Select()
	if err != nil {
		panic(err)
	}

	if len(users) == 0 {
		var errMsg struct {
			Success bool
			Error   string
		}
		errMsg.Success = false
		errMsg.Error = "Invalid Username/Password"
		c.JSON(http.StatusTeapot, errMsg)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].HashedPassword), []byte(password))

	var msg struct {
		Username       string
		Password       string
		User           User
		CompareSuccess bool
	}

	msg.Username = username
	msg.Password = password
	msg.User = users[0]
	msg.CompareSuccess = err == nil

	c.JSON(http.StatusOK, msg)
}

func getLoginEndpoint(c *gin.Context) {
	var msg struct {
		Endpoint string
	}
	msg.Endpoint = "loginEndpoint"

	c.JSON(http.StatusOK, msg)
}

func postSignupEndpoint(c *gin.Context) {

	username := c.PostForm("username")
	if username == "" {
		var errMsg struct {
			Success bool
			Error   string
		}
		errMsg.Success = false
		errMsg.Error = "Username is empty"
		c.JSON(http.StatusTeapot, errMsg)
		return
	}

	password := c.PostForm("password")
	if password == "" {
		var errMsg struct {
			Success bool
			Error   string
		}
		errMsg.Success = false
		errMsg.Error = "Password is empty"
		c.JSON(http.StatusTeapot, errMsg)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	msg := &User{
		Username:       username,
		Password:       password,
		Email:          "dkljf",
		HashedPassword: string(hashedPassword),
	}
	db := DB()

	err = db.Insert(msg)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, msg)
}

func homepageEndpoint(c *gin.Context) {
	var msg struct {
		Endpoint string
	}
	msg.Endpoint = "homepageEndpoint"

	c.JSON(http.StatusOK, msg)

}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("In the AuthRequired")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}
