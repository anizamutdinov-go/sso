package auth

import (
	"context"
	ssov1 "github.com/anizamutdinov-go/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyNumeric = 0
	emptyString  = ""
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if req.GetEmail() == emptyString {
		return nil, status.Error(codes.InvalidArgument, "EMail not found")
	}
	if req.GetPassword() == emptyString {
		return nil, status.Error(codes.InvalidArgument, "Password not found")
	}
	if req.GetAppId() == emptyNumeric {
		return nil, status.Error(codes.InvalidArgument, "ApplicationID not found")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if req.GetEmail() == emptyString {
		return nil, status.Error(codes.InvalidArgument, "EMail not found")
	}
	if req.GetPassword() == emptyString {
		return nil, status.Error(codes.InvalidArgument, "Password not found")
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}
	return &ssov1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	if req.GetUserId() == emptyNumeric {
		return nil, status.Error(codes.InvalidArgument, "UserID not found")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}
	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
