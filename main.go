package main

import ( "fmt" ; "log" ; "encoding/json" ; "net/http" ; "strconv" ; "github.com/gorilla/mux" ; "database/sql"; _ "github.com/go-sql-driver/mysql"  )

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8000"
	DRIVER_NAME = "mysql"
	DATASOURCE_NAME = "root:G22mEyct140522FU@/moviesdb"
)

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc http.HandlerFunc
}

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Name string `json:"name"`
}

var routes []Route = []Route{
	{"getDB","GET","/", getCurrentDb},
	{"getMovieRec","GET","/movies", readRecord},
	{"createMovieRec","POST","/movie", createRecord},
	{"updateMovieRec","PUT","/movie/{id}", updateRecord},
	{"deleteMovieRec","DELETE","/movie/{id}", deleteRecord},
}

var movies []Movie

var db *sql.DB // declare a global private db instance 
var connectionError error

func init(){
	// connect to database
	db, connectionError =sql.Open(DRIVER_NAME,DATASOURCE_NAME)
	if connectionError != nil {
        log.Fatal("error connecting to database :: ", connectionError)
    }
}

func getCurrentDb(w http.ResponseWriter, r *http.Request){
	rows,err := db.Query("SELECT DATABASE() as db") // perform select query on db to get current db name
	if err != nil{
		log.Print("error executing query :: ", err)
		return
	}
	var db string
	// iterate over records
	for rows.Next() {
		rows.Scan(&db) // copy db values variable db
	}
	fmt.Fprintf(w, "Current Database is :: %s", db) // write to http response stream
}

func createRecord(w http.ResponseWriter, r *http.Request){
	vals := r.URL.Query() // fetch request query params
	id,ok := vals["id"]
	isbn,name := vals["isbn"],vals["name"] // get name from request
	if ok {
		log.Print("going to insert database for name : ", name[0])
		stmt,err := db.Prepare("INSERT INTO movies SET id=?,isbn=?,name=?") // prepare an insert statment with a name placeholder that will be replaced with the name
		if err != nil {
			log.Print("error preparing query::", err)
			return
		}
		result, err := stmt.Exec(id[0],isbn[0],name[0]) // execute the statement
		if err != nil {
			log.Print("error executing query ::", err)
			return
		}
		lastId,err := result.LastInsertId() // get the id of last inserted record
		fmt.Fprintf(w,"Last inserted record Id is :: %s ", strconv.FormatInt(lastId,10)) // log id
	}else {
		fmt.Fprintf(w,"error occured while creating record in db for name :: %s",name[0])
	}
}

func readRecord(w http.ResponseWriter , r *http.Request){
	log.Print("printing all records from ")
	rows,err := db.Query("SELECT * FROM movies")
	if err != nil{
		log.Print("error executing select query :: ", err)
		return
	}

	var (
		id string 
		isbn string 
		name string
	)
	for rows.Next(){
		rows.Scan(&id,&isbn,&name);
		movie := Movie{ID:id, Isbn:isbn, Name:name }
		movies = append(movies, movie)
	}
	json.NewEncoder(w).Encode(movies)
}

func updateRecord(w http.ResponseWriter, r *http.Request){
	// get variables from dynamic route(url)
	vars := mux.Vars(r) 
	id := vars["id"]

	// get query parameters from url
	vals := r.URL.Query()
	name,ok := vals["name"]
	isbn := vals["isbn"]
	if ok {
		log.Printf("preparing to update database record for %s", name[0])

		stmt, err := db.Prepare("UPDATE movies SET name=?, isbn=? WHERE id=?") 
		if err != nil {
			log.Printf("error occured when preparing query :: ", err)
			return
		} 

		result, err := stmt.Exec(name[0],isbn[0],id)
		if err != nil {
			log.Printf("error occured when executing query :: ", err)
			return
		} 

		rowsAffected,err := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows updated in database are :: %d",rowsAffected)
	}else {
		fmt.Fprintf(w, "Error occurred while updating record in database for id :: %s", id)
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	stmt,err := db.Prepare("DELETE FROM movies WHERE id=?")
	if err != nil {
		log.Printf("error occured when preparing query :: ", err)
		return
	} 
	results,err := stmt.Exec(id)
	if err != nil {
		log.Printf("error occured when executing query :: ", err)
		return
	} 

	rowsAffected,err := results.RowsAffected()
	fmt.Fprintf(w, "Number of rows updated in database are :: %d",rowsAffected)

}

func addRoutes(r *mux.Router) *mux.Router{
	for _,route := range routes {
		r.HandleFunc(route.Pattern,route.HandlerFunc).Methods(route.Method)
	}
	return r
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	routeHandlers := addRoutes(r)

	fmt.Printf("Starting server on port 8000 \n")

	defer db.Close()

	err := http.ListenAndServe(CONN_HOST + ":" + CONN_PORT,routeHandlers)
	log.Fatal("ListenAndServe",err)
}