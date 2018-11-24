package main

import (
	// "os"
	"fmt"
	"strings"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

func main() {
	// Echo instance
	e := echo.New()
	
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// mongoHost := os.Getenv("mongoHost")
	// mongoHost := os.Getenv("port")
	mongoHost := viper.GetString("Mongo.Host")
	mongoUser := viper.GetString("Mongo.User")
	mongoPass := viper.GetString("Mongo.Password")
	port := viper.GetString("port")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	session, err := mgo.Dial(mongoUser + ":"+ mongoPass + "@" + mongoHost)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	h := handler {
		m: session,
	}
	// Routes
	e.GET("/", hello)
	e.GET("/todos/:id", h.view)
	e.GET("/todos", h.list)
	e.PUT("/todos/:id", h.done)
	e.POST("/todos", h.create)
	e.DELETE("/todos/:id", h.delete)

	// Start server
	e.Logger.Fatal(e.Start(port))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type todo struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
	Topic string `json:"topic" bson:"topic"`
	Done bool `json:"done" bson:"done"`
}

type handler struct {
	m *mgo.Session
}

func (h *handler)view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo

	col := session.DB("workshop").C("todosNaja")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handler)list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var ts []todo

	col := session.DB("workshop").C("todosNaja")
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ts)
}

func (h *handler)create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	t.ID = bson.NewObjectId()
	col := session.DB("workshop").C("todosNaja")
	if err := col.Insert(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handler)done(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo

	col := session.DB("workshop").C("todosNaja")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	t.Done = true
	if err := col.UpdateId(id, t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handler)delete(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo

	col := session.DB("workshop").C("todosNaja")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	
	if err := col.RemoveId(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result" : "Remove Success",
	})
}