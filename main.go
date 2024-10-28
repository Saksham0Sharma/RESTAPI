package main

import (
	handlers "myapiproj/Handlers"
	"myapiproj/fileparser"

	"github.com/gin-gonic/gin"
)

var output []map[string]string

func fetchAllOutput(c *gin.Context) {
	c.JSON(200, output)
}

func main() {

	route := gin.Default()
	route.GET("/output/all", fetchAllOutput)
	handlers.HandleRequests(route)
	fileparser.HandleCSV_ExcelParsing(route)
	route.Run()

}

//route.GET("/ping", handlePing)
////route.GET("/output",output)
//route.POST("/restaurant", addRestaurant)
//route.POST("/parsefile", parseFile)
//fileparser.HandleFileParsing(route)
