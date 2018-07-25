package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"github.com/labstack/echo/middleware"
	"fmt"
)

type PersonInfo struct {
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Phone   string `db:"phone" json:"phone"`
	Email   string `db:"email" json:"email"`
}

type Education struct {
	ID			 int		`db:"id" json:"id"`
	School       string   	`db:"school" json:"school"`
	DateAttended string   	`db:"date_attended" json:"date_attended"`
	Notes        []string 	`db:"notes" json:"notes"`
	ResumeID	 int 		`db:"resume_id" json:"resume_id"`
}

type Employment struct {
	ID			 int		`db:"id" json:"id"`
	Company      string   	`db:"company" json:"company"`
	DateAttended string   	`db:"date_attended" json:"date_attended"`
	Position     string   	`db:"position" json:"position"`
	Notes        []string 	`db:"notes" json:"notes"`
	ResumeID	 int 		`db:"resume_id" json:"resume_id"`
}

type Volunteer struct {
	ID			 int		`db:"id" json:"id"`
	Company      string   	`db:"company" json:"company"`
	DateAttended string   	`db:"date_attended" json:"date_attended"`
	Position     string   	`db:"position" json:"position"`
	Notes        []string 	`db:"notes" json:"notes"`
	ResumeID	 int 		`db:"resume_id" json:"resume_id"`
}

type Resume struct {
	PersonInfo  PersonInfo   `db:"person_info" json:"person_info"`
	Educations  []Education  `db:"educations" json:"educations"`
	Employments []Employment `db:"employments" json:"employments"`
	Volunteers  []Volunteer  `db:"volunteers" json:"volunteers"`
}

var r = Resume{
	PersonInfo: PersonInfo{},
	Educations: []Education{},
	Employments: []Employment{},
	Volunteers: []Volunteer{},
}

//Initializes DB, returns a pointer to sql DB
func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath) // Open connection to DB
	checkErr(err)		// Here we check for any db errors then exit

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}

//Migrate the schema
func migrate(db *sql.DB) {
	// SQL statement to create a task table, with no records in it.
	sql := `
    CREATE TABLE IF NOT EXISTS resume(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        full_name VARCHAR,
		address VARCHAR,
		phone VARCHAR,
		email VARCHAR
    );
	
	CREATE TABLE IF NOT EXISTS education(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		school TEXT,
	  	date_attended TEXT,
		resume_id INTEGER,
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	
	CREATE TABLE IF NOT EXISTS employment(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		company TEXT,
		position TEXT,
	  	date_attended TEXT,
		resume_id INTEGER,
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	
	CREATE TABLE IF NOT EXISTS volunteer(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		company TEXT,
		position TEXT,
	  	date_attended TEXT,
		resume_id INTEGER,
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	
	CREATE TABLE IF NOT EXISTS edu_note(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		education_id INTEGER,
		resume_id INTEGER,
		detail TEXT,
		FOREIGN KEY(education_id) REFERENCES education(id),
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);

	CREATE TABLE IF NOT EXISTS emp_note(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		employment_id INTEGER,
		resume_id INTEGER,
		detail TEXT,
		FOREIGN KEY(employment_id) REFERENCES employment(id),
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);

	CREATE TABLE IF NOT EXISTS vol_note(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		volunteer_id INTEGER,
		resume_id INTEGER,
		detail TEXT,
		FOREIGN KEY(volunteer_id) REFERENCES volunteer(id),
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	`

	// Execute the SQL statement
	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	checkErr(err)
}

func GetData(db *sql.DB) echo.HandlerFunc {
	rows, err := db.Query("SELECT id, full_name, address, phone, email from resume WHERE id = 1")
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var address string
		var phone string
		var email string
		err = rows.Scan(&id, &name, &address, &phone, &email)
		checkErr(err)
		r.PersonInfo.Name = name
		r.PersonInfo.Address = address
		r.PersonInfo.Phone = phone
		r.PersonInfo.Email = email
	}
	r.Educations = getEduData(db)
	r.Employments = getEmpData(db)
	r.Volunteers = getVolData(db)

	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, r)
	}
}

func PutData(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		res := new(Resume)
		if err = c.Bind(res); err != nil {
			return
		}
		r = *res
		// R now reflects the input that we used BIND to get, we can update the database
		sql := "INSERT OR REPLACE INTO resume(id, full_name, address, phone, email) VALUES(?, ?, ?, ?, ?)"
		// Create a prepared SQL statement
		stmt, err := db.Prepare(sql)
		checkErr(err)
		// Replace the '?' in our prepared statement with values in our resume object
		_, err2 := stmt.Exec(1, r.PersonInfo.Name, r.PersonInfo.Address, r.PersonInfo.Phone, r.PersonInfo.Email)
		checkErr(err2)

		setEduData(db)
		setEmpData(db)
		setVolData(db)

		return c.JSON(http.StatusOK, res) //changed r to res
	}
}
func getEduData(db *sql.DB) []Education {
	eduArray := []Education{}
	// Find all educations belonging to the first resume
	eduRows, err := db.Query("SELECT id, school, date_attended, resume_id FROM education WHERE resume_id = 1")
	//var count = 1
	// For each education, get data and create new Education object, then add it to our eduArray
	for eduRows.Next() {
		var id int
		var school string
		var dateAttended string
		var resumeID int
		//var notes []string
		err = eduRows.Scan(&id, &school, &dateAttended, &resumeID)
		checkErr(err)

		//construct new Education object based on data we just read in, and add it to our array of Education objects
		edu := Education{
			ID: id,
			School: school,
			DateAttended: dateAttended,
			ResumeID: resumeID,
		}

		eduArray = append(eduArray, edu)

	}
	// Get NOTES for each School section in Education
	for idx, edu := range eduArray {
		var notes []string
		noteRows, err := db.Query("SELECT id, education_id, detail, resume_id FROM edu_note WHERE education_id = " + strconv.Itoa(edu.ID) + " AND resume_id = " + strconv.Itoa(edu.ResumeID))
		checkErr(err)
		for noteRows.Next() {
			var id int
			var secID int
			var txt string
			var resID int
			err = noteRows.Scan(&id, &secID, &txt, &resID)

			checkErr(err)
			notes = append(notes, txt)
		}
		eduArray[idx].Notes = notes
	}
	return eduArray
}
func getEmpData(db *sql.DB) []Employment {
	empArray := []Employment{}
	// Find all employments belonging to the first resume
	empRows, err := db.Query("SELECT id, company, date_attended, resume_id FROM employment WHERE resume_id = 1")
	// For each employment, get data and create new Employment object, then add it to our empArray
	for empRows.Next() {
		var id int
		var company string
		var dateAttended string
		var resumeID int
		err = empRows.Scan(&id, &company, &dateAttended, &resumeID)
		checkErr(err)

		//construct new Employment object based on data we just read in, and add it to our array of Employment objects
		job := Employment{
			ID: id,
			Company: company,
			DateAttended: dateAttended,
			ResumeID: resumeID,
		}
		empArray = append(empArray, job)
	}
	// Get NOTES for each Company section in Employment
	for idx, job := range empArray {
		var notes []string
		noteRows, err := db.Query("SELECT id, employment_id, detail, resume_id FROM emp_note WHERE employment_id = " + strconv.Itoa(job.ID) + " AND resume_id = " + strconv.Itoa(job.ResumeID))
		checkErr(err)
		for noteRows.Next() {
			var id int
			var secID int
			var txt string
			var resID int
			err = noteRows.Scan(&id, &secID, &txt, &resID)
			checkErr(err)
			notes = append(notes, txt)
		}
		empArray[idx].Notes = notes
	}
	return empArray
}
func getVolData(db *sql.DB) []Volunteer {
	volArray := []Volunteer{}
	// Find all volunteers belonging to the first resume
	volRows, err := db.Query("SELECT id, company, date_attended, resume_id FROM volunteer WHERE resume_id = 1")
	// For each volunteer, get data and create new Volunteer object, then add it to our volArray
	for volRows.Next() {
		var id int
		var company string
		var dateAttended string
		var resumeID int
		err = volRows.Scan(&id, &company, &dateAttended, &resumeID)
		checkErr(err)

		//construct new Volunteer object based on data we just read in, and add it to our array of Volunteer objects
		job := Volunteer{
			ID: id,
			Company: company,
			DateAttended: dateAttended,
			ResumeID: resumeID,
		}
		volArray = append(volArray, job)
	}
	// Get NOTES for each Company section in Volunteer
	for idx, job := range volArray {
		var notes []string
		noteRows, err := db.Query("SELECT id, volunteer_id, detail, resume_id FROM vol_note WHERE volunteer_id = " + strconv.Itoa(job.ID) + " AND resume_id = " + strconv.Itoa(job.ResumeID))
		checkErr(err)
		for noteRows.Next() {
			var id int
			var secID int
			var txt string
			var resID int
			err = noteRows.Scan(&id, &secID, &txt, &resID)
			checkErr(err)
			notes = append(notes, txt)
		}
		volArray[idx].Notes = notes
	}
	return volArray
}

func setEduData(db *sql.DB) {
	//var thisID int
	db.Exec("DELETE FROM education WHERE resume_id = 1")
	for i, edu := range r.Educations {
		sql := "INSERT INTO education(id, school, date_attended, resume_id) VALUES(?, ?, ?, ?)"
		stmt, err := db.Prepare(sql)
		checkErr(err)
		_, err = stmt.Exec(i+1, edu.School, edu.DateAttended, 1)

		//delete Note record for current education
		db.Exec("DELETE FROM edu_note WHERE education_id = " + strconv.Itoa(i + 1))
		for _, note := range edu.Notes {
			sql := "INSERT INTO edu_note(education_id, resume_id, detail) VALUES(?, ?, ?)"
			stmt, err := db.Prepare(sql)
			checkErr(err)
			_, err = stmt.Exec(i+1, 1, note)
		}
	}
}
func setEmpData(db *sql.DB) {
	db.Exec("DELETE FROM employment WHERE resume_id = 1")
	for i, job := range r.Employments {
		sql := "INSERT INTO employment(id, company, date_attended, resume_id) VALUES(?, ?, ?, ?)"
		stmt, err := db.Prepare(sql)
		checkErr(err)
		_, err = stmt.Exec(i+1, job.Company, job.DateAttended, 1)

		//delete Note record for current education
		db.Exec("DELETE FROM emp_note WHERE employment_id = " + strconv.Itoa(i + 1))
		for _, note := range job.Notes {
			sql := "INSERT INTO emp_note(employment_id, resume_id, detail) VALUES(?, ?, ?)"
			stmt, err := db.Prepare(sql)
			checkErr(err)
			_, err = stmt.Exec(i+1, 1, note)
		}
	}
}
func setVolData(db *sql.DB) {
	db.Exec("DELETE FROM volunteer WHERE resume_id = 1")
	for i, job := range r.Volunteers {
		sql := "INSERT INTO volunteer(id, company, date_attended, resume_id) VALUES(?, ?, ?, ?)"
		stmt, err := db.Prepare(sql)
		checkErr(err)
		_, err = stmt.Exec(i+1, job.Company, job.DateAttended, 1)

		//delete Note record for current education
		db.Exec("DELETE FROM vol_note WHERE volunteer_id = " + strconv.Itoa(i + 1))
		for _, note := range job.Notes {
			sql := "INSERT INTO vol_note(volunteer_id, resume_id, detail) VALUES(?, ?, ?)"
			stmt, err := db.Prepare(sql)
			checkErr(err)
			_, err = stmt.Exec(i+1, 1, note)
		}
	}
}
func main() {
	// Initialize the DB
	db := initDB("resume.db")
	// close database connection before exiting program.
	defer db.Close()
	migrate(db)		// Schema migration

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "src/public/resume.html") // using to serve a static file that will contain our VueJS client code.
	e.Static("/static", "src/public")	// using to serve all files contained in public folder, and must be accessed
										// through the /static folder ("localhost:1323/static/resume.css")

	// Route => handler
	//e.GET("/resumedb", handlers.GetData(db))
	//e.PUT("/resumedb", handlers.PutData(db))

	e.GET("/resumejson", GetData(db))
	e.POST("/resumejson", PutData(db))


	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func checkErr(err error, args ...string) {
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err, args)
		return
	}
}