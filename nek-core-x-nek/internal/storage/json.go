package storage

import (
	"encoding/json"
	"errors"
	"github.com/qopix/nek/internal/model"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu    sync.Mutex
	path  string
	state model.State
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, state: model.State{Nodes: []model.Node{}, Users: []model.User{}, Subscriptions: []model.Subscription{}}}
	b, e := os.ReadFile(path)
	if e == nil {
		if e = json.Unmarshal(b, &s.state); e != nil {
			return nil, e
		}
		if s.state.Nodes == nil {
			s.state.Nodes = []model.Node{}
		}
		if s.state.Users == nil {
			s.state.Users = []model.User{}
		}
		if s.state.Subscriptions == nil {
			s.state.Subscriptions = []model.Subscription{}
		}
	} else if !errors.Is(e, os.ErrNotExist) {
		return nil, e
	}
	if e = s.save(); e != nil {
		return nil, e
	}
	return s, nil
}
func (s *Store) View() model.State {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, _ := json.Marshal(s.state)
	var out model.State
	_ = json.Unmarshal(b, &out)
	return out
}
func (s *Store) Update(fn func(*model.State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := fn(&s.state); e != nil {
		return e
	}
	return s.save()
}
func (s *Store) save() error {
	if e := os.MkdirAll(filepath.Dir(s.path), 0700); e != nil {
		return e
	}
	b, e := json.MarshalIndent(s.state, "", "  ")
	if e != nil {
		return e
	}
	tmp := s.path + ".tmp"
	if e = os.WriteFile(tmp, b, 0600); e != nil {
		return e
	}
	return os.Rename(tmp, s.path)
}
