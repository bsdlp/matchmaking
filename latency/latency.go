package latency

import (
	"fmt"
	"net"
	"time"

	fastping "github.com/tatsushid/go-fastping"
	"golang.org/x/net/context"
)

// Checker implements latency checker server
type Checker Config

// DefaultPingCount is the default number of times we should ping a target.
const DefaultPingCount = 10

func ping(p *fastping.Pinger, ip string, out chan time.Duration) {
	p.AddIP(ip)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		out <- rtt
	}
	p.OnIdle = func() {
		return
	}
	p.RunLoop()
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
	go ping(p, in.IP, out)

	var totalDuration time.Duration
	pingCount, ok := ctx.Value("pingCount").(int)
	if !ok {
		pingCount = DefaultPingCount
	}

	for i := 0; i < pingCount; i++ {
		duration := <-out
		totalDuration += duration
	}
	p.Stop()

	// average latency in milliseconds
	var averageLatency int64
	averageLatency = totalDuration.Nanoseconds() / 1e6 / int64(pingCount)

	r = &Result{
		Location: c.ID,
		Latency:  averageLatency,
		User:     in.User,
		Pinging:  false,
	}
	return
}
