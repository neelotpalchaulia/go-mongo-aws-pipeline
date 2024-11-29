package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin" // A web framework for Go
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB collection variable
var collection *mongo.Collection

// Define a struct that represents a Task in MongoDB
type Task struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title"`
}

func main() {
	// Connect to MongoDB server
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017") // MongoDB URI
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create a MongoDB collection instance for CRUD operations
	collection = client.Database("taskdb").Collection("tasks")

	// Initialize Gin web framework for routing
	r := gin.Default()

	// Define routes
	r.POST("/tasks", createTask) // POST request to create a task
	r.GET("/tasks", getTasks)    // GET request to fetch tasks

	// Start the server
	r.Run(":8080") // The server listens on port 8080
}

// Create a task
func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil { // Bind incoming JSON data to the Task struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Insert task into MongoDB
	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task) // Respond with created task
}

// Get all tasks
func getTasks(c *gin.Context) {
	var tasks []Task
	cursor, err := collection.Find(context.TODO(), bson.M{}) // Find all tasks in MongoDB
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for cursor.Next(context.TODO()) {
		var task Task
		cursor.Decode(&task)        // Decode the cursor into a Task object
		tasks = append(tasks, task) // Add task to the tasks list
	}
	c.JSON(http.StatusOK, tasks) // Respond with the list of tasks
}
