package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var tasks = []Task{}

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.GET("/tasks", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "tasks": tasks,
    })
  })
  r.POST("/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.BindJSON(&newTask); err != nil {
			return
		}
		tasks = append(tasks, newTask)
		c.IndentedJSON(http.StatusCreated, newTask)
	})
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}