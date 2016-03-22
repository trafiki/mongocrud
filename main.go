package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//==============================================================================
// If you are connecting to a mgo database, you don't need to always redial
// else you are creating multiple connections, a mongo session can be used
// multiple times as many as you want. Hence I will only create it once then
// get the session to do what i want in the other functions.

// ses provides internal global session for our database operation,we will
// create new sessions of this one, once we have successfully connected.
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

	// We are setting the mod of operation to Monotonic to ensure we get
	// assurances that our writes and reads are synchornized and we not
	// cause inconsistencies.
	// NOTE: Never forget to do this.
	session.SetMode(mgo.Monotonic, true)

	ses = session
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

// NOTE:Question, if why do I want to declare this outside, do they
// have a relevance to not being declared within the main() function,
// Is your plan to provide these as global profiles forever?
// TODO: move them into your main function, else describe the reason why
// you have them here, it better be good, constructive and have a really
// solid reasoning behind it. O_O
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
	// Adding a profile to the database
	tunde.CreateOrUpdateProfile()
	saheed.CreateOrUpdateProfile()

	// Get all profiles in the Profiles collection.
	GetProfiles()

	// Get a profille from the Profiles collection.
	ReadProfile(tunde.Name)

	// TODO: WHAT happend to the UD (Update Delete) part of CRUD O_O?
}

// GetProfiles returns all the profile in the Profiles Collection
func GetProfiles() []Profile {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")

	var profiles []Profile
	// NOTE: Indent your code
	// TODO: Do you really think not checking your error is appropriate here, if am running this
	// on production server, my whole app will have a fault because i dont even know
	// if i found a valid record.
	// Those that make sense to you? O_O
	err = c.Find(bson.M{}).All(&profiles)

	return profiles
}

// ReadProfile returns the profile in the Profiles Collection with name equal to the id parameter (id == name)
func ReadProfile(id string) Profile {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")
	// NOTE:What happened to property indented code?
	// NOTE: why not do a:   "var profile Profile".
	// When declaring variables you do plan to initialize with a value
	// use the declaration not the assignative form.
	profile := Profile{}

	// TODO: Do you really think not checking your error is appropriate here, if am running this
	// on production server, my whole app will have a fault because i dont even know
	// if i found a valid record.
	// Those that make sense to you? O_O
	err = c.Find(bson.M{"name": id}).One(&profile)

	// NOTE: use pointers when returning something, why do i want to have a copy?
	// What if i want to change something?
	// Do i really want a copy of the profile and not a reference here
	// Does having a copy make sense O_O to having a reference?
	return profile
}

// DeleteProfile deletes the profile in the Profiles Collection with name equal to the id parameter (id == name)
func DeleteProfile(id string) bool {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")
	// TODO: Do you really think not checking your error is appropriate here, if am running this
	// on production server, my whole app will have a fault because i dont even know
	// if i found a valid record.
	// Those that make sense to you? O_O
	err = c.RemoveId(id)

	return true
}

// NOTE: not that it is wrong, but why not have a function that you pass a
// Profile to, to add it into the db. but you can peek and choose as you like.
// CreateOrUpdateProfile Creates or Updates the profile in the Profiles Collection with id parameter
func (p *Profile) CreateOrUpdateProfile() bool {
	// Get a session.
	session := getSession()
	defer session.Close()

	c := session.DB("ProfileService").C("Profiles")
	// NOTE: Indent your code
	_, err = c.UpsertId(p.Name, p)
	// NOTE: Indent your code
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())
		return false
	}
	// NOTE: Indent your code
	return true
}
