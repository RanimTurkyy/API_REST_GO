package main
 
import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "sync"
    "time"
 
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
// Simuler une longue task
func task(id int) {
    log.Printf("Task %d running...\n", id)
    time.Sleep(3 * time.Second)
    log.Printf("Task %d completed!\n", id)
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
    })
 
    // Liste des tâches
    r.GET("/tasks", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "tasks": tasks,
        })
    })
 
    // Ajouter une tâche
    r.POST("/tasks", func(c *gin.Context) {
        var newTask Task
        if err := c.BindJSON(&newTask); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
            return
        }
 
        tasks = append(tasks, newTask)
        saveTasksToFile() // Sauvegarde après ajout
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
        saveTasksToFile() // Sauvegarde après modification
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
        saveTasksToFile() // Sauvegarde après suppression
        c.JSON(http.StatusNoContent, nil)
    })
 
    // /tasks/process
    r.POST("/tasks/process", func(c *gin.Context) {
        go func() { // Corrected the function declaration
            log.Println("Starting processing...")
            time.Sleep(5 * time.Second)
            log.Println("Processing done.")
        }()
 
        c.JSON(http.StatusOK, gin.H{
            "message": "Processing...",
        })
    })
    // Exécuter des tasks en parallèle
    r.GET("/tasks/parallel", func(c *gin.Context) {
        var wg sync.WaitGroup
 
        for i := 1; i <= 3; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                task(i)
            }()
        }
 
        wg.Wait()
        c.JSON(http.StatusOK, gin.H{"message": "All tasks completed!"})
    })
 
    r.Run() // Écoute sur 0.0.0.0:8080 (ou localhost:8080 sur Windows)
}