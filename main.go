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

func getTaskByID(id string) (*Task, int) {
	for i, t := range tasks {
		if t.ID == id {
			return &t, i
		}
	}
	return nil, -1
}


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

  // Modifier une tâche
	r.PUT("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedTask Task

		if err := c.BindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
			return
		}

		task, index := getTaskByID(id)
		if task == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
			return
		}

		// Mise à jour de la tâche
		tasks[index].Title = updatedTask.Title
		c.JSON(http.StatusOK, tasks[index])
	})

	// Supprimer une tâche
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, index := getTaskByID(id)
		if index == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
			return
		}

		// Suppression de la tâche
		tasks = append(tasks[:index], tasks[index+1:]...)
		c.JSON(http.StatusNoContent, nil)
	})

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}