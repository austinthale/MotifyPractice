package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo"
	"net/http"
	"database/sql"
	"github.com/labstack/echo/middleware"
	"encoding/json"
)

/*
var rnd *renderer.Render
*/

/*
// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}
*/
/*
func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./src/public/*.html",
	}

	rnd = renderer.New(opts)
}
*/
/*
// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
*/

type PersonInfo struct {
	Name    string `json:"name"`
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
}

func displayInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, r)
}

//Initializes DB, returns a pointer to sql DB
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

//Migrate the schema
func migrate(db *sql.DB) {
	sql := `
    CREATE TABLE IF NOT EXISTS Employment(
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
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "src/public/resume.html") //using to serve a static file that will contain our VueJS client code.
	e.Static("/static", "src/public")

	//Encoding JSON to write it to the server
	http.HandleFunc("/resumejson", func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(r)
	})

	// Route => handler
	e.GET("/resumejson", displayInfo)



	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

/*func resume(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "resume", nil)
}*/