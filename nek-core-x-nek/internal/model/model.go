package model

import "time"

type Node struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	AgentURL  string    `json:"agent_url,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
type Subscription struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Protocol  string    `json:"protocol"`
	NodeID    string    `json:"node_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
type State struct {
	Nodes         []Node         `json:"nodes"`
	Users         []User         `json:"users"`
	Subscriptions []Subscription `json:"subscriptions"`
}
