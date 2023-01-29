// user grpc service package
package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvasanth07/showcase/common"
	pb "github.com/iamvasanth07/showcase/common/protos/user"
	"github.com/iamvasanth07/showcase/user/config"
	"github.com/iamvasanth07/showcase/user/model"
	"github.com/iamvasanth07/showcase/user/repo"
	"github.com/iamvasanth07/showcase/user/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type IUserService interface {
	Create(context.Context, *pb.CreateUserRequest) (*pb.CreateUserRequest, error)
	Update(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	Delete(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	GetAll(context.Context, *pb.GetAllUserRequest) (*pb.GetAllUserResponse, error)
	Get(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error)
	Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)
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
	user := ProtoToUser(req.User)
	user.Password = req.User.Password
	err := s.db.Create(user)
	if err != nil {
		return nil, err
	}
	getUser, err := s.db.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		User: UserToProto(getUser),
	}, nil
}

func (s *UserServer) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := utils.ValidateUserUpdate(req.User); err != nil {
		return nil, err
	}

	user := ProtoToUser(req.User)
	err := s.db.Update(user)
	if err != nil {
		return nil, err
	}
	getUser, err := s.db.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{
		User: UserToProto(getUser),
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
		User: UserToProto(user),
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
	users, err := s.db.FindAll(req.Paginate.Page, req.Paginate.Limit)
	if err != nil {
		return nil, err
	}
	var res []*pb.User
	var meta *pb.Metadata
	for _, user := range users {
		res = append(res, UserToProto(user))
	}
	meta = &pb.Metadata{
		Total: int32(len(res)),
		Page:  req.Paginate.Page,
		Limit: req.Paginate.Limit,
	}
	return &pb.GetAllUserResponse{Users: res, Metadata: meta}, nil
}

// Login user return token
func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := utils.ValidateUserLogin(req); err != nil {
		return nil, err
	}
	user, err := s.db.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}
	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

// function to generate JWT token with expiry time
func (s *UserServer) generateJWTToken(user *model.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(s.settings.JWT.Expiry)).Unix(),
		Issuer:    user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.settings.JWT.Secret))
}

func RunServer() {

	logger := log.New(os.Stdout, "user-service: ", log.LstdFlags)
	settings := config.GetSettings()
	logger.Println("Initializing user service with settings...")
	logger.Printf("%v, %v, %v", settings.Database, settings.Server, settings.Logger)
	conn, err := initDB(settings)
	logger.Println("Migration database...")
	err = migrateDB(conn)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	db := repo.NewUserRepo(conn)

	// Starting HTTP server for gRPC gateway
	// go runHTTPServer(settings, db, logger)
	// Starting gRPC server
	runGRPCServer(settings, db, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	logger.Println("Stopping the server")

	os.Exit(0)

}

func initDB(settings *config.Settings) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", settings.Database.Host, settings.Database.Port, settings.Database.User, settings.Database.Password, settings.Database.Name, settings.Database.SslMode)
	conn, err := common.GetDBConnection(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return conn, nil
}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}

func runGRPCServer(settings *config.Settings, db *repo.UserRepo, logger *log.Logger) {
	userServer := NewUserServer(db, logger, settings)
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.Server.GrpcHost, settings.Server.GrcpPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	pb.RegisterUserServiceServer(s, userServer)
	logger.Println("GRPC Server started on port: " + settings.Server.GrcpPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
}

// func runHTTPServer(settings *config.Settings, db *repo.UserRepo, logger *log.Logger) {
// 	userServer := NewUserServer(db, logger, settings)
// 	grpcMux := runtime.NewServeMux()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	err := pb.RegisterUserServiceHandlerServer(ctx, grpcMux, userServer)
// 	if err != nil {
// 		log.Fatalf("failed to register the handler to the server: %v", err)
// 	}
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 	})
// 	mux.Handle("/", grpcMux)
// 	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.Server.HTTPHost, settings.Server.HTTPPort))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	logger.Println("HTTP Server started on port: " + settings.Server.HTTPPort)
// 	if err := http.Serve(lis, mux); err != nil {
// 		log.Fatalf("failed to serv http server: %v", err)
// 	}
// }

func GetHTTPHandler() http.Handler {
	logger := log.New(os.Stdout, "user-api-service: ", log.LstdFlags)
	settings := config.GetSettings()
	logger.Println("Settings: ", settings.Database.Host, settings.Database.Port, settings.Database.User, settings.Database.Password, settings.Database.Name, settings.Database.SslMode)
	logger.Println("Initializing user http service with settings...")
	logger.Printf("%v, %v, %v", settings.Database, settings.Server, settings.Logger)
	conn, err := initDB(settings)
	db := repo.NewUserRepo(conn)
	userServer := NewUserServer(db, logger, settings)
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterUserServiceHandlerServer(ctx, grpcMux, userServer)
	if err != nil {
		log.Fatalf("failed to register the handler to the server: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	return mux
}
