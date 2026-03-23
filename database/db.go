package database
import(
	"log"
	"todo-app-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
// DB is a global variable that holds the database connection. It is of type *gorm.DB, which is a pointer to a GORM database instance. This variable can be accessed throughout the application to perform database operations such as creating, reading, updating, and deleting records. By using a global variable for the database connection, it allows for easy access to the database from different parts of the application without needing to pass the connection around as a parameter.
var DB *gorm.DB
// Connect is a function that establishes a connection to the PostgreSQL database using GORM. It constructs the Data Source Name (DSN) string with the necessary connection parameters such as host, user, password, dbname, port, and sslmode. Then it attempts to open a connection to the database using gorm.Open() with the postgres driver and the DSN. If there is an error during the connection process, it logs the error and terminates the application using log.Fatal(). If the connection is successful, it assigns the database instance to the global variable DB and performs auto-migration for the Todo model to ensure that the corresponding table is created in the database if it does not already exist.
func Connect(){
	// dsn (Data Source Name) is a string that contains the necessary information to connect to a PostgreSQL database. It includes the host, user, password, database name, port, and SSL mode. In this case, the DSN specifies that the database is hosted on localhost, the username is "postgres", the password is "Hanna", the database name is "todo_app", the port is 5432, and SSL mode is disabled. This DSN will be used by GORM to establish a connection to the PostgreSQL database.
	dsn:= "host=localhost user=postgres password=Hanna dbname=todo_app port=5432 sslmode=disable"
	// gorm.Open(postgres.Open(dsn), &gorm.Config{}) attempts to establish a connection to the PostgreSQL database using the GORM library. The postgres.Open(dsn) function creates a new PostgreSQL driver instance with the provided DSN (Data Source Name), which contains the connection details for the database. The &gorm.Config{} argument allows for additional configuration options for GORM, but in this case, it is left empty, meaning that default settings will be used. If the connection is successful, it returns a *gorm.DB instance representing the database connection; if there is an error during the connection process, it returns an error that can be handled accordingly.
	db, err:= gorm.Open(postgres.Open(dsn),&gorm.Config{})
	// if err!=nil checks if there was an error during the database connection process. If there was an error (i.e., if err is not nil), it logs a fatal error message using log.Fatal(), which will print the message "Failed to connect to database:" followed by the specific error details contained in the err variable. This will also terminate the application immediately, preventing any further execution since a successful database connection is crucial for the application's functionality.
	if err!=nil{
		log.Fatal("Failed to connect to database:", err)
	}
	// DB=db assigns the successfully established database connection (represented by the db variable) to the global variable DB. This allows the rest of the application to access the database connection through the DB variable, enabling it to perform various database operations such as creating, reading, updating, and deleting records using GORM's ORM features.
	DB=db
	// log.Println("Database connected successfully") logs a message to the console indicating that the database connection was established successfully. This is useful for debugging and confirming that the application has successfully connected to the database before proceeding with any further operations.
	log.Println("Database connected successfully")
	// DB.AutoMigrate(&models.Todo{}) is a method call that performs auto-migration for the Todo model using GORM. This means that it will automatically create or update the corresponding table in the database based on the structure of the Todo struct defined in the models package. If the table does not exist, it will be created; if it already exists, it will be updated to match any changes in the struct. This is a convenient way to ensure that the database schema is in sync with the application's data models without needing to manually write SQL migration scripts. After the auto-migration is complete, it logs a message indicating that the database migration was successful.
	DB.AutoMigrate(&models.Todo{})
	// log.Println("Database migrated successfully") logs a message to the console indicating that the database migration process was completed successfully. This message serves as a confirmation that the auto-migration for the Todo model was executed without any issues, and the corresponding table in the database is now up to date with the structure defined in the Todo struct.
	log.Println("Database migrated successfully")
}