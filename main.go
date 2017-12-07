package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gorilla/sessions"
	"log"
	"net/http"
	tb "twitchboard/twitchboard"
	"github.com/go-pg/pg"
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

func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{&tb.User{}} {
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil

}


func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	db := tb.DB()
	log.Print(db)

	err := CreateSchema(db)
	if err != nil {
		panic(err)
	}
	user1 := &tb.User{
		Username:   "admin",
		Email: "admin1@admin",
	}
	err = db.Insert(user1)
	if err != nil {
		panic(err)
	}

	var users []tb.User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	log.Println(users)


	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

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

	r.POST("/login", postLoginEndpoint)
	r.GET("/login", getLoginEndpoint)
	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}

func postLoginEndpoint(c *gin.Context) {



	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getLoginEndpoint(c *gin.Context) {
	var msg struct {
		Endpoint string
	}
	msg.Endpoint = "loginEndpoint"

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