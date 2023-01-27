// user grpc service package
package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/iamvasanth07/showcase/common"
	pb "github.com/iamvasanth07/showcase/common/protos"
	"github.com/iamvasanth07/showcase/user/config"
	"github.com/iamvasanth07/showcase/user/model"
	"github.com/iamvasanth07/showcase/user/repo"
	"github.com/iamvasanth07/showcase/user/utils"
	"google.golang.org/grpc"
)

type IUserService interface {
	Create(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	Update(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	Delete(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	GetAll(context.Context, *pb.GetAllUserRequest) (*pb.GetAllUserResponse, error)
	Get(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error)
}

type UserServer struct {
	db       *repo.UserRepo
	log      *log.Logger
	settings *config.Settings
	pb.UnimplementedUserServiceServer
}

func NewUserServer(db *repo.UserRepo, log *log.Logger, settings *config.Settings) *UserServer {
	return &UserServer{
		db:       db,
		log:      log,
		settings: settings,
	}
}

func (s *UserServer) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := utils.ValidateUserCreate(req.User); err != nil {
		return nil, err
	}
	user := &model.User{
		Name:  req.User.Name,
		Email: req.User.Email,
		Phone: req.User.Phone,
	}
	err := s.db.Create(user)
	if err != nil {
		return nil, err
	}
	getUser, err := s.db.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    getUser.Id,
			Name:  getUser.Name,
			Email: getUser.Email,
			Phone: getUser.Phone,
		},
	}, nil
}

func (s *UserServer) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := utils.ValidateUserUpdate(req.User); err != nil {
		return nil, err
	}

	user := &model.User{
		Name:  req.User.Name,
		Email: req.User.Email,
		Phone: req.User.Phone,
	}
	err := s.db.Update(user)
	if err != nil {
		return nil, err
	}
	getUser, err := s.db.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:    getUser.Id,
			Name:  getUser.Name,
			Email: getUser.Email,
			Phone: getUser.Phone,
		},
	}, nil
}

func (s *UserServer) Get(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	if err := utils.ValidateUserGet(req.Id); err != nil {
		return nil, err
	}
	user, err := s.db.FindByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
	}, nil
}

func (s *UserServer) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := utils.ValidateUserDelete(req.Id); err != nil {
		return nil, err
	}
	err := s.db.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{
		Id: req.Id,
	}, nil
}

func (s *UserServer) GetAll(ctx context.Context, req *pb.GetAllUserRequest) (*pb.GetAllUserResponse, error) {

	if err := utils.ValidateUserGetAll(req); err != nil {
		return nil, err
	}
	users, err := s.db.FindAll(req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		return nil, err
	}
	var res []*pb.User
	var meta *pb.Metadata
	for _, user := range users {
		res = append(res, &pb.User{
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		})
	}
	meta = &pb.Metadata{
		Total: int32(len(res)),
		Page:  req.Pagination.Page,
		Limit: req.Pagination.Limit,
	}
	return &pb.GetAllUserResponse{Users: res, Metadata: meta}, nil
}

func RunServer() {

	logger := log.New(os.Stdout, "user-service: ", log.LstdFlags)
	settings := config.GetSettings()
	logger.Println("Initializing user service with settings...")
	logger.Printf("%v, %v, %v", settings.Database, settings.Server, settings.Logger)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", settings.Database.Host, settings.Database.Port, settings.Database.User, settings.Database.Password, settings.Database.Name, settings.Database.SslMode)
	conn, err := common.GetDBConnection(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	go func() {
		db := repo.NewUserRepo(conn)
		userServer := NewUserServer(db, logger, settings)
		userServer.log.Println("Server started on port: " + settings.Server.Port)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", settings.Server.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		opts := []grpc.ServerOption{}
		s := grpc.NewServer(opts...)
		pb.RegisterUserServiceServer(s, userServer)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	logger.Println("Stopping the server")

	os.Exit(0)

}
