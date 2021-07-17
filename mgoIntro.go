package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type Movie1 struct {
	Name string `bson:"name"`
	Year string `bson:"year"`
	Directors []string `bson:"directors"`
	Writers []string `bson:"writers"`
	BoxOffice1 `bson:"boxOffice1"`
}

type BoxOffice1 struct {
	Budget uint64 `bson:"budget"`
	Gross uint64 `bson:"gross"`
}

func main() {
	session, err := mgo.Dial("172.17.0.3")

	if err != nil {
		panic(err)
	}

	defer session.Close()

	c := session.DB("appdb").C("movies")

	darkNight := &Movie1{
		Name: "The Dark Night",
		Year: "2008",
		Directors: []string{"Christopher Nolan"},
		Writers: []string{"Jonathan Nolan", "Christopher Nolan"},
		BoxOffice1: BoxOffice1 {
			Budget: 185000000,
			Gross: 533316061,
		},
	}

	err = c.Insert(darkNight)
	
	if err != nil {
		log.Fatal(err)
	}
	
	result := Movie1{}

	err = c.Find(bson.M{"boxOffice1.budget": bson.M{"$gt": 150000000}}).One(&result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Movie1:", result.Name)
}