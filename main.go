package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ses *mgo.Session

func getSession() *mgo.Session {
	if ses != nil {
		// We do a copy because we want to have a unique session from what
		// another might be using but we still want to multiplex on the same
		// connection underneath it, its much cleaner.
		return ses.Copy()
	}

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		panic(err)
	}

	// We are setting the mode of operation to Monotonic to ensure we get
	// assurances that our writes and reads are synchornized and we not
	// cause inconsistencies.
	session.SetMode(mgo.Monotonic, true)

	ses = session.Copy()
	return session
}

//==============================================================================

// Profile - is the memory representation of one user profile
type Profile struct {
	Name        string `json:"username"`
	Password    string `json:"password"`
	Age         int    `json:"age"`
	LastUpdated time.Time
}

func main() {

	// Adding profiles to the database
	CreateProfile("Tom", "tominson", 27)
	CreateProfile("Dick", "dickinson", 27)
	CreateProfile("Harry", "harrison", 27)

	// Get all profiles in the Profiles collection.
	GetProfiles()

	// Get a profille from the Profiles collection.
	ReadProfile("Tom")

	// Delete a profile from the profiles collection.
	DeleteProfile("Dick")

	// Update a profile in the profiles collection.
	UpdateProfile("Harry", "harry", 25)

}

// GetProfiles returns all the profile in the Profiles Collection
func GetProfiles() *[]Profile {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")

	var profiles []Profile

	err := c.Find(bson.M{}).All(&profiles)
	if err != nil {
		log.Println("Error getting profiles: ", err.Error())
		return &profiles
	}

	return &profiles
}

// ReadProfile returns the profile in the Profiles Collection with name equal to the id parameter (id == name)
func ReadProfile(id string) *Profile {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")

	var profile Profile

	err := c.Find(bson.M{"name": id}).One(&profile)
	if err != nil {
		log.Println("Error reading profile: ", err.Error())
		return &profile
	}

	return &profile
}

// DeleteProfile deletes the profile in the Profiles Collection with name equal to the id parameter (id == name)
func DeleteProfile(id string) bool {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")

	err := c.Remove(bson.M{"name": id})
	if err != nil {
		log.Println("Error removing profile: ", err.Error())
		return false
	}

	fmt.Println(id, "has been removed.")

	return true
}

// CreateProfile creates a profile for the Profiles collection.
func CreateProfile(id string, password string, age int) bool {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")

	index := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		log.Println(err.Error())
	}

	err = c.Insert(&Profile{id, password, age, time.Now()})

	if err != nil {
		log.Println("Error creating Profile: ", err.Error())
		return false
	}

	return true
}

// UpdateProfile updates a profile in the profiles collection with the id parameter.
func UpdateProfile(id string, password string, age int) bool {
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")

	querier := bson.M{"name": id}
	change := bson.M{"$set": bson.M{"Password": password, "Age": age, "LastUpdated": time.Now()}}

	err := c.Update(querier, change)
	if err != nil {
		log.Println("Error updating profile: ", err.Error())
		return false
	}

	return true
}
