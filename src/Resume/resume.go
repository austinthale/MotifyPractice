package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)
/*
type PersonInfo struct {
	Name    string `db:"name" json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type Education struct {
	School       string   `json:"school"`
	DateAttended string   `json:"date_attended"`
	Notes        []string `json:"notes"`
}

type Employment struct {
	Company      string   `json:"company"`
	DateAttended string   `json:"date_attended"`
	Position     string   `json:"position"`
	Notes        []string `json:"notes"`
}

type Volunteer struct {
	Company      string   `json:"company"`
	DateAttended string   `json:"date_attended"`
	Position     string   `json:"position"`
	Notes        []string `json:"notes"`
}

type Resume struct {
	PersonInfo  PersonInfo   `json:"person_info"`
	Educations  []Education  `json:"educations"`
	Employments []Employment `json:"employments"`
	Volunteers  []Volunteer  `json:"volunteers"`
}

var r = Resume{
	PersonInfo: PersonInfo{
		Name:    "Austin Hale",
		Address: "757 S 320 W, Provo, UT 84601",
		Phone:   "+1-559-346-7123",
		Email:   "austin.t.hale89@gmail.com",
	},
	Educations: []Education{
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
		},
	},
	Employments: []Employment{
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
		},
	},
	Volunteers: []Volunteer{
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
		},
	},
}*/

type resume struct {
	id 		int
	full_name 	string
	address string
	phone	string
	email	string
}
type education struct {
	id			 int
	school       string
	DateAttended string
	Notes        []string
}

type employment struct {
	id			 int
	Company      string
	DateAttended string
	Position     string
	Notes        []string
}

type volunteer struct {
	id 			 int
	Company      string
	DateAttended string
	Position     string
	Notes        []string
}

type Note struct {
	id 			 int
	section_id 	 int
	detail 		 string
	note_type   		 int
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


/*func displayInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, r)
}

func saveInfo(c echo.Context) (err error) {
	res := new(Resume)
	if err = c.Bind(res); err != nil {
		return
	}
	r = *res
	return c.JSON(http.StatusOK, res) //changed r to res
}*/

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
	//e.GET("/resumejson", displayInfo)

	//e.POST("/resumejson", saveInfo)


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