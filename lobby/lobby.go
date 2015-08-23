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
	sync.Mutex
	ID     string
	Users  []*User
	Config Config
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
	sync.Mutex
	Lobbies []LobbyServer
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

// NewLobby creates a new lobby.
func NewLobby(cfg Config) (l LobbyServer) {
	l = &Lobby{
		ID:     uuid.NewV4().String(),
		Config: cfg,
	}
	return
}

// validatePublic checks to make sure the lobby is not full. Should be protected
// by a lock.
func (l *Lobby) validatePublic(user *User) (err error) {
	if len(l.Users) >= l.Config.LobbySize {
		err = errors.New("Lobby is full.")
	}
	return
}

// validatePrivate checks lobby fullness as well as whether or not user was
// invited by an existing member of the lobby. Should be protected by a lock.
func (l *Lobby) validatePrivate(user *User) (err error) {
	return
}

// validate calls corresponding privacy level validator to determine if user is
// allowed to join the lobby. Should be protected by a lock.
func (l *Lobby) validate(user *User) (err error) {
	switch {
	case l.Config.PrivacyLevel == PublicMatch:
		err = l.validatePublic(user)
		return
	case l.Config.PrivacyLevel == PrivateMatch:
		err = l.validatePrivate(user)
		return
	}
	err = fmt.Errorf(
		"Validator for privacy level %v not implemented", l.Config.PrivacyLevel)
	return
}

// Join adds a user to the lobby. Locks the lobby to validate and add user.
func (l *Lobby) Join(ctx context.Context, h *Hello) (j *Joined, err error) {
	l.Lock()
	err = l.validate(h.User)
	if err != nil {
		return
	}
	l.Users = append(l.Users, h.User)
	l.Unlock()
	return
}

// Leave removes a user from the lobby
func (l *Lobby) Leave(ctx context.Context, g *Goodbye) (left *Left, err error) {
	l.Lock()

	// remove user from lobby
	users := l.Users[:0]
	for _, u := range l.Users {
		if u.Id != g.User.Id {
			users = append(users, u)
		}
	}
	l.Users = users

	left = &Left{
		Lobby: l.ID,
		Left:  true,
	}
	l.Unlock()
	return
}

// Check returns lobby status
func (l *Lobby) Check(ctx context.Context, s *Sup) (status *Status, err error) {
	return
}
