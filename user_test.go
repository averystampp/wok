package wok

import (
	"testing"
)

var config = &DbConfig{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "docker",
	Dbname:   "postgres",
}

func TestCreateUser(t *testing.T) {
	db, err := DbConnect(config)
	if err != nil {
		t.Error(err)
	}
	Database = db
	u := new(User)

	u.Email = "test"
	u.Password = "password"

	if err := CreateUser(u); err != nil {
		t.Error(err)
	}

}

func TestLogin(t *testing.T) {
	username := "test"
	password := "password"

	_, err := Login(username, password)

	if err != nil {
		t.Error(err)
	}
}

func TestLogout(t *testing.T) {
	var id string
	row := Database.QueryRow("SELECT session_id from user where email=test")
	row.Scan(&id)

	if err := Logout(id); err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	var id int
	row := Database.QueryRow("SELECT id from user where email=test")
	row.Scan(&id)

	if err := DeleteUser(id); err != nil {
		t.Error(err)
	}
}
