package controllers

import (
	"net/http"
	"toko_sembako_acen/models"
	"toko_sembako_acen/repository"

	"github.com/gin-gonic/gin"
)

func GetData(ctx *gin.Context) {
	var example []*models.Example
	repository.Get(&example)
	ctx.JSON(http.StatusOK, &example)

}
