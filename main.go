package main

import (
	"net/http"
	"strconv"

	"assignment-todos/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type HTTPError struct {
	StatusCode int `json:"code"`
	Message string `json:"message"`
}

var todos = [] Todo{
	{ID: 1, Description: "Buy groceries", Completed: false},
	{ID: 2, Description: "Clean the house", Completed: true},
	{ID: 3, Description: "Go for a run", Completed: false},
	{ID: 4, Description: "Study for exams", Completed: false},
	{ID: 5, Description: "Write a blog post", Completed: true},
}

var (
	errorBadRequest     = newHTTPError(http.StatusBadRequest, "Bad Request")
	errorDataNotFound   = newHTTPError(http.StatusNotFound, "Todo not found")
)

func newHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{StatusCode: code, Message: msg}
}

// @title My API TODOS without Database Example
// @version 1.0
// @description This is a sample API that demonstrates CRUD operations for managing TODO items without using a database. It is built using Gin-Gonic framework for handling HTTP requests and responses, and Swagger for API documentation.
// @host localhost:8080
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.GET("/todos", getAllTodos)          // Endpoint to get all TODO items.
	r.GET("/todos/:id", getTodoByID)      // Endpoint to get a TODO item by its ID.
	r.POST("/todos", addTodo)             // Endpoint to add a new TODO item.
	r.PUT("/todos/:id", editTodo)         // Endpoint to edit an existing TODO item by its ID.
	r.DELETE("/todos/:id", deleteTodo)    // Endpoint to delete a TODO item by its ID.

	// Run the Gin-Gonic server on port 8080.
	r.Run(":8080")
}

// @Summary Get all todos
// @Description Get all todos
// @Produce json
// @Success 200 {array} Todo
// @Router /todos [get]
func getAllTodos(c *gin.Context) {
	c.JSON(http.StatusOK, todos)
}

// @Summary Get todo by ID
// @Description Get todo by ID
// @Param id path int true "Todo ID"
// @Produce json
// @Success 200 {object} Todo
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /todos/{id} [get]
func getTodoByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, todo := range todos {
		if todo.ID == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}

	c.JSON(http.StatusNotFound, errorDataNotFound)
}

// @Summary Add a new todo
// @Description Add a new todo
// @Accept json
// @Produce json
// @Param todo body Todo true "Todo Object"
// @Success 200 {object} Todo
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /todos [post]
func addTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, errorBadRequest)
		return
	}

	// Generate a new ID (In a real-world scenario, you would use a database with auto-incrementing IDs)
	newTodo.ID = len(todos) + 1

	todos = append(todos, newTodo)
	c.JSON(http.StatusOK, newTodo)
}

// @Summary Edit an existing todo
// @Description Edit an existing todo
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body Todo true "Todo Object"
// @Success 200 {object} Todo
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /todos/{id} [put]
func editTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, todo := range todos {
		if todo.ID == id {
			var updatedTodo Todo
			if err := c.ShouldBindJSON(&updatedTodo); err != nil {
				c.JSON(http.StatusBadRequest, errorBadRequest)
				return
			}

			updatedTodo.ID = id
			todos[i] = updatedTodo

			c.JSON(http.StatusOK, updatedTodo)
			return
		}
	}

	c.JSON(http.StatusNotFound, errorDataNotFound)
}

// @Summary Delete a todo
// @Description Delete a todo
// @Param id path int true "Todo ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /todos/{id} [delete]
func deleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, todo := range todos {
		if todo.ID == id {
			// Remove the todo from the slice
			todos = append(todos[:i], todos[i+1:]...)

			c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, errorDataNotFound)
}
