// Code generated by protoc-gen-grain. DO NOT EDIT.
// versions:
//  protoc-gen-grain v0.1.0
//  protoc           v4.24.3
// source: protos.proto

package main

import (
	errors "errors"
	fmt "fmt"
	actor "github.com/asynkron/protoactor-go/actor"
	cluster "github.com/asynkron/protoactor-go/cluster"
	proto "google.golang.org/protobuf/proto"
	slog "log/slog"
	time "time"
)

var xUserActorFactory func() UserActor

// UserActorFactory produces a UserActor
func UserActorFactory(factory func() UserActor) {
	xUserActorFactory = factory
}

// GetUserActorGrainClient instantiates a new UserActorGrainClient with given Identity
func GetUserActorGrainClient(c *cluster.Cluster, id string) *UserActorGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &UserActorGrainClient{Identity: id, cluster: c}
}

// GetUserActorKind instantiates a new cluster.Kind for UserActor
func GetUserActorKind(opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &UserActorActor{
			Timeout: 60 * time.Second,
		}
	}, opts...)
	kind := cluster.NewKind("UserActor", props)
	return kind
}

// GetUserActorKind instantiates a new cluster.Kind for UserActor
func NewUserActorKind(factory func() UserActor, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	xUserActorFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &UserActorActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind("UserActor", props)
	return kind
}

// UserActor interfaces the services available to the UserActor
type UserActor interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
	Connect(*Empty, cluster.GrainContext) (*Empty, error)
}

// UserActorGrainClient holds the base data for the UserActorGrain
type UserActorGrainClient struct {
	Identity string
	cluster  *cluster.Cluster
}

// Connect requests the execution on to the cluster with CallOptions
func (g *UserActorGrainClient) Connect(r *Empty, opts ...cluster.GrainCallOption) (*Empty, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "UserActor", reqMsg, opts...)
	if err != nil {
		return nil, err
	}
	switch msg := resp.(type) {
	case *cluster.GrainResponse:
		result := &Empty{}
		err = proto.Unmarshal(msg.MessageData, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case *cluster.GrainErrorResponse:
		return nil, errors.New(msg.Err)
	default:
		return nil, errors.New("unknown response")
	}
}

// UserActorActor represents the actor structure
type UserActorActor struct {
	ctx     cluster.GrainContext
	inner   UserActor
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *UserActorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: //pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xUserActorFactory()
		a.inner.Init(a.ctx)

		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}
	case *actor.ReceiveTimeout:
		ctx.Poison(ctx.Self())
	case *actor.Stopped:
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
		case 0:
			req := &Empty{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] Connect(Empty) proto.Unmarshal failed.", slog.Any("error", err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			r0, err := a.inner.Connect(req, a.ctx)
			if err != nil {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			bytes, err := proto.Marshal(r0)
			if err != nil {
				ctx.Logger().Error("[Grain] Connect(Empty) proto.Marshal failed", slog.Any("error", err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			resp := &cluster.GrainResponse{MessageData: bytes}
			ctx.Respond(resp)
		}
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}
