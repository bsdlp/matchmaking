package lobby

import (
	"errors"
	"fmt"
	"sync"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// Lobby holds state for a matchmaking lobby, implements MatchMakingServer.
type Lobby struct {
	ID     string
	Users  []*User
	Config Config
	Mutex  sync.Mutex
}

// Config holds parameters for a matchmaking lobby.
type Config struct {
	TeamSize     int
	TeamCount    int
	LobbySize    int // Should always be TeamSize * TeamCount
	PrivacyLevel int
	Checks       []string
}

// State is an object that holds all of the lobbies and provides a mutex lock.
type State struct {
	Lobbies []*Lobby
	Mutex   sync.Mutex
}

// PrivacyCheck describes the privacy validation interface.
type PrivacyCheck interface {
	Check(lobby *Lobby, user *User) (err error)
}

// Privacy levels
const (
	PublicMatch = iota
	PrivateMatch
)

// Lobbies is the State used for testing.
var Lobbies = &State{}

// NewLobby creates a new lobby.
func NewLobby(cfg Config) (lobby *Lobby) {
	lobby = &Lobby{
		ID:     uuid.NewV4().String(),
		Config: cfg,
	}
	return
}

// ValidatePublic checks to make sure the lobby is not full. Should be protected
// by a lock.
func (l *Lobby) ValidatePublic(user *User) (err error) {
	if len(l.Users) >= l.Config.LobbySize {
		err = errors.New("Lobby is full.")
	}
	return
}

// ValidatePrivate checks lobby fullness as well as whether or not user was
// invited by an existing member of the lobby. Should be protected by a lock.
func (l *Lobby) ValidatePrivate(user *User) (err error) {
	return
}

// Validate calls corresponding privacy level validator to determine if user is
// allowed to join the lobby. Should be protected by a lock.
func (l *Lobby) Validate(user *User) (err error) {
	switch {
	case l.Config.PrivacyLevel == PublicMatch:
		err = l.ValidatePublic(user)
		return
	case l.Config.PrivacyLevel == PrivateMatch:
		err = l.ValidatePrivate(user)
		return
	}
	err = fmt.Errorf(
		"Validator for privacy level %v not implemented", l.Config.PrivacyLevel)
	return
}

// Join adds a user to the lobby. Locks the lobby to validate and add user.
func (l *Lobby) Join(user *User) (err error) {
	l.Mutex.Lock()
	err = l.Validate(user)
	if err != nil {
		return
	}
	l.Users = append(l.Users, user)
	l.Mutex.Unlock()
	return
}

// FindLobby returns lobby that matches provided lobby ID. If
// the attributes do not match any lobby, then a new lobby is created with
// provided attributes. This method should get overridden in code used for
// more than testing.
func (s *State) FindLobby(l *Lobby) (rl *Lobby, err error) {
	s.Mutex.Lock()
	for _, _l := range s.Lobbies {
		_l.Mutex.Lock()
		// If lobby IDs are an exact match, return it.
		// TODO: Greedy filter (missing search attributes don't count against
		// filter)
		if _l.ID == l.ID {
			rl = _l
			return
		}
		_l.Mutex.Unlock()
	}

	// If no lobby is found with matching attributes, create a new lobby.
	rl = NewLobby(l.Config)

	// Add lobby to state.
	s.Lobbies = append(s.Lobbies, rl)
	s.Mutex.Unlock()
	return
}

// Join adds a user to the lobby.
func (s *State) Join(ctx context.Context, hello *Hello) (joined *Joined, err error) {
	l, err := s.FindLobby(&Lobby{ID: hello.Lobby})
	if err != nil {
		return
	}

	err = l.Join(hello.User)
	if err != nil {
		return
	}

	joined = &Joined{
		Lobby: l.ID,
		User:  l.Users,
	}

	return
}
