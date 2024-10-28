package handlers

import (
	"github.com/gin-gonic/gin"
)

//var output []map[string]string

func handlePing(c *gin.Context) {
	c.JSON(200, "PING ACTIVE")
}

func HandleRequests(router *gin.Engine) {
	router.GET("/ping", handlePing)
	//router.GET("/Outputs", fetchAllOutput)
	//router.POST("/restaurants", AddRestaurant)
}

//func fetchAllOutput(c *gin.Context) {
//    c.JSON(200, output)
//}//

/*func addRestaurant(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, "there is error")
		return
	}

	var restaurant []map[string]string

	err = json.Unmarshal(jsonData, &restaurant)
	if err != nil {
		c.JSON(500, "format error")
        return 0
	}

	restaurants = append(restaurants, restaurant...)
	c.JSON(200, "success")

}*/
