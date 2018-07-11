package main

import (
	"github.com/labstack/echo"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-echo-vue/handlers"
)


func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}


func migrate(db *sql.DB) {
	sql := `
    CREATE TABLE IF NOT EXISTS tasks(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name VARCHAR NOT NULL
    );
    `

	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}

func main() {
	db := initDB("storage.db")
	migrate(db)
	// Create a new instance of Echo
	e := echo.New()

	e.File("/", "public/index.html") //using to serve a static file that will contain our VueJS client code.
	//e.GET("/tasks", func(c echo.Context) error { return c.JSON(200, "GET Tasks") })
	//e.PUT("/tasks", func(c echo.Context) error { return c.JSON(200, "PUT Tasks") })
	//e.DELETE("/tasks/:id", func(c echo.Context) error { return c.JSON(200, "DELETE Task "+c.Param("id")) })
	// WE ARE CHANGING THE PREVIOUS FUNCTIONS
	// so that they return a function that satisfies the interface.
	// (the following ones will not follow the function signature required by Echo.)
	// This is a trick I used so we can pass around the db instance
	// from handler to handler without having to create a
	// new one each time we want to use the database
	e.GET("/tasks", handlers.GetTasks(db))
	e.PUT("/tasks", handlers.PutTask(db))
	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	// Start as a web server
	//e.Run(standard.New(":8000")) DEPRECATED
	e.Start(":8000")
}