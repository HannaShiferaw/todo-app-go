package main	
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-app-go/models"
	"todo-app-go/database"
)
func main(){
	// to create a server
	r:=gin.Default()
	// to connect to database
database.Connect()
	// defining get route at /
	r.GET("/", func(c *gin.Context){
		// sends json request
		c.JSON(200,gin.H{
			"message":"API is running",
		})
	})
	// defining post route at /todos
r.POST("/todos", func(c *gin.Context){
	// ShouldBindJSON is a method provided by the Gin framework that binds the incoming JSON payload from the request body to the specified struct. In this case, it attempts to bind the JSON data to the todo variable of type models.Todo. If the binding is successful, the todo variable will be populated with the data from the JSON request. If there is an error during binding (e.g., if the JSON structure does not match the Todo struct), it will return an error, which is handled in the subsequent lines of code. 
	var todo models.Todo
	// c.ShouldBindJSON(&todo) attempts to bind the incoming JSON payload from the request body to the todo variable of type models.Todo. If there is an error during this binding process (e.g., if the JSON structure does not match the Todo struct), it will return an error. The code checks for this error, and if it exists, it responds with a 400 Bad Request status and a JSON object containing the error message. This ensures that the client is informed about any issues with the JSON payload they sent in the request.
	if err:=c.ShouldBindJSON(&todo);
	// if
	 err!=nil{
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) sends a JSON response with a 400 Bad Request status code and a JSON object containing the error message. The error message is obtained from the err variable, which contains details about what went wrong during the JSON binding process. This response informs the client that their request was invalid and provides information about the specific error that occurred.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return is used to exit the current function early. In this context, it is used to stop further execution of the code in the handler function after sending the error response. This ensures that if there is an error during JSON binding, the function will not continue to execute any subsequent code that may rely on a successfully bound todo variable, preventing potential issues or unintended behavior.
		return
	}
	// database.DB.Create(&todo) is a method call that creates a new record in the database using the GORM library. The &todo argument is a pointer to the todo variable, which contains the data that should be inserted into the database. When this method is called, GORM will generate an SQL INSERT statement based on the fields of the Todo struct and execute it against the connected database. If the operation is successful, a new record will be created in the corresponding table with the values from the todo variable. The result of this operation, including any potential errors, is stored in the result variable for further handling.
	result:= database.DB.Create(&todo)
	// if result.Error!=nil checks if there was an error during the database creation process. If there was an error (i.e., if result.Error is not nil), it sends a JSON response with a 500 Internal Server Error status code and a JSON object containing the error message. The error message is obtained from result.Error.Error(), which provides details about what went wrong during the database operation. This response informs the client that there was an issue with the server while trying to create the new record in the database.
	if result.Error!=nil{
		// c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()}) sends a JSON response with a 500 Internal Server Error status code and a JSON object containing the error message. The error message is obtained from the result.Error.Error() method, which provides details about what went wrong during the database operation. This response informs the client that there was an issue with the server while trying to create the new record in the database.
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
// c.JSON(http.StatusOK, todo) sends a JSON response with a 200 OK status code and the todo object as the response body. This indicates that the request was successful and returns the created todo item in the response. The todo variable contains the data that was just inserted into the database, including any auto-generated fields such as ID or timestamps, which can be useful for the client to have after creating a new record.
	c.JSON(http.StatusOK, todo)
})
r.GET("/todos", func(c *gin.Context){
	var todos[]models.Todo
	// to fetch all todos from database
	result:= database.DB.Find(&todos)
	if result.Error!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":result.Error.Error()})
			return
		}
		// return all todos as json
	c.JSON(http.StatusOK, todos)
	})

r.GET("/todos/:id", func(c *gin.Context){
	var todo models.Todo
	id:= c.Param("id")
	result :=database.DB.First(&todo,id)
	if result.Error!=nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
})

r.Run()
}