package main

import "fmt"

// -- behavior/interface
type UserSaver interface {
	SaveUser(name string)
}

// -- concrete implementation (SQL Database)
type SQLDatabase struct {
	dbUrl string
}

func (s *SQLDatabase) SaveUser(name string) {
	fmt.Printf("Insert INTO users VALUES ('%s') via %s\n", name, s.dbUrl)
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

func main() {
	// BAD: Manual assembly of concrete types
	sqlDB := &SQLDatabase{dbUrl: "mysql://localhost:3306"}

	userService := NewUserService(sqlDB)

	userService.RegisterUser("Alice")
}
