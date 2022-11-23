package main

import ( "fmt" ; "log" ; "strings" ; "encoding/json" ; "math/rand" ; "net/http" ; "strconv" ; "github.com/gorilla/mux")

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8000"
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
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}
var routes []Route = []Route{
	{"getMovies","GET","/movies", getMovies},
	{"getMovie","GET","/movies/{id}", getMovie},
	{"createMovie","POST","/movies", createMovie},
	{"updateMovie","PUT","/movies/{id}", updateMovie},
	{"deleteMovie","DELETE","/movies/{id}", deleteMovie},
}

var movies []Movie
var moviesV1 []Movie
var moviesV2 []Movie

func init(){
	movies = []Movie{
		{ID:"1",Isbn:"3849574", Name:"Spider Man", Director: &Director{Firstname:"Hello", Lastname:"Bagle"}},
	}
	moviesV1 = []Movie{
		{ID:"1",Isbn:"3849574", Name:"Wonder Woman", Director: &Director{Firstname:"Hello", Lastname:"Bagle"}},
	}
	moviesV2 = []Movie{
		{ID:"1",Isbn:"3849574", Name:"Super Man", Director: &Director{Firstname:"Hello", Lastname:"Bagle"}},
	}
}

func getMovies(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", CONN_HOST + ":" + CONN_PORT)
	if strings.HasPrefix(r.URL.Path,"/v1"){
		json.NewEncoder(w).Encode(moviesV1)
	}else if strings.HasPrefix(r.URL.Path,"/v2"){
		json.NewEncoder(w).Encode(moviesV2)
	}else {
		json.NewEncoder(w).Encode(movies)
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, elem := range movies {
		if elem.ID == params["id"] {
			movies = append(movies[:idx], movies[idx + 1 :]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for _, elem := range movies {
		if elem.ID == params["id"] {
			json.NewEncoder(w).Encode(elem)
			return 
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	movies = append(movies,movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, elem := range movies {
		if elem.ID == params["id"] {
			movies = append(movies[:idx],movies[idx + 1 :]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies,movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
}

func addRoutes(r *mux.Router) *mux.Router{
	for _,route := range routes {
		r.HandleFunc(route.Pattern,route.HandlerFunc).Methods(route.Method)
	}
	return r
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	version1 := r.PathPrefix("/v1").Subrouter()
	version2 := r.PathPrefix("/v2").Subrouter()

	routeHandlers := addRoutes(r)

	addRoutes(version1)
	addRoutes(version2)

	fmt.Printf("Starting server on port 8000 \n")


	err := http.ListenAndServe(CONN_HOST + ":" + CONN_PORT,routeHandlers)
	log.Fatal("ListenAndServe",err)
}