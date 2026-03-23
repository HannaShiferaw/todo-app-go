package models

import (
		"gorm.io/gorm"
)
type Todo struct{
// gorm.Model provides fields ID, CreatedAt, UpdatedAt, DeletedAt to the struct and also implements the gorm.Model interface, which allows it to be used with GORM's ORM features. By embedding gorm.Model, the Todo struct inherits these fields and can be easily managed by GORM for database operations such as creating, updating, and deleting records. This means that when you use GORM to interact with the database, it will automatically handle the ID generation, timestamping, and soft deletion for the Todo struct without requiring you to manually define those fields or logic.
	gorm.Model
	// json:"title" specifies that when the struct is serialized to JSON, the Title field will be represented as "title" in the JSON output. This is useful for ensuring that the JSON keys are in a specific format, such as lowercase, which is common in JSON APIs. When you marshal a Todo struct to JSON, the Title field will appear as "title" in the resulting JSON object.
	Title string `json:"title"`
	// json:"description" specifies that when the struct is serialized to JSON, the Description field will be represented as "description" in the JSON output. This allows for consistent naming conventions in the JSON representation of the struct, making it easier to work with in APIs and other contexts where JSON is used for data exchange. When you marshal a Todo struct to JSON, the Description field will appear as "description" in the resulting JSON object.
	Description string `json:"description"`
	// json:"completed" specifies that when the struct is serialized to JSON, the Completed field will be represented as "completed" in the JSON output. This allows for consistent naming conventions in the JSON representation of the struct, making it easier to work with in APIs and other contexts where JSON is used for data exchange. When you marshal a Todo struct to JSON, the Completed field will appear as "completed" in the resulting JSON object.
	Completed bool `json:"completed"`
}
