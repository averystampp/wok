package wok

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// users id in database
	ID int `json:"id"`

	// first name of user
	FirstName string `json:"firstname"`

	//last name of user
	LastName string `json:"lastname"`

	// users password
	Password string `json:"password"`

	// session id of user, conforms to uuid
	SessionID uuid.UUID `json:"sessionid"`

	// users role, defaults to "user" but "admin" and "inactive" are also available
	Role string `json:"role"`

	Logged_in string `json:"loggedin"`
}

// takes a user object and a db, parses user struct to ensure minimum fields are filled in
// creates a new user with supplied info and a uuid, also sets the users role to "user"
// returns the full user object for post processing like a direcct login on frontend, also
// returns an error.
func CreateUser(user *User) error {

	// Check if firstname, lastname, and password are blank
	// returns error if they are. These fields are required
	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	// Create new UUID for user when their account is created
	user.SessionID = uuid.New()

	// Set the role to a user, this is the standard user authorization
	user.Role = "user"

	// set the user to logged out, this requires them to sign after creating their account
	user.Logged_in = "false"

	// hash users password before inserting into the database
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// set the users password as a string of the hash
	user.Password = string(hash)

	// Insert user values into users table in database,
	// if there is an error it will return an empty user and the error
	_, err = Database.Exec("INSERT INTO users (email, lastname, password, role, session_id, logged_in) VALUES ($1, $2, $3, $4, $5, $6)",
		user.FirstName, user.LastName, user.Password, user.Role, user.SessionID, user.Logged_in)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// return nil for no errors
	return nil
}

// takes a username and password, then queries db for the hashed password
// and compares the user submitted password to the hash, will return a uuid on success,
// status unauthorized if the hash and pass do not match, or a server error if the database
// cannot return a response.
func Login(username, password string) (string, error) {
	qs := "SELECT password FROM users WHERE email=$1 LIMIT 1"

	// query db with the supplied username, then writes the hashed password to var pass
	var pass string
	row := Database.QueryRow(qs, username)
	row.Scan(&pass)

	// compare users password to the hash pulled from the db
	// if they don't match then server responds with a unauthorized request
	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	// create a new uuid for the user, if there is an error then the handler will
	// return the error and break the request
	uuid := uuid.New().String()

	// update the user where their firstname is equal to input
	// set theie logged in status to true, and also update their session_id
	// this ensures that the user always has a fresh uuid upon login
	update := "UPDATE users SET session_id=$1, logged_in=$2 WHERE email=$3"
	_, err := Database.Exec(update, uuid, "true", username)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("cannot process request")
	}

	// return the new session_id and nil for error
	return uuid, nil

}

// logs user out, takes their id and sets the logged in status to false...
// probably could hack something better together, the users uuid should be invalidated
// need to add a context to the user to create an expiry e.g. ctx.Set("expire in", time.Minutes * 30)
func Logout(id string) error {
	qs := "UPDATE users SET logged_in=FALSE WHERE session_id=$1"

	_, err := Database.Exec(qs, id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("cannot process request")
	}

	return nil
}

func DeleteUser(id int) error {
	qs := "DELETE FROM users WHERE ID=$1"
	_, err := Database.Exec(qs, id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("cannot process request")
	}

	return nil
}

// same as user function but will create a an admin.
// This by default is only accesable with the CLI flag "createuser"
func CreateAdmin(user *User) error {
	// Check if firstname, lastname, and password are blank
	// returns error if they are. These fields are required
	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	// Create new UUID for user when their account is created
	user.SessionID = uuid.New()

	// Set the role to a user, this is the standard user authorization
	user.Role = "admin"

	// set the user to logged out, this requires them to sign after creating their account
	user.Logged_in = "false"

	// hash users password before inserting into the database
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// set the users password as a string of the hash
	user.Password = string(hash)

	// Insert user values into users table in database,
	// if there is an error it will return an empty user and the error
	_, err = Database.Exec("INSERT INTO users (email, lastname, password, role, session_id, logged_in) VALUES ($1, $2, $3, $4, $5, $6)",
		user.FirstName, user.LastName, user.Password, user.Role, user.SessionID, user.Logged_in)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// return nil for no errors
	return nil
}
