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
	sync.Mutex
	ID        string // UUID
	Location  string // Location UUID
	User      string // User UUID
	IP        net.IP
	PingCount int // Number of times target has been pinged in this session
	TotalRTT  time.Duration
}

// State holds global state for this location.
type State struct {
	sync.Mutex
	PingSessions map[string]*Session // map[IPAddress]*Session
	PingChecker  Checker
	Pinger       *fastping.Pinger
}

// DefaultPingLimit is the default number of times we should ping a target.
const DefaultPingLimit = 5

// AverageLatency calculates average latency in ms
func (s *Session) AverageLatency() (averageLatency int64) {
	s.Lock()
	averageLatency = s.TotalRTT.Nanoseconds() / 1e6 / int64(s.PingCount)
	s.Unlock()
	return
}

// NewSession creates a new ping session.
func (state *State) NewSession(user string, ip net.IP) (s *Session, err error) {
	s = &Session{
		ID:       uuid.NewV4().String(),
		Location: state.PingChecker.ID,
		User:     user,
		IP:       ip,
	}

	state.Lock()
	state.PingSessions[ip.String()] = s
	err = state.Pinger.AddIP(ip.String())
	state.Unlock()
	if err != nil {
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

	s.Lock()
	s.TotalRTT += rtt
	s.PingCount++
	if s.PingCount > state.PingChecker.PingLimit {
		state.Pinger.RemoveIPAddr(addr)
	}
	s.Unlock()

	return
}

// Ping checks latency to user
func (state *State) Ping(ctx context.Context, in *Request) (r *Result, err error) {
	// validate IP address
	ip := net.ParseIP(in.IP)
	if ip == nil {
		err = fmt.Errorf("'%s' is not a valid ip address", in.IP)
		return
	}

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
			Done:     true,
		}
	case !ok:
		// Create a new session and add it to state.
		s, err = state.NewSession(in.User, ip)
		if err != nil {
			return
		}
		r = &Result{
			Location: state.PingChecker.ID,
			User:     s.User,
			Done:     false,
		}
	}
	return
}
