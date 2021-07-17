package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	session *mgo.Session
	collection *mgo.Collection	
}

type Movie struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Year string `json:"year" bson:"year"`
	Directors []string `json:"directors" bson:"directors"`
	Writers []string `json:"writers" bson:"writers"`
	BoxOffice BoxOffice `json:"boxOffice" bson:"boxOffice"`
}

type BoxOffice struct {
	Budget uint64 `json:"budget" bson:"budget"`
	Gross uint64 `json:"gross" bson:"gross"`
}

func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	var movie Movie

	err := db.collection.Find(bson.M{
			"_id": bson.ObjectIdHex(vars["id"]),
		}).One(&movie)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _:= json.Marshal(movie)

		// fmt.Printf("type is %T\n",response)

		w.Write(response)
	}
}

func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := db.collection.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})

	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Deleted successfully!"))
	}
}

func (db *DB) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	vars := mux.Vars(r)
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &movie)

	// create hash id
	err := db.collection.Update(bson.M{"_id": bson.ObjectIdHex(vars["id"])},
								bson.M{"$set": &movie})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Updated Successfully!"))
	}
}

func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &movie)

	// create hash id
	movie.ID = bson.NewObjectId()

	err := db.collection.Insert(movie)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _:= json.Marshal(movie)

		// fmt.Printf("type is %T\n",response)

		w.Write(response)
	}
}

func main() {
	session, err := mgo.Dial("172.17.0.3")

	if err != nil {
		panic(err)
	}

	defer session.Close()

	c := session.DB("appdb").C("movies")

	db := &DB{session: session, collection: c}

	r := mux.NewRouter()

	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.GetMovie).Methods("get")
	r.HandleFunc("/v1/movies", db.PostMovie).Methods("post")
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.UpdateMovie).Methods("put")
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.DeleteMovie).Methods("delete")

	srv := &http.Server {
		Handler: r,
		Addr: "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}