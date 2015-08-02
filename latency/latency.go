package latency

import (
	"fmt"
	"net"
	"time"

	fastping "github.com/tatsushid/go-fastping"
	"golang.org/x/net/context"
)

type pinger fastping.Pinger

// Checker implements latency checker server
type Checker Config

// DefaultPingCount is the default number of times we should ping a target.
const DefaultPingCount = 10

func (p *pinger) ping(ip string, out <-chan time.Duration) {
	p.AddIP(in.IP)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		out <- rtt
	}
	p.OnIdle = func() {
		return
	}
	err = p.RunLoop()
	if err != nil {
		return
	}
	return
}

// GetContext returns context for Checker session.
func (c *Checker) GetContext(user string) (ctx context.Context) {
	ctx = context.WithValue(context.Context{}, "user", user)
	return
}

// Ping checks latency to user
func (c *Checker) Ping(ctx context.Context, in *Request) (r *Result, err error) {
	out := make(chan time.Duration)

	// validate IP address
	ip := net.ParseIP(in.IP)
	if ip == nil {
		err = fmt.Errorf("'%v' is not a valid ip address")
		return
	}

	// Ping in a goroutine
	p := fastping.NewPinger()
	go p.ping(in.IP, out)

	var totalDuration time.Duration
	pingCount, ok := ctx.Value("pingCount").(int)
	if !ok {
		pingCount := DefaultPingCount
	}

	for i := 0; i < pingCount; i++ {
		duration := <-out
		totalDuration += totalDuration
	}
	p.Stop()

	// average latency in milliseconds
	var averageLatency float64
	averageLatency = totalDuration.Nanoseconds() / 1e6 / pingCount

	r = &Result{
		Location: c.ID,
		Latency:  averageLatency,
		User:     in.User,
		Pinging:  false,
	}
	return
}
