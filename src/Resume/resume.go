package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type PersonInfo struct {
	Name    string
	Address string
	Phone   string
	Email   string
}

type Education struct {
	School			string
	DateAttended 	string
	Notes 			[]string
}

type Employment struct {
	Company 		string
	DateAttended	string
	Position 		string
	Notes 			[]string
}

type Volunteer struct {
	Company 		string
	DateAttended	string
	Position 		string
	Notes 			[]string
}

func displayInfo(c echo.Context) error {
	p := PersonInfo{
		Name:    "Austin Hale",
		Address: "757 S 320 W,\nProvo, UT 84601",
		Phone:   "+1-559-346-7123",
		Email:   "austin.t.hale89@gmail.com",
	}
	//if err := c.Bind(p); err != nil {
	//	return err
	//}

	//Get Name, Address, Phone, and Email
	name := p.Name
	address := p.Address
	phone := p.Phone
	email := p.Email
	education := getEducation()
	career := getJobs()
	volunteer := getVolunteer()
	return c.String(http.StatusOK, "" + name + "\n" + address + "\n" + phone + "\n" + email + "\n" + education + "\n" + career + "\n" + volunteer)
}

func getEducation() string {
	schoolStr := "\n\nEducation:\n"
	s1 := Education{
		School: 		"Utah Valley University",
		DateAttended: 	"Aug 2015 - Present",
		Notes:			nil,
	}
	s1.Notes = append(s1.Notes, "Cumulative GPA of 3.8, expected to graduate in Spring 2019.")
	s1.Notes = append(s1.Notes, "Experienced in object-oriented and game development, design, testing, and debugging in teams and independently.")
	schoolStr += s1.School + "\n" + s1.DateAttended + "\n"
	for i := 0; i < len(s1.Notes); i++  {
		schoolStr += s1.Notes[i] + "\n"
	}
	s2 := Education{
		School: 		"BYU-Idaho",
		DateAttended: 	"Apr 2011 – Dec 2012",
		Notes:			nil,
	}
	s2.Notes = append(s2.Notes, "Completed General Education classes and Intro to Programming classes")
	schoolStr += s2.School + "\n" + s2.DateAttended + "\n"
	for i := 0; i < len(s2.Notes); i++  {
		schoolStr += s2.Notes[i] + "\n"
	}
	return schoolStr
}

func getJobs() string {
	jobStr := "\n\nEmployment:\n"
	e1 := Employment{
		Company: 		"Utah Valley University CS Tutor Lab",
		DateAttended: 	"Jan 2018 - Present",
		Position:		"Computer Science Tutor and Grader",
		Notes:			nil,
	}
	e1.Notes = append(e1.Notes, "Coached students 1-on-1, teaching them to think critically for problem solving, debugging, and effective programming.")
	e1.Notes = append(e1.Notes, "Facilitated numerous programming courses as a grader and class assistant, providing additional guidance and simplifying programming concepts (e.g. CPU Scheduling simulation, AVL Trees, Advanced Algorithms, etc.)")
	jobStr += e1.Company + "\n" + e1.DateAttended + "\n"
	for i := 0; i < len(e1.Notes); i++  {
		jobStr += e1.Notes[i] + "\n"
	}
	e2 := Employment{
		Company: 		"Frontier Communications",
		DateAttended: 	"Feb 2016 – Aug 2016",
		Position:		"Customer Service Representative",
		Notes:			nil,
	}
	e2.Notes = append(e2.Notes, "Averaged a 94 NPS by resolving customer concerns through explaining policies with professionalism and simplicity.")
	jobStr += e2.Company + "\n" + e2.DateAttended + "\n"
	for i := 0; i < len(e2.Notes); i++  {
		jobStr += e2.Notes[i] + "\n"
	}
	return jobStr
}

func getVolunteer() string {
	volunteerStr := "\nVolunteer Experience:\n"
	v1 := Volunteer{
		Company: 		"Xing Zhou Bilingual School; Guangdong, China",
		DateAttended: 	"Mar 2015 - Jul 2015",
		Position:		"English Teacher",
		Notes:			nil,
	}
	v1.Notes = append(v1.Notes, "Adapted to Chinese culture while teaching English to 12 classes of 30+ students from 1st-8th grade")
	volunteerStr += v1.Company + "\n" + v1.DateAttended + "\n"
	for i := 0; i < len(v1.Notes); i++  {
		volunteerStr += v1.Notes[i] + "\n"
	}
	v2 := Volunteer{
		Company: 		"Church of Jesus Christ of Latter-day Saints; Salta, Argentina",
		DateAttended: 	"Feb 2009 - Feb 2011",
		Position:		"Full-time Representative",
		Notes:			nil,
	}
	v2.Notes = append(v2.Notes, "Organized and taught 10-60 minute lessons to individuals and families while mastering Spanish")
	v2.Notes = append(v2.Notes, "Demonstrated leadership by planning service and community projects as Branch President")
	volunteerStr += v2.Company + "\n" + v2.DateAttended + "\n"
	for i := 0; i < len(v2.Notes); i++  {
		volunteerStr += v2.Notes[i] + "\n"
	}
	return volunteerStr
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", displayInfo)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}