package controllers

import (
	"net/http"
	"server/apps/messaging/services"
	"server/errors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roomId struct {
	Id string `json:"id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	var rm *services.Room
	rmId := c.Param("id")

	rmObjectId, err := primitive.ObjectIDFromHex(rmId)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		errors.PrintError(errors.GetErrorKey(), errors.Wrap(err, err.Error()))
		return
	}

	rm, ok := services.GetRoom(rmObjectId)
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)

	services.ServeWs(rm, conn)
}

func NewRoom(c *gin.Context) {
	room := services.NewRoom()
	rmIdStream := make(chan string, 1)

	go room.Serve(rmIdStream)

	roomId := roomId{<-rmIdStream}

	c.JSON(http.StatusCreated, roomId)
}