package user

import (
	"fmt"
	"sync"

	"github.com/satori/go.uuid"

	"golang.org/x/net/context"
)

// State holds the global state for users.
type State struct {
	sync.Mutex
	Users []*User
}

// ValidateAndSetName checks if a new name is provided and updates the user's
// name.
func (u *User) ValidateAndSetName(name string) (err error) {
	if name == "" {
		err = fmt.Errorf("%v is not a valid name for a user", name)
		return
	}
	u.Name = name
	return
}

// Search looks for Users in state matching provided User attributes. ID
// overrides all other parameters.
func (state *State) Search(ctx context.Context, in *User) (ul *UserList, err error) {
	ul = &UserList{Users: []*User{}}
	for _, u := range state.Users {
		if u.ID == in.ID {
			ul = &UserList{Users: []*User{u}}
			return
		}
		if u.Name == in.Name {
			ul.Users = append(ul.Users, u)
		}
	}
	return
}

// Update applies provided Delta to existing User in state.
func (state *State) Update(ctx context.Context, in *Delta) (u *User, err error) {
	ul, err := state.Search(ctx, &User{ID: in.User})
	if err != nil {
		return
	}

	u = ul.Users[0]
	err = u.ValidateAndSetName(in.Name)
	if err != nil {
		return
	}
	return
}

// Delete deletes User specified by provided User object from state.
func (state *State) Delete(ctx context.Context, in *User) (u *User, err error) {
	state.Lock()
	var index = -1
	for i, _u := range state.Users {
		if _u.ID == in.ID {
			index = i
			u = _u
		}
	}

	// If the user is found, pop it from state.
	if index != -1 {
		state.Users = append(state.Users[:index], state.Users[index+1:]...)
	} else {
		err = fmt.Errorf("User %+v not found in state", in)
	}

	state.Unlock()
	return
}

// Create creates User specified by provided User object and adds it to state.
func (state *State) Create(ctx context.Context, in *User) (u *User, err error) {
	state.Lock()
	in.ID = uuid.NewV4().String()
	state.Users = append(state.Users, in)
	state.Unlock()

	u = in
	return
}
