package auth

import (
	"context"
	"time"
	"voting_weSockets/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtKey = []byte("my_secret_key")

type AuthServiceServer struct {
	users map[string]*models.User
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type RegisterRequest struct {
	Username string
	Password string
}

type RegisterResponse struct {
	Success bool
	Message string
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token   string
	Message string
}

type LogoutRequest struct {
	Token string
}

type LogoutResponse struct {
	Success bool
}

func (s *AuthServiceServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error hashing password")
	}
	s.users[req.Username] = &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	return &RegisterResponse{Success: true, Message: "User registered successfully"}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, exists := s.users[req.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return &LoginResponse{Message: "Invalid credentials"}, nil
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: req.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error signing token")
	}
	return &LoginResponse{Token: tokenString, Message: "Login successful"}, nil
}

func (s *AuthServiceServer) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// Invalidate the JWT token (this implementation does not store tokens, so it does nothing)
	return &LogoutResponse{Success: true}, nil
}

func RegisterAuthServiceServer(s *grpc.Server) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "auth.AuthService",
		HandlerType: (*AuthServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Register",
				Handler:    _AuthService_Register_Handler,
			},
			{
				MethodName: "Login",
				Handler:    _AuthService_Login_Handler,
			},
			{
				MethodName: "Logout",
				Handler:    _AuthService_Logout_Handler,
			},
		},
	}, &AuthServiceServer{users: make(map[string]*models.User)})
}

func _AuthService_Register_Handler(srv interface{}, ctx context.Context, codec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := codec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*AuthServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*AuthServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, codec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := codec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*AuthServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Logout_Handler(srv interface{}, ctx context.Context, codec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := codec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*AuthServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*AuthServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}
