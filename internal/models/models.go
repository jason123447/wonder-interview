package models

import "slices"

type User struct {
	ID           int    `json:"id"`
	Account      string `json:"account"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
}

var MockUsers = []User{
	{ID: 1, Account: "user1", Password: "", PasswordHash: "$2y$14$KZ3xNka0Xxzsm/uJy2g4gOM/QsQvEUF75TEceDAEtTohgm6QFEU.6"},
	{ID: 2, Account: "user2", Password: "", PasswordHash: "$2y$14$KZ3xNka0Xxzsm/uJy2g4gOM/QsQvEUF75TEceDAEtTohgm6QFEU.6"},
}

func FindUserByID(userID int) *User {
	idx := slices.IndexFunc(MockUsers, func(u User) bool {
		return u.ID == userID
	})
	if idx == -1 {
		return nil
	}
	return &MockUsers[idx]
}

func FindUserByAccount(account string) *User {
	idx := slices.IndexFunc(MockUsers, func(u User) bool {
		return u.Account == account
	})
	if idx == -1 {
		return nil
	}
	return &MockUsers[idx]
}
