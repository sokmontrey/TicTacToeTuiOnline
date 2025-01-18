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

func (httpr *HttpResponder) Error(errMsg string) {
	httpr.c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": errMsg,
	})
}

func (httpr *HttpResponder) NumPlayersParamError(numPlayersParam string) {
	httpr.Error(fmt.Sprintf("Invalid parameter for number of player (%s)", numPlayersParam))
}

func (httpr *HttpResponder) MaxNumPlayerError(maxPlayers int) {
	httpr.Error(fmt.Sprintf("Number of players must not exceed %d", maxPlayers))
}

func (httpr *HttpResponder) RoomCreated(id string) {
	httpr.c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"id":      id,
	})
}
