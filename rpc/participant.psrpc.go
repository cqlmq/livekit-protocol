// Code generated by protoc-gen-psrpc v0.6.0, DO NOT EDIT.
// source: rpc/participant.proto

package rpc

import (
	"context"

	"github.com/livekit/psrpc"
	"github.com/livekit/psrpc/pkg/client"
	"github.com/livekit/psrpc/pkg/info"
	"github.com/livekit/psrpc/pkg/rand"
	"github.com/livekit/psrpc/pkg/server"
	"github.com/livekit/psrpc/version"
)
import livekit1 "github.com/livekit/protocol/livekit"
import livekit6 "github.com/livekit/protocol/livekit"

var _ = version.PsrpcVersion_0_6

// ============================
// Participant Client Interface
// ============================

type ParticipantClient[ParticipantTopicType ~string] interface {
	RemoveParticipant(ctx context.Context, participant ParticipantTopicType, req *livekit6.RoomParticipantIdentity, opts ...psrpc.RequestOption) (*livekit6.RemoveParticipantResponse, error)

	MutePublishedTrack(ctx context.Context, participant ParticipantTopicType, req *livekit6.MuteRoomTrackRequest, opts ...psrpc.RequestOption) (*livekit6.MuteRoomTrackResponse, error)

	UpdateParticipant(ctx context.Context, participant ParticipantTopicType, req *livekit6.UpdateParticipantRequest, opts ...psrpc.RequestOption) (*livekit1.ParticipantInfo, error)

	UpdateSubscriptions(ctx context.Context, participant ParticipantTopicType, req *livekit6.UpdateSubscriptionsRequest, opts ...psrpc.RequestOption) (*livekit6.UpdateSubscriptionsResponse, error)

	ForwardParticipant(ctx context.Context, participant ParticipantTopicType, req *livekit6.ForwardParticipantRequest, opts ...psrpc.RequestOption) (*livekit6.ForwardParticipantResponse, error)

	// Close immediately, without waiting for pending RPCs
	Close()
}

// ================================
// Participant ServerImpl Interface
// ================================

type ParticipantServerImpl interface {
	RemoveParticipant(context.Context, *livekit6.RoomParticipantIdentity) (*livekit6.RemoveParticipantResponse, error)

	MutePublishedTrack(context.Context, *livekit6.MuteRoomTrackRequest) (*livekit6.MuteRoomTrackResponse, error)

	UpdateParticipant(context.Context, *livekit6.UpdateParticipantRequest) (*livekit1.ParticipantInfo, error)

	UpdateSubscriptions(context.Context, *livekit6.UpdateSubscriptionsRequest) (*livekit6.UpdateSubscriptionsResponse, error)

	ForwardParticipant(context.Context, *livekit6.ForwardParticipantRequest) (*livekit6.ForwardParticipantResponse, error)
}

// ============================
// Participant Server Interface
// ============================

type ParticipantServer[ParticipantTopicType ~string] interface {
	RegisterRemoveParticipantTopic(participant ParticipantTopicType) error
	DeregisterRemoveParticipantTopic(participant ParticipantTopicType)
	RegisterMutePublishedTrackTopic(participant ParticipantTopicType) error
	DeregisterMutePublishedTrackTopic(participant ParticipantTopicType)
	RegisterUpdateParticipantTopic(participant ParticipantTopicType) error
	DeregisterUpdateParticipantTopic(participant ParticipantTopicType)
	RegisterUpdateSubscriptionsTopic(participant ParticipantTopicType) error
	DeregisterUpdateSubscriptionsTopic(participant ParticipantTopicType)
	RegisterForwardParticipantTopic(participant ParticipantTopicType) error
	DeregisterForwardParticipantTopic(participant ParticipantTopicType)
	RegisterAllParticipantTopics(participant ParticipantTopicType) error
	DeregisterAllParticipantTopics(participant ParticipantTopicType)

	// Close and wait for pending RPCs to complete
	Shutdown()

	// Close immediately, without waiting for pending RPCs
	Kill()
}

// ==================
// Participant Client
// ==================

type participantClient[ParticipantTopicType ~string] struct {
	client *client.RPCClient
}

// NewParticipantClient creates a psrpc client that implements the ParticipantClient interface.
func NewParticipantClient[ParticipantTopicType ~string](bus psrpc.MessageBus, opts ...psrpc.ClientOption) (ParticipantClient[ParticipantTopicType], error) {
	sd := &info.ServiceDefinition{
		Name: "Participant",
		ID:   rand.NewClientID(),
	}

	sd.RegisterMethod("RemoveParticipant", false, false, true, true)
	sd.RegisterMethod("MutePublishedTrack", false, false, true, true)
	sd.RegisterMethod("UpdateParticipant", false, false, true, true)
	sd.RegisterMethod("UpdateSubscriptions", false, false, true, true)
	sd.RegisterMethod("ForwardParticipant", false, false, true, true)

	rpcClient, err := client.NewRPCClient(sd, bus, opts...)
	if err != nil {
		return nil, err
	}

	return &participantClient[ParticipantTopicType]{
		client: rpcClient,
	}, nil
}

func (c *participantClient[ParticipantTopicType]) RemoveParticipant(ctx context.Context, participant ParticipantTopicType, req *livekit6.RoomParticipantIdentity, opts ...psrpc.RequestOption) (*livekit6.RemoveParticipantResponse, error) {
	return client.RequestSingle[*livekit6.RemoveParticipantResponse](ctx, c.client, "RemoveParticipant", []string{string(participant)}, req, opts...)
}

func (c *participantClient[ParticipantTopicType]) MutePublishedTrack(ctx context.Context, participant ParticipantTopicType, req *livekit6.MuteRoomTrackRequest, opts ...psrpc.RequestOption) (*livekit6.MuteRoomTrackResponse, error) {
	return client.RequestSingle[*livekit6.MuteRoomTrackResponse](ctx, c.client, "MutePublishedTrack", []string{string(participant)}, req, opts...)
}

func (c *participantClient[ParticipantTopicType]) UpdateParticipant(ctx context.Context, participant ParticipantTopicType, req *livekit6.UpdateParticipantRequest, opts ...psrpc.RequestOption) (*livekit1.ParticipantInfo, error) {
	return client.RequestSingle[*livekit1.ParticipantInfo](ctx, c.client, "UpdateParticipant", []string{string(participant)}, req, opts...)
}

func (c *participantClient[ParticipantTopicType]) UpdateSubscriptions(ctx context.Context, participant ParticipantTopicType, req *livekit6.UpdateSubscriptionsRequest, opts ...psrpc.RequestOption) (*livekit6.UpdateSubscriptionsResponse, error) {
	return client.RequestSingle[*livekit6.UpdateSubscriptionsResponse](ctx, c.client, "UpdateSubscriptions", []string{string(participant)}, req, opts...)
}

func (c *participantClient[ParticipantTopicType]) ForwardParticipant(ctx context.Context, participant ParticipantTopicType, req *livekit6.ForwardParticipantRequest, opts ...psrpc.RequestOption) (*livekit6.ForwardParticipantResponse, error) {
	return client.RequestSingle[*livekit6.ForwardParticipantResponse](ctx, c.client, "ForwardParticipant", []string{string(participant)}, req, opts...)
}

func (s *participantClient[ParticipantTopicType]) Close() {
	s.client.Close()
}

// ==================
// Participant Server
// ==================

type participantServer[ParticipantTopicType ~string] struct {
	svc ParticipantServerImpl
	rpc *server.RPCServer
}

// NewParticipantServer builds a RPCServer that will route requests
// to the corresponding method in the provided svc implementation.
func NewParticipantServer[ParticipantTopicType ~string](svc ParticipantServerImpl, bus psrpc.MessageBus, opts ...psrpc.ServerOption) (ParticipantServer[ParticipantTopicType], error) {
	sd := &info.ServiceDefinition{
		Name: "Participant",
		ID:   rand.NewServerID(),
	}

	s := server.NewRPCServer(sd, bus, opts...)

	sd.RegisterMethod("RemoveParticipant", false, false, true, true)
	sd.RegisterMethod("MutePublishedTrack", false, false, true, true)
	sd.RegisterMethod("UpdateParticipant", false, false, true, true)
	sd.RegisterMethod("UpdateSubscriptions", false, false, true, true)
	sd.RegisterMethod("ForwardParticipant", false, false, true, true)
	return &participantServer[ParticipantTopicType]{
		svc: svc,
		rpc: s,
	}, nil
}

func (s *participantServer[ParticipantTopicType]) RegisterRemoveParticipantTopic(participant ParticipantTopicType) error {
	return server.RegisterHandler(s.rpc, "RemoveParticipant", []string{string(participant)}, s.svc.RemoveParticipant, nil)
}

func (s *participantServer[ParticipantTopicType]) DeregisterRemoveParticipantTopic(participant ParticipantTopicType) {
	s.rpc.DeregisterHandler("RemoveParticipant", []string{string(participant)})
}

func (s *participantServer[ParticipantTopicType]) RegisterMutePublishedTrackTopic(participant ParticipantTopicType) error {
	return server.RegisterHandler(s.rpc, "MutePublishedTrack", []string{string(participant)}, s.svc.MutePublishedTrack, nil)
}

func (s *participantServer[ParticipantTopicType]) DeregisterMutePublishedTrackTopic(participant ParticipantTopicType) {
	s.rpc.DeregisterHandler("MutePublishedTrack", []string{string(participant)})
}

func (s *participantServer[ParticipantTopicType]) RegisterUpdateParticipantTopic(participant ParticipantTopicType) error {
	return server.RegisterHandler(s.rpc, "UpdateParticipant", []string{string(participant)}, s.svc.UpdateParticipant, nil)
}

func (s *participantServer[ParticipantTopicType]) DeregisterUpdateParticipantTopic(participant ParticipantTopicType) {
	s.rpc.DeregisterHandler("UpdateParticipant", []string{string(participant)})
}

func (s *participantServer[ParticipantTopicType]) RegisterUpdateSubscriptionsTopic(participant ParticipantTopicType) error {
	return server.RegisterHandler(s.rpc, "UpdateSubscriptions", []string{string(participant)}, s.svc.UpdateSubscriptions, nil)
}

func (s *participantServer[ParticipantTopicType]) DeregisterUpdateSubscriptionsTopic(participant ParticipantTopicType) {
	s.rpc.DeregisterHandler("UpdateSubscriptions", []string{string(participant)})
}

func (s *participantServer[ParticipantTopicType]) RegisterForwardParticipantTopic(participant ParticipantTopicType) error {
	return server.RegisterHandler(s.rpc, "ForwardParticipant", []string{string(participant)}, s.svc.ForwardParticipant, nil)
}

func (s *participantServer[ParticipantTopicType]) DeregisterForwardParticipantTopic(participant ParticipantTopicType) {
	s.rpc.DeregisterHandler("ForwardParticipant", []string{string(participant)})
}

func (s *participantServer[ParticipantTopicType]) allParticipantTopicRegisterers() server.RegistererSlice {
	return server.RegistererSlice{
		server.NewRegisterer(s.RegisterRemoveParticipantTopic, s.DeregisterRemoveParticipantTopic),
		server.NewRegisterer(s.RegisterMutePublishedTrackTopic, s.DeregisterMutePublishedTrackTopic),
		server.NewRegisterer(s.RegisterUpdateParticipantTopic, s.DeregisterUpdateParticipantTopic),
		server.NewRegisterer(s.RegisterUpdateSubscriptionsTopic, s.DeregisterUpdateSubscriptionsTopic),
		server.NewRegisterer(s.RegisterForwardParticipantTopic, s.DeregisterForwardParticipantTopic),
	}
}

func (s *participantServer[ParticipantTopicType]) RegisterAllParticipantTopics(participant ParticipantTopicType) error {
	return s.allParticipantTopicRegisterers().Register(participant)
}

func (s *participantServer[ParticipantTopicType]) DeregisterAllParticipantTopics(participant ParticipantTopicType) {
	s.allParticipantTopicRegisterers().Deregister(participant)
}

func (s *participantServer[ParticipantTopicType]) Shutdown() {
	s.rpc.Close(false)
}

func (s *participantServer[ParticipantTopicType]) Kill() {
	s.rpc.Close(true)
}

var psrpcFileDescriptor6 = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xcb, 0x4a, 0xc3, 0x40,
	0x14, 0x25, 0x88, 0x2e, 0xa6, 0x08, 0x76, 0x54, 0x28, 0xc1, 0x47, 0x5f, 0xeb, 0x04, 0xf4, 0x0f,
	0x5c, 0x08, 0x2e, 0x84, 0x12, 0x75, 0xe3, 0x46, 0x92, 0xc9, 0xd5, 0x0e, 0x4d, 0x72, 0xc7, 0x99,
	0x9b, 0x4a, 0x57, 0x2e, 0x04, 0xc1, 0x9d, 0xdf, 0xe2, 0x17, 0x4a, 0x9b, 0x26, 0x99, 0xb6, 0x58,
	0xcc, 0x72, 0xce, 0x39, 0xf7, 0x9c, 0xfb, 0x18, 0x76, 0xac, 0x95, 0xf0, 0x55, 0xa8, 0x49, 0x0a,
	0xa9, 0xc2, 0x8c, 0x3c, 0xa5, 0x91, 0x90, 0xef, 0x68, 0x25, 0xdc, 0x7d, 0x54, 0x24, 0x31, 0x33,
	0x05, 0xe6, 0x1e, 0x25, 0x72, 0x0a, 0x13, 0x49, 0x4f, 0x29, 0xc6, 0x90, 0x94, 0x28, 0x2f, 0x51,
	0x8d, 0x98, 0x16, 0xd8, 0xc5, 0xf7, 0x2e, 0x6b, 0x8d, 0x6a, 0x4f, 0xfe, 0xce, 0xda, 0x01, 0xa4,
	0x38, 0x05, 0x1b, 0xec, 0x7a, 0xcb, 0x4a, 0x2f, 0x40, 0x4c, 0x2d, 0xe6, 0x26, 0x86, 0x8c, 0x24,
	0xcd, 0xdc, 0x7e, 0xad, 0x58, 0xaf, 0x0e, 0xc0, 0x28, 0xcc, 0x0c, 0xf4, 0x87, 0x3f, 0x5f, 0x4e,
	0xf7, 0xc0, 0x71, 0x4f, 0x58, 0xcb, 0x9a, 0x82, 0xdb, 0x8f, 0x8e, 0xc3, 0x67, 0x8c, 0xdf, 0xe6,
	0x04, 0xa3, 0x3c, 0x4a, 0xa4, 0x19, 0x43, 0x7c, 0xaf, 0x43, 0x31, 0xe1, 0xa7, 0x95, 0xff, 0x9c,
	0x9c, 0x77, 0xb1, 0xc0, 0x03, 0x78, 0xcd, 0xc1, 0x90, 0x7b, 0xf6, 0x17, 0xdd, 0x28, 0x7a, 0xca,
	0xda, 0x0f, 0x2a, 0x0e, 0x69, 0x65, 0xf6, 0x5e, 0x65, 0xbd, 0xc1, 0x95, 0xe9, 0x9d, 0x4a, 0x62,
	0xaf, 0x26, 0x7b, 0xc6, 0x7f, 0xe6, 0x7e, 0x3a, 0xec, 0xb0, 0x30, 0xbf, 0xcb, 0x23, 0x23, 0xb4,
	0x2c, 0x6e, 0xc9, 0x07, 0x6b, 0xd1, 0x2b, 0x6c, 0x19, 0x3e, 0xdc, 0x2e, 0x6a, 0xb4, 0x80, 0x0f,
	0x87, 0xf1, 0x6b, 0xd4, 0x6f, 0xa1, 0x8e, 0xed, 0x15, 0xd4, 0xc7, 0xdd, 0x24, 0xcb, 0x36, 0x06,
	0x5b, 0x35, 0x4d, 0xba, 0xb8, 0xea, 0x3d, 0x9e, 0xbf, 0x48, 0x1a, 0xe7, 0x91, 0x27, 0x30, 0xf5,
	0x97, 0xb6, 0xfe, 0xe2, 0xbb, 0x0a, 0x4c, 0x7c, 0xad, 0x44, 0xb4, 0xb7, 0x78, 0x5d, 0xfe, 0x06,
	0x00, 0x00, 0xff, 0xff, 0x09, 0x26, 0xa2, 0x60, 0x13, 0x03, 0x00, 0x00,
}
