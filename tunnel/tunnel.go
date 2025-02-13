// Package tunnel provides gre network tunnelling
package tunnel

import (
	"errors"
	"time"

	"github.com/micro/go-micro/transport"
)

const (
	// send over one link
	Unicast Mode = iota
	// send to all channel listeners
	Multicast
	// send to all links
	Broadcast
)

var (
	// DefaultDialTimeout is the dial timeout if none is specified
	DefaultDialTimeout = time.Second * 5
	// ErrDialTimeout is returned by a call to Dial where the timeout occurs
	ErrDialTimeout = errors.New("dial timeout")
	// ErrDiscoverChan is returned when we failed to receive the "announce" back from a discovery
	ErrDiscoverChan = errors.New("failed to discover channel")
	// ErrLinkNotFound is returned when a link is specified at dial time and does not exist
	ErrLinkNotFound = errors.New("link not found")
)

// Mode of the session
type Mode uint8

// Tunnel creates a gre tunnel on top of the go-micro/transport.
// It establishes multiple streams using the Micro-Tunnel-Channel header
// and Micro-Tunnel-Session header. The tunnel id is a hash of
// the address being requested.
type Tunnel interface {
	// Init initializes tunnel with options
	Init(opts ...Option) error
	// Address returns the address the tunnel is listening on
	Address() string
	// Connect connects the tunnel
	Connect() error
	// Close closes the tunnel
	Close() error
	// Links returns all the links the tunnel is connected to
	Links() []Link
	// Dial allows a client to connect to a channel
	Dial(channel string, opts ...DialOption) (Session, error)
	// Listen allows to accept connections on a channel
	Listen(channel string, opts ...ListenOption) (Listener, error)
	// String returns the name of the tunnel implementation
	String() string
}

// Link represents internal links to the tunnel
type Link interface {
	// Id returns the link unique Id
	Id() string
	// Delay is the current load on the link (lower is better)
	Delay() int64
	// Length returns the roundtrip time as nanoseconds (lower is better)
	Length() int64
	// Current transfer rate as bits per second (lower is better)
	Rate() float64
	// State of the link e.g connected/closed
	State() string
	// honours transport socket
	transport.Socket
}

// The listener provides similar constructs to the transport.Listener
type Listener interface {
	Accept() (Session, error)
	Channel() string
	Close() error
}

// Session is a unique session created when dialling or accepting connections on the tunnel
type Session interface {
	// The unique session id
	Id() string
	// The channel name
	Channel() string
	// The link the session is on
	Link() string
	// a transport socket
	transport.Socket
}

// NewTunnel creates a new tunnel
func NewTunnel(opts ...Option) Tunnel {
	return newTunnel(opts...)
}
