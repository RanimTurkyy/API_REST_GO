package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var tasks = []Task{}

// Fonction pour sauvegarder les tâches dans tasks.json
func saveTasksToFile() {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatalf("Erreur de sérialisation JSON: %v", err)
	}
	err = ioutil.WriteFile("tasks.json", data, 0644)
	if err != nil {
		log.Fatalf("Erreur d'écriture dans le fichier: %v", err)
	}
}

// Fonction pour charger les tâches depuis tasks.json
func loadTasksFromFile() {
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Printf("Aucun fichier trouvé, démarrage avec une liste vide.")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Erreur de lecture du fichier: %v", err)
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatalf("Erreur de désérialisation JSON: %v", err)
	}
}

// Fonction pour récupérer une tâche par ID
func getTaskByID(id string) (*Task, int) {
	for i, t := range tasks {
		if t.ID == id {
			return &t, i
		}
	}
	return nil, -1
}

func main() {
	// Charger les tâches depuis le fichier JSON au démarrage
	loadTasksFromFile()

	r := gin.Default()

	// Route de test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
