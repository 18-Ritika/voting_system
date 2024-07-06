package models

import (
	"github.com/dgrijalva/jwt-go"
	"sync"
)

type User struct {
	ID        int
	Name      string
	VotedTeam string
}

type Vote struct {
	ID   int
	Team Team
}

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type VotingSystem struct {
	UserVote  map[User]Vote
	VoteCount map[Team]int
}

func NewVotingSystem() *VotingSystem {
	return &VotingSystem{
		UserVote:  make(map[User]Vote),
		VoteCount: make(map[Team]int),
	}
}

type Database struct {
	syMutex sync.Mutex
	result  map[string]VotingSystem
}

type Message struct {
	Type  string `json:"type"`
	Vote  string `json:"vote,omitempty"`
	Token string `json:"token,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
