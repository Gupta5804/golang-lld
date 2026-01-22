package main

import (
	"fmt"
)

// -- behavior/interface
type UserSaver interface {
	SaveUser(name string)
}

// -- concrete implementation (SQL Database)
type SQLDatabase struct {
	dbUrl string
}

func (s *SQLDatabase) SaveUser(name string) {
	fmt.Printf("Insert INTO users VALUES ('%s') via SQL DB:%s\n", name, s.dbUrl)
}

// -- concrete implementation (File Database)
type FileDatabase struct {
	filename string
}

func (f *FileDatabase) SaveUser(name string) {
	fmt.Printf("Insert INTO users VALUES ('%s') via File name: %s\n", name, f.filename)
}

// -- Business Logic --
type UserService struct {
	saver UserSaver
}

func (u *UserService) RegisterUser(name string) {
	// logic: validate name
	fmt.Println("Validating user...")

	// delegation
	u.saver.SaveUser(name)
}

// -- builder/constructor function --
func NewUserService(u UserSaver) *UserService {
	return &UserService{saver: u}
}
func NewDatabase(env string) (UserSaver, error) {
	// if env == "production" {
	// 	return &SQLDatabase{dbUrl: "sqlDB"}, nil
	// } else if env == "local" {
	// 	return &FileDatabase{filename: "file_01"}, nil
	// } else {
	// 	return nil, fmt.Errorf("unknown environment: %s", env)
	// }
	switch env {
	case "production":
		return &SQLDatabase{dbUrl: "sqlDB"}, nil
	case "local":
		return &FileDatabase{filename: "file_01"}, nil
	default:
		return nil, fmt.Errorf("unknown environment : %s", env)
	}
}

func main() {
	// BAD: Manual assembly of concrete types

	db, err := NewDatabase("loca")
	// sqlDB := &SQLDatabase{dbUrl: "mysql://localhost:3306"}

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic("aborting")
	}
	userService := NewUserService(db)

	userService.RegisterUser("Alice")
}
