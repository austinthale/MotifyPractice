package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strconv"
)

type PersonInfo struct {
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Phone   string `db:"phone" json:"phone"`
	Email   string `db:"email" json:"email"`
}

type Education struct {
	id			 int		`db:"id" json:"id"`
	School       string   	`db:"school" json:"school"`
	DateAttended string   	`db:"date_attended" json:"date_attended"`
	Notes        []string 	`db:"notes" json:"notes"`
}

type Employment struct {
	id			 int		`db:"id" json:"id"`
	Company      string   	`db:"company" json:"company"`
	DateAttended string   	`db:"date_attended" json:"date_attended"`
	Position     string   	`db:"position" json:"position"`
	Notes        []string 	`db:"notes" json:"notes"`
}

type Volunteer struct {
	id			 int		`db:"id" json:"id"`
	Company      string   	`db:"company" json:"company"`
	DateAttended string   	`db:"date_attended" json:"date_attended"`
	Position     string   	`db:"position" json:"position"`
	Notes        []string 	`db:"notes" json:"notes"`
}

type Resume struct {
	PersonInfo  PersonInfo   `db:"person_info" json:"person_info"`
	Educations  []Education  `db:"educations" json:"educations"`
	Employments []Employment `db:"employments" json:"employments"`
	Volunteers  []Volunteer  `db:"volunteers" json:"volunteers"`
}

var r = Resume{
	PersonInfo: PersonInfo{ /*
		Name:    "Austin Hale",
		Address: "757 S 320 W, Provo, UT 84601",
		Phone:   "+1-559-346-7123",
		Email:   "austin.t.hale89@gmail.com",*/
	},
	Educations: []Education{ /*
		{
			School:       "Utah Valley University",
			DateAttended: "Aug 2015 - Present",
			Notes: []string{
				"Cumulative GPA of 3.8, expected to graduate in Spring 2019.",
				"Experienced in object-oriented and game development, design, testing, and debugging in teams and independently."},
		},
		{
			School:       "BYU-Idaho",
			DateAttended: "Apr 2011 – Dec 2012",
			Notes: []string{
				"Completed General Education classes and Intro to Programming classes"},
		},*/
	},
	Employments: []Employment{ /*
		{
			Company:      "Utah Valley University CS Tutor Lab",
			DateAttended: "Jan 2018 - Present",
			Position:     "Computer Science Tutor and Grader",
			Notes: []string{
				"Coached students 1-on-1, teaching them to think critically for problem solving, debugging, and effective programming.",
				"Facilitated numerous programming courses as a grader and class assistant, providing additional guidance and simplifying programming concepts (e.g. CPU Scheduling simulation, AVL Trees, Advanced Algorithms, etc.)",
			},
		},
		{
			Company:      "Frontier Communications",
			DateAttended: "Feb 2016 – Aug 2016",
			Position:     "Customer Service Representative",
			Notes: []string{
				"Averaged a 94 NPS by resolving customer concerns through explaining policies with professionalism and simplicity.",
			},
		}, */
	},
	Volunteers: []Volunteer{ /*
		{
			Company:      "Xing Zhou Bilingual School; Guangdong, China",
			DateAttended: "Mar 2015 - Jul 2015",
			Position:     "English Teacher",
			Notes: []string{
				"Adapted to Chinese culture while teaching English to 12 classes of 30+ students from 1st-8th grade",
			},
		},
		{
			Company:      "Church of Jesus Christ of Latter-day Saints; Salta, Argentina",
			DateAttended: "Feb 2009 - Feb 2011",
			Position:     "Full-time Representative",
			Notes: []string{
				"Organized and taught 10-60 minute lessons to individuals and families while mastering Spanish",
				"Demonstrated leadership by planning service and community projects as Branch President",
			},
		}, */
	},
}

type resume struct {
	id 			int			`db:"id"`
	full_name 	string		`db:"name" json:"name"`
	address 	string		`db:"address"`
	phone		string		`db:"phone"`
	email		string		`db:"email"`
}
type education struct {
	id			 int		`db:"id"`
	school       string		`db:"school"`
	DateAttended string		`db:"date_attended"`
	Notes        []string	`db:"notes"`
}

type employment struct {
	id			 int		`db:"id"`
	Company      string		`db:"company"`
	DateAttended string		`db:"date_attended"`
	Position     string		`db:"position"`
	Notes        []string
}

type volunteer struct {
	id 			 int		`db:"id"`
	Company      string		`db:"company"`
	DateAttended string		`db:"date_attended"`
	Position     string		`db:"position"`
	Notes        []string	`db:"notes"`
}

type note struct {
	id 			 int		`db:"id"`
	section_id 	 int		`db:"section_id"`
	detail 		 string		`db:"detail"`
	note_type    int		`db:"note_type"`
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
	DELETE FROM resume;
	CREATE TABLE IF NOT EXISTS education(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		school TEXT,
	  	date_attended TEXT,
		resume_id INTEGER,
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	DELETE FROM education;
	CREATE TABLE IF NOT EXISTS employment(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		company TEXT,
		position TEXT,
	  	date_attended TEXT,
		resume_id INTEGER,
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	DELETE FROM employment;
	CREATE TABLE IF NOT EXISTS volunteer(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		company TEXT,
		position TEXT,
	  	date_attended TEXT,
		resume_id INTEGER,
		FOREIGN KEY(resume_id) REFERENCES resume(id)
	);
	DELETE FROM volunteer;
	CREATE TABLE IF NOT EXISTS note(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		section_id INTEGER,
		detail TEXT,
		type TEXT
	);
	DELETE FROM note;
    `
	// Execute the SQL statement
	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	checkErr(err)
}

func insertRecords(db *sql.DB) {
	// Begin transaction
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("INSERT INTO resume(id, full_name, address, phone, email) VALUES(?, ?, ?, ?, ?)")
	checkErr(err)
	//defer stmt.Close()

	// Create empty slice of resume struct pointers.
	resumes := []*resume{}
	resumes = append(resumes, &resume{id: 1, full_name: "Austin T. Hale", address: "123 Main st.", phone: "555-555-5555", email: "test@gmail.com"})

	_, err = stmt.Exec(resumes[0].id, resumes[0].full_name, resumes[0].address, resumes[0].phone, resumes[0].email)
	checkErr(err)

	tx.Commit()
	_, err = db.Exec("INSERT INTO education(id, school, date_attended, resume_id) values(1, 'UVU', 'Aug 2015 - Present',  1),(2, 'BYU-Idaho', 'Apr 2011 – Dec 2012', 1)")
	_, err = db.Exec("INSERT INTO employment(id, company, position, date_attended, resume_id) values(1, 'UVU CS Tutor Lab', 'Tutor', 'Jan 2018 - Present',  1),(2, 'Frontier Communications', 'Customer Service Rep','Feb 2016 – Aug 2016', 1)")
	_, err = db.Exec("INSERT INTO volunteer(id, company, position, date_attended, resume_id) values(1, 'Xing Zhou Bilingual School; Guangdong, China', 'English Teacher', 'Mar 2015 - Jul 2015',  1),(2, 'Church of Jesus Christ of Latter-day Saints', 'Missionary','Feb 2009 - Feb 2011', 1)")

	_, err = db.Exec("INSERT INTO note(id, section_id, detail, type) values(1, 1, 'Cumulative GPA of 3.8, expected to graduate in Spring 2019.', 'Education'),(2, 1, 'Experienced in object-oriented and game development, design, testing, and debugging in teams and independently.', 'Education')")
	_, err = db.Exec("INSERT INTO note(id, section_id, detail, type) values(3, 2, 'Completed General Education classes and Intro to Programming classes', 'Education'),(4, 1, 'Coached students 1-on-1, teaching them to think critically for problem solving, debugging, and effective programming.', 'Employment')")
	_, err = db.Exec("INSERT INTO note(id, section_id, detail, type) values(5, 1, 'Facilitated numerous programming courses as a grader and class assistant, providing additional guidance and simplifying programming concepts (e.g. CPU Scheduling simulation, AVL Trees, Advanced Algorithms, etc.)', 'Employment'),(6, 2, 'Averaged a 94 NPS by resolving customer concerns through explaining policies with professionalism and simplicity.', 'Employment')")
	_, err = db.Exec("INSERT INTO note(id, section_id, detail, type) values(7, 1, 'Adapted to Chinese culture while teaching English to 12 classes of 30+ students from 1st-8th grade', 'Volunteer'),(8, 2, 'Organized and taught 10-60 minute lessons to individuals and families while mastering Spanish', 'Volunteer')")
	_, err = db.Exec("INSERT INTO note(id, section_id, detail, type) values(9, 2, 'Demonstrated leadership by planning service and community projects as Branch President', 'Volunteer')")
}

func readRecords(db *sql.DB) {
	rows, err := db.Query("SELECT * from resume")
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
		fmt.Println(id, name, address, phone, email)
	}
	err = rows.Err()
	checkErr(err)
}


func displayInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, r)
}

func saveInfo(c echo.Context) (err error) {
	res := new(Resume)
	if err = c.Bind(res); err != nil {
		return
	}
	r = *res
	return c.JSON(http.StatusOK, res) //changed r to res
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
	//log.Println(eduArray)
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
		//eduDetails, err := db.Query("SELECT id, section_id, detail, type FROM notes WHERE type = 'Education' AND section_id = ?")

		//notes = append(notes, )

		//construct new Education object based on data we just read in, and add it to our array of Education objects
		edu := Education{
			id: id,
			School: school,
			DateAttended: dateAttended,
		}

		eduArray = append(eduArray, edu)

	}
	// Get NOTES for each School section in Education
	for idx, edu := range eduArray {
		var notes []string
		noteRows, err := db.Query("SELECT id, section_id, detail, type FROM note WHERE section_id = " + strconv.Itoa(edu.id) + " AND type = 'Education'")
		checkErr(err)
		for noteRows.Next() {
			var id int
			var secID int
			var txt string
			var noteType string
			err = noteRows.Scan(&id, &secID, &txt, &noteType)
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

		//notes = append(notes, )

		//construct new Employment object based on data we just read in, and add it to our array of Employment objects
		job := Employment{
			id: id,
			Company: company,
			DateAttended: dateAttended,
		}
		empArray = append(empArray, job)
	}
	// Get NOTES for each Company section in Employment
	for idx, job := range empArray {
		var notes []string
		noteRows, err := db.Query("SELECT id, section_id, detail, type FROM note WHERE section_id = " + strconv.Itoa(job.id) + " AND type = 'Employment'")
		checkErr(err)
		for noteRows.Next() {
			var id int
			var secID int
			var txt string
			var noteType string
			err = noteRows.Scan(&id, &secID, &txt, &noteType)
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

		//notes = append(notes, )

		//construct new Volunteer object based on data we just read in, and add it to our array of Volunteer objects
		job := Volunteer{
			id: id,
			Company: company,
			DateAttended: dateAttended,
		}
		volArray = append(volArray, job)
	}
	// Get NOTES for each Company section in Volunteer
	for idx, job := range volArray {
		var notes []string
		noteRows, err := db.Query("SELECT id, section_id, detail, type FROM note WHERE section_id = " + strconv.Itoa(job.id) + " AND type = 'Employment'")
		checkErr(err)
		for noteRows.Next() {
			var id int
			var secID int
			var txt string
			var noteType string
			err = noteRows.Scan(&id, &secID, &txt, &noteType)
			checkErr(err)
			notes = append(notes, txt)
		}
		volArray[idx].Notes = notes
	}
	return volArray
}
func main() {
	// Initialize the DB
	db := initDB("resume.db")
	// close database connection before exiting program.
	defer db.Close()
	migrate(db)		// Schema migration

	insertRecords(db)
	readRecords(db)

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
		//log.Fatal(err)
		return
	}
}