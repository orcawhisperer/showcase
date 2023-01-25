// user grpc service package
package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/iamvasanth07/showcase/user/model"
	pb "github.com/iamvasanth07/showcase/user/proto"
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
	log *log.Logger
	db  *repo.UserRepo
	pb.UnimplementedUserServiceServer
}

func NewUserServer(db *repo.UserRepo) *UserServer {
	return &UserServer{
		log: log.New(os.Stdout, "user-service: ", log.LstdFlags),
		db:  db,
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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	conn, err := utils.GetDBConnection(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	db := repo.NewUserRepo(conn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	lis, err := net.Listen("tcp", os.Getenv("USER_SVC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewUserServer(db))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("server started on port: " + os.Getenv("USER_SVC_PORT"))
}
