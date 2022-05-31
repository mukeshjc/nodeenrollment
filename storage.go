package nodeenrollment

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// MessageWithId is a proto message that is required to implement a GetId
// function, which will be immediately satisfied by any message with an
// `string id = X;` parameter.
type MessageWithId interface {
	proto.Message
	GetId() string
}

// Storage is an interface for to store values. The interface operates on
// proto.Message or MessageWithId (which embeds a proto.Message but requires a
// GetId() function), which is satisifed by all types in this library and
// provides some type safety vs. any.
//
// The interface can be used for multiple types of message via a type switch on,
// allowing various implementations to then read or write the correct data from
// e.g. separate storage locations.
type Storage interface {
	// Store stores the given message
	Store(context.Context, MessageWithId) error

	// Load loads values into the given message. The message must be populated
	// with the ID of the value to load. If not found, the returned error should
	// be ErrNotFound.
	Load(context.Context, MessageWithId) error

	// Remove removes the given message. The ID field of the message must be
	// populated, and only the ID field of the message is considered.
	Remove(context.Context, MessageWithId) error

	// List returns a list of IDs; the type of the message is used to
	// disambiguate what to list, and can be a nil pointer to the type.
	List(context.Context, proto.Message) ([]string, error)
}
