package domain

import "sync"

type User struct {
	ID           string
	IsSuspended  bool
	PasswordHash string
	mu           sync.RWMutex
}

func(u *User) Suspend(){
	u.mu.Lock()
	defer u.mu.Unlock()
	u.IsSuspended = true
}
func(u *User) Restore(){
	u.mu.Lock()
	defer u.mu.Unlock()
	u.IsSuspended = false
}
func(u *User) ChangePassword(newHash string){
	u.mu.Lock()
	defer u.mu.Unlock()
	u.PasswordHash = newHash
}

