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
	PingCount int // Number of times to ping target
	TotalRTT  time.Duration
	State     int
	RWMutex   sync.RWMutex

	pings chan time.Duration
}

// Ping States
const (
	Waiting = iota
	Pinging
	Done
)

// State holds global state for this location.
type State struct {
	PingSessions []*Session
	PingChecker  Checker
	Mutex        sync.Mutex
}

// DefaultPingCount is the default number of times we should ping a target.
const DefaultPingCount = 10

// Ping the host based on attributes from its Session.
func (s *Session) Ping() {
	s.RWMutex.Lock()
	s.State = Pinging

	p := fastping.NewPinger()
	p.AddIP(ip)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		s.pings <- rtt
	}
	p.OnIdle = func() {
		return
	}
	p.RunLoop()

	var totalDuration time.Duration
	for i := 0; i < pingCount; i++ {
		duration := <-s.pings
		s.TotalRTT += duration
	}
	p.Stop()

	s.State = Done
	s.RWMutex.Unlock()
	return
}

// AverageLatency calculates average latency in ms
func (s *Session) AverageLatency(totalRTT time.Duration) (averageLatency int64) {
	averageLatency = totalRTT.Nanoseconds() / 1e6 / int64(s.PingCount)
}

// NewSession creates a new ping session.
func NewSession(state *State, in *Request) (newSession *Session) {
	// validate IP address
	ip := net.ParseIP(in.IP)
	if ip == nil {
		err = fmt.Errorf("'%v' is not a valid ip address")
		return
	}

	pings := make(chan time.Duration)

	newSession = &Session{
		ID:        uuid.NewV4(),
		Location:  state.PingChecker.ID,
		User:      in.User,
		IP:        ip,
		PingCount: state.PingChecker.PingCount,
		TotalRTT:  time.ParseDuration("0"),
		State:     Waiting,
		pings:     pings,
	}
	return
}

// RenderResult returns a result response based on session.
func (s *Session) RenderResult() (r *Result) {
	switch {
	case s.State == Done:
		// average latency in milliseconds
		r = &Result{
			Location: state.PingChecker.ID,
			Latency:  s.AverageLatency(s.TotalRTT),
			User:     s.User,
			Pinging:  false,
		}
		return
	case s.State != Done:
		r = &Result{
			Location: state.PingChecker.ID,
			User:     s.User,
			Pinging:  true,
		}
		return
	}
}

// FilterByRequest filters sessions based on request parameters.
func (s *Session) FilterByRequest(in *Request) (ok bool) {
	switch {
	case s.User != in.User:
		ok = false
		return
	case s.IP != in.IP:
		ok = false
		return
	}
	return
}

// GetSession gets a session based on User and IP.
func (state *State) GetSession(in *Request) (s *Session, ok bool) {
	for session := range state.PingSessions {
		found = session.FilterByRequest(in)
		if found {
			s = session
			ok = true
			return
		}
	}
	return
}

// Ping checks latency to user
func (state *State) Ping(ctx context.Context, in *Request) (r *Result, err error) {
	var s *Session
	s, ok = state.GetSession(in)
	if !ok {
		// Create a new session and add it to state.
		s = func(state, in) (s *Session) {
			s = NewSession(state, in)
			state.Mutex.Lock()
			state.PingSessions = append(state.PingSessions, s)
			state.Mutex.Unlock()
		}(state, in)

		// Ping in a goroutine
		go s.Ping()
	}

	r = s.RenderResult()
	return
}
