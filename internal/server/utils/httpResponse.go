package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpResponder struct {
	c *gin.Context
}

func NewHttpResponder(c *gin.Context) HttpResponder {
	return HttpResponder{c}
}

func (httpr *HttpResponder) ResponseError(errMsg string) {
	httpr.c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": errMsg,
		"data":    nil,
	})
}

func (httpr *HttpResponder) InvalidNumPlayerParam(numPlayersParam string) {
	httpr.ResponseError(fmt.Sprintf("Invalid parameter for number of player (%s)", numPlayersParam))
}

func (httpr *HttpResponder) MaxNumPlayer(maxPlayers int) {
	httpr.ResponseError(fmt.Sprintf("Number of players must not exceed %d", maxPlayers))
}

func (httpr *HttpResponder) RoomNotFound(idParam string) {
	httpr.ResponseError(fmt.Sprintf("Room not found (%s)", idParam))
}

func (httpr *HttpResponder) RoomFull(id string) {
	httpr.ResponseError(fmt.Sprintf("Room with id: %s is full", id))
}

func (httpr *HttpResponder) ResponseOk(message string, data interface{}) {
	httpr.c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func (httpr *HttpResponder) RoomCreated(id string) {
	httpr.ResponseOk("Room created with id: "+id, gin.H{
		"id": id,
	})
}

func (httpr *HttpResponder) RoomAvailable(id string) {
	httpr.ResponseOk("Room with id: "+id+" is available", nil)
}
