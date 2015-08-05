package latency

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	fastping "github.com/tatsushid/go-fastping"
	"golang.org/x/net/context"
)

// Checker implements latency checker server
type Checker Config

// Session is the object that holds a user's ping service session.
type Session struct {
	ID        string // UUID
	Location  string // Location UUID
	User      string // User UUID
	IP        net.IP
	PingCount int // Number of times target has been pinged in this session
	TotalRTT  time.Duration
	Mutex     sync.Mutex
}

// State holds global state for this location.
type State struct {
	PingSessions map[string]*Session // map[IPAddress]*Session
	PingChecker  Checker
	Pinger       *fastping.Pinger
	Mutex        sync.Mutex
}

// DefaultPingLimit is the default number of times we should ping a target.
const DefaultPingLimit = 5

// NewSession creates a new ping session.
func NewSession(state *State, in *Request) (newSession *Session, err error) {
	// validate IP address
	ip := net.ParseIP(in.IP)
	if ip == nil {
		err = fmt.Errorf("'%s' is not a valid ip address", in.IP)
		return
	}

	newSession = &Session{
		ID:       uuid.NewV4().String(),
		Location: state.PingChecker.ID,
		User:     in.User,
		IP:       ip,
	}
	return
}

// AverageLatency calculates average latency in ms
func (s *Session) AverageLatency() (averageLatency int64) {
	s.Mutex.Lock()
	averageLatency = s.TotalRTT.Nanoseconds() / 1e6 / int64(s.PingCount)
	s.Mutex.Unlock()
	return
}

// FilterByRequest filters sessions based on request parameters.
func (s *Session) FilterByRequest(in *Request) (ok bool) {
	switch {
	case s.User != in.User:
		ok = false
		return
	case s.IP.String() != in.IP:
		ok = false
		return
	}
	return
}

// NewState returns a new State object.
func NewState(id string) (state *State) {
	pinger := fastping.NewPinger()
	state = &State{
		PingSessions: make(map[string]*Session),
		PingChecker: Checker{
			ID:        id,
			PingLimit: DefaultPingLimit,
		},
		Pinger: pinger,
	}

	pinger.OnRecv = state.onRecv
	pinger.RunLoop()
	return
}

func (state *State) onRecv(addr *net.IPAddr, rtt time.Duration) {
	s := state.PingSessions[addr.String()]

	s.Mutex.Lock()
	s.TotalRTT += rtt
	s.PingCount++
	if s.PingCount > state.PingChecker.PingLimit {
		state.Pinger.RemoveIPAddr(addr)
	}
	s.Mutex.Unlock()

	return
}

// Ping checks latency to user
func (state *State) Ping(ctx context.Context, in *Request) (r *Result, err error) {
	ip := net.ParseIP(in.IP)

	var s *Session
	var ok bool
	s, ok = state.PingSessions[ip.String()]
	switch {
	case ok:
		// TODO: implement cache instead of deleting immediately or something.
		delete(state.PingSessions, ip.String())

		r = &Result{
			Location: state.PingChecker.ID,
			Latency:  s.AverageLatency(),
			User:     s.User,
			Pinging:  false,
		}
	case !ok:
		// Create a new session and add it to state.
		s, err = NewSession(state, in)
		if err != nil {
			return
		}
		state.Mutex.Lock()
		state.PingSessions[s.IP.String()] = s
		err = state.Pinger.AddIP(ip.String())
		state.Mutex.Unlock()
		if err != nil {
			return
		}

		r = &Result{
			Location: state.PingChecker.ID,
			User:     s.User,
			Pinging:  true,
		}
	}
	return
}
