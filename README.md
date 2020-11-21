# Cake - the CPace AKE

Cake is CPace implemented with the Ristretto group. Easy to use, hard to misuse, it's a piece of cake.

The API is minimal as to allow you to focus on the rest of your application and ensure the highest standards for this.
A simple round trip (2 messages) is enough.

CPace + AKE = Cake

# What is CPace?

tl;dr
> CPace is an authentication protocol that allows two parties that share the same password to authenticate one to another.

It allows for secure mutual authentication with plaintext passwords without transmitting them in clear over the wire.
The protocol spits out a very strong shared session secret on success. This secret can be used to derive encryption keys for your session, a session token, or whatever.

# Why Cake?

Authentication and crypto can be hard sometimes. This aims at giving you something as easy and fool-proof as possible to do so.
No hassle with complex configurations, no need to understand the underlying cryptography.

# Gimme an example - How am I supposed to hold this ?

Let's say you have a client program, that we'll call the initiator, and want to authenticate to a peer that already has the clear-text password and knows the client's identity.

```Go
// Both parties have these values
var (
	clientID = []byte("client")
	serverID = []byte("server")
	secret   = []byte("password")
)
```

On the client, do something like this:

```Go
import "github.com/bytemare/cake"

client, err := cake.Client(serverID, username, secret)
	if err != nil {
		panic(err)
	}

message1, err := client.Start()
	if err != nil {
		panic(err)
	}
````

Send this ```message1``` to the peer. Depending on your setup and needs, you can also send the initiator's ID.

The peer, let's call it the responder (or server), receives the authentication request from the client. Note that if you have a database with users and passwords, the lookup is up to you, and not covered in Cake.

```Go
import "github.com/bytemare/cake"

server, err := cake.Server(serverID, username, secret)
	if err != nil {
		panic(err)
	}

message2, err := server.AuthenticateClient(message1)
	if err != nil {
		panic(err)
	}

sessionKey := server.SessionKey()
```

Send ```message2``` back to the initiator. Note that the server can already extract the session key !

The client needs to ingest this last message to complete the implicit authentication of the server and get the session key.

```Go
err := client.Finish(message2)
	if err != nil {
		panic(err)
	}

sessionKey := client.SessionKey()
```

The ```sessionKey``` is the same for both peers, and can be used for whatever secret you need.

# Under the hood

You don't need to understand the following to use Cake. But if you're hungry, continue reading !

1. Cake is a wrapper to the [CPace](https://github.com/bytemare/cpace) implementation. It's a fancy Diffie Hellman throwing the password onto an elliptic curve.
1. JSON encoded [PAKE messages](https://github.com/bytemare/pake)
1. Strong cryptographic defaults with Ristretto255, SHA512, and Argon2id.
1. Uses [HashToGroup](https://github.com/bytemare/cryptotools) to map the secret to the Ristretto group.

### Work in progress

- Even more testing, and fuzzing.
- Compilable to wasm.

# Third party compatibility

If you want to implement another client or responder compatible with this, you'll need to be able to parse the exchanged
messages, which have the same JSON encoded format. You can read about them [here](https://github.com/bytemare/pake).

The initiator's message must have the field set to the SID, if not, you must give it to your responder in another way (not covered in cake, see CPace).
The responder's message has auth set to nil, and is ignored by the initiator.

# Important note

This is possible thanks to the tremendous and relentless work of people who want greater security for the people, and giving their knowledge to the community.
Special thanks to
- [CPace](https://datatracker.ietf.org/doc/draft-irtf-cfrg-cpace) : [Björn Haase](https://github.com/BjoernMHaase)
- The [Ristretto](https://datatracker.ietf.org/doc/draft-irtf-cfrg-ristretto255-decaf448) team
- The [Hash-to-curve](https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve) team

Many others.
