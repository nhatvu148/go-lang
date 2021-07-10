package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// User is a struct that represents a user in our application
type User struct {
	FullName string `json:"fullName"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Post is a struct that represents a single post
type Post struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author User   `json:"author"`
}

var db *sql.DB
var err error
var posts []Post = []Post{}

func main() {
	host := flag.String("host", "localhost", "Host name")
	user := flag.String("user", "root", "User name")
	password := flag.String("password", "123456789", "Password")
	database := flag.String("database", "jmu", "Database")
	// shipInfoID := flag.Int("shipInfoID", 1, "Ship information ID")
	// startTime := flag.String("startTime", "", "Start time")
	// endTime := flag.String("endTime", "", "End time")
	// outDir := flag.String("outDir", ".", "Output Directory")
	// jpt_root := flag.String("jpt_root", "D:/AKIYAMA/Trunk_Rev56717_ForWeb/bin/Release/x64", "Jupiter Root Directory")
	// outCsvDir := flag.String("outCsvDir", fmt.Sprintf("%s/TechnoStar/00", os.TempDir()), "Output CSV Directory")
	flag.Parse()

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", *user, *password, *host, *database))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	const port string = ":8000"

	router.HandleFunc("/api/jmudt/fatigue", getJmuData(db, "jmu", "fatigue")).Methods("GET")
	router.HandleFunc("/api/jmudt/moment", getJmuData(db, "jmu", "moment")).Methods("GET")
	router.HandleFunc("/api/jmudt/state_calc", getJmuData(db, "jmu", "state_calc")).Methods("GET")
	router.HandleFunc("/api/jmudt/operation", getJmuData(db, "statistics", "operation")).Methods("GET")
	router.HandleFunc("/api/jmudt/gyro", getJmuData(db, "statistics", "gyro")).Methods("GET")
	router.HandleFunc("/api/jmudt/trim01", getJmuData(db, "statistics", "trim01")).Methods("GET")
	router.HandleFunc("/api/jmudt/trim02", getJmuData(db, "statistics", "trim02")).Methods("GET")
	router.HandleFunc("/api/jmudt/trim03", getJmuData(db, "statistics", "trim03")).Methods("GET")
	router.HandleFunc("/api/jmudt/trim04", getJmuData(db, "statistics", "trim04")).Methods("GET")
	router.HandleFunc("/api/jmudt/waves", getJmuData(db, "statistics", "waves")).Methods("GET")
	router.HandleFunc("/api/jmudt/state_statistics", getJmuData(db, "statistics", "state_statistics")).Methods("GET")
	router.HandleFunc("/api/jmudt/statistics_process", getJmuData(db, "statistics", "statistics_process")).Methods("GET")
	router.HandleFunc("/api/jmudt/shipmeasurepoint", getJmuData(db, "statistics", "shipmeasurepoint")).Methods("GET")
	router.HandleFunc("/api/jmudt/rainflow", getJmuData(db, "statistics", "rainflow")).Methods("GET")
	router.HandleFunc("/api/jmudt/compweather", getJmuData(db, "statistics", "compweather")).Methods("GET")
	router.HandleFunc("/api/jmudt/msmwallowable", getJmuData(db, "statistics", "msmwallowable")).Methods("GET")
	router.HandleFunc("/api/jmudt/shipinfo", getJmuData(db, "ship_master", "shipinfo")).Methods("GET")
	router.HandleFunc("/api/jmudt/shipowner", getJmuData(db, "ship_master", "shipowner")).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Up and running...")
	})

	log.Printf("Server listening on port %s\n", port)
	http.ListenAndServe(port, router)
}
