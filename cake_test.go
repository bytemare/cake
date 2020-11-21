package cake

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	clientID = []byte("client")
	serverID = []byte("server")
	secret   = []byte("password")
)

func TestFail(t *testing.T) {
	// Initiator ID is too big
	bigClient := []byte(strings.Repeat("a", 1<<16+1))

	if _, err := Client(serverID, bigClient, secret); err == nil {
		t.Error("expected error when client ID is too big")
	} else {
		e := errors.New("CPACE - SETUP : id exceeds authorised length")
		assert.EqualErrorf(t, err, e.Error(), "wrong error on big client ID. Expected %q, got %q", e, err)
	}

	if _, err := Server(serverID, bigClient, secret); err == nil {
		t.Error("expected error when client ID is too big")
	} else {
		e := errors.New("CPACE - SETUP : peer ID exceeds authorised length")
		assert.EqualErrorf(t, err, e.Error(), "wrong error on big client ID. Expected %q, got %q", e, err)
	}

	// Server ID is too big
	bigServer := []byte(strings.Repeat("a", 1<<16+1))

	if _, err := Client(bigServer, clientID, secret); err == nil {
		t.Error("expected error when peer ID is too big")
	} else {
		e := errors.New("CPACE - SETUP : peer ID exceeds authorised length")
		assert.EqualErrorf(t, err, e.Error(), "wrong error on big peer ID. Expected %q, got %q", e, err)
	}

	if _, err := Server(bigServer, clientID, secret); err == nil {
		t.Error("expected error when client ID is too big")
	} else {
		e := errors.New("CPACE - SETUP : id exceeds authorised length")
		assert.EqualErrorf(t, err, e.Error(), "wrong error on big client ID. Expected %q, got %q", e, err)
	}
}
