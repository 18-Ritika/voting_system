package voting

import (
	"context"
	"sync"
	"voting_weSockets/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VotingServiceServer struct {
	sessions map[string]*models.Session
	mu       sync.Mutex
}

type CreateSessionRequest struct {
	SessionID string
}

type CreateSessionResponse struct {
	Success bool
}

type JoinSessionRequest struct {
	SessionID string
	UserID    string
}

type JoinSessionResponse struct {
	Success bool
}

type CastVoteRequest struct {
	SessionID string
	UserID    string
	Vote      string
}

type CastVoteResponse struct {
	Success bool
}

func (s *VotingServiceServer) CreateSession(ctx context.Context, req *CreateSessionRequest) (*CreateSessionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[req.SessionID] = &models.Session{
		SessionID: req.SessionID,
		Votes:     make(map[string]string),
		Clients:   []*models.Client{},
	}
	return &CreateSessionResponse{Success: true}, nil
}

func (s *VotingServiceServer) JoinSession(ctx context.Context, req *JoinSessionRequest) (*JoinSessionResponse, error) {
	_, ok := s.sessions[req.SessionID]
	if !ok {
		return &JoinSessionResponse{Success: false}, status.Error(codes.NotFound, "session not found")
	}
	// Add user connection to session clients (handled elsewhere)
	return &JoinSessionResponse{Success: true}, nil
}

func (s *VotingServiceServer) CastVote(ctx context.Context, req *CastVoteRequest) (*CastVoteResponse, error) {
	session, ok := s.sessions[req.SessionID]
	if !ok {
		return &CastVoteResponse{Success: false}, status.Error(codes.NotFound, "session not found")
	}
	s.mu.Lock()
	session.Votes[req.UserID] = req.Vote
	s.mu.Unlock()
	// Broadcast vote update to all clients (handled elsewhere)
	return &CastVoteResponse{Success: true}, nil
}

func RegisterVotingServiceServer(s *grpc.Server) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "voting.VotingService",
		HandlerType: (*VotingServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "CreateSession",
				Handler:    _VotingService_CreateSession_Handler,
			},
			{
				MethodName: "JoinSession",
				Handler:    _VotingService_JoinSession_Handler,
			},
			{
				MethodName: "CastVote",
				Handler:    _VotingService_CastVote_Handler,
			},
		},
	}, &VotingServiceServer{sessions: make(map[string]*models.Session)})
}

func _VotingService_CreateSession_Handler(srv interface{}, ctx context.Context, codec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
	if err := codec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*VotingServiceServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/voting.VotingService/CreateSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*VotingServiceServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_JoinSession_Handler(srv interface{}, ctx context.Context, codec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinSessionRequest)
	if err := codec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*VotingServiceServer).JoinSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/voting.VotingService/JoinSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*VotingServiceServer).JoinSession(ctx, req.(*JoinSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_CastVote_Handler(srv interface{}, ctx context.Context, codec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CastVoteRequest)
	if err := codec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*VotingServiceServer).CastVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/voting.VotingService/CastVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*VotingServiceServer).CastVote(ctx, req.(*CastVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}
