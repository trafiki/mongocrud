package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Profile - is the memory representation of one user profile
type Profile struct {
	Name        string `json:"username"`
	Password    string `json:"password"`
	Age         int    `json:"age"`
	LastUpdated time.Time
}

var tunde = &Profile{
	Name:        "tunde",
	Password:    "akerele",
	Age:         27,
	LastUpdated: time.Now(),
}

var saheed = &Profile{
	Name:        "saheed",
	Password:    "ajibulu",
	Age:         27,
	LastUpdated: time.Now(),
}

func main() {

	tunde.CreateOrUpdateProfile()
	saheed.CreateOrUpdateProfile()
	GetProfiles()
	ReadProfile(tunde.Name)
}

// GetProfiles - Returns all the profile in the Profiles Collection
func GetProfiles() []Profile {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return nil
	}
	defer session.Close()

	// Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	var profiles []Profile
	err = c.Find(bson.M{}).All(&profiles)

	return profiles
}

// ReadProfile returns the profile in the Profiles Collection with name equal to the id parameter (id == name)
func ReadProfile(id string) Profile {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return Profile{}
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	profile := Profile{}
	err = c.Find(bson.M{"name": id}).One(&profile)

	return profile
}

// DeleteProfile deletes the profile in the Profiles Collection with name equal to the id parameter (id == name)
func DeleteProfile(id string) bool {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return false
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	err = c.RemoveId(id)

	return true
}

// CreateOrUpdateProfile Creates or Updates (Upsert) the profile in the Profiles Collection with id parameter
func (p *Profile) CreateOrUpdateProfile() bool {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return false
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	_, err = c.UpsertId(p.Name, p)
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())
		return false
	}
	return true
}
