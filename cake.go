// Package cake is a CPace PAKE that's dead simple to use. Using it really is a piece of cake.
package cake

import (
	"github.com/bytemare/cpace"
	"github.com/bytemare/cryptotools"
	"github.com/bytemare/cryptotools/encoding"
	"github.com/bytemare/cryptotools/hash"
	"github.com/bytemare/cryptotools/hashtogroup"
	"github.com/bytemare/cryptotools/ihf"
	"github.com/bytemare/pake"
)

// Initiator represents the initiator role in the CPace protocol.
type Initiator struct {
	c pake.Pake
}

// Responder represents the responder role in the CPace protocol.
type Responder struct {
	c pake.Pake
}

func bake(role pake.Role, serverID, clientID, secret []byte) (cp *cpace.Parameters, ct *cryptotools.Parameters) {
	ct = &cryptotools.Parameters{
		Group:  hashtogroup.Ristretto255Sha512,
		Hash:   hash.SHA512,
		IHF:    ihf.Argon2id,
		IHFLen: ihf.DefaultLength,
	}

	cp = &cpace.Parameters{
		Secret:   secret,
		Encoding: encoding.JSON,
	}

	if role == pake.Initiator {
		cp.ID = clientID
		cp.PeerID = serverID
	} else {
		cp.ID = serverID
		cp.PeerID = clientID
	}

	return cp, ct
}

// Client returns a client instance with the input values and is ready to use.
func Client(serverID, clientID, secret []byte) (*Initiator, error) {
	client, err := cpace.Client(bake(pake.Initiator, serverID, clientID, secret))
	if err != nil {
		return nil, err
	}

	return &Initiator{c: client}, nil
}

// Server returns a server instance with the input values relative to the client and is ready to accept the Client request.
func Server(serverID, clientID, secret []byte) (*Responder, error) {
	server, err := cpace.Server(bake(pake.Responder, serverID, clientID, secret))
	if err != nil {
		return nil, err
	}

	return &Responder{c: server}, nil
}

// Start initiates the protocol with the client. The output message should be send to the peer.
func (i *Initiator) Start() (message []byte, err error) {
	return i.c.Authenticate(nil)
}

// Finish wraps up the protocol given the message from the peer (server)
// On success, generates the session key internally, that can then be retrieved with SessionKey().
// If not, returns the encountered error.
func (i *Initiator) Finish(message []byte) error {
	_, err := i.c.Authenticate(message)
	return err
}

// SessionKey returns the session key if the Finish function returned successfully, or nil otherwise.
func (i *Initiator) SessionKey() []byte {
	return i.c.SessionKey()
}

// AuthenticateClient reads the client's message and completes the server's part of the protocol.
// On success, the output response should be send to the client. If not, returns the encountered error.
func (r *Responder) AuthenticateClient(message []byte) (response []byte, err error) {
	return r.c.Authenticate(message)
}

// SessionKey returns the session key if the AuthenticateClient function returned successfully, or nil otherwise.
func (r *Responder) SessionKey() []byte {
	return r.c.SessionKey()
}
