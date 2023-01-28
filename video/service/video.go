// Video grpc service package

package service

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvasanth07/showcase/common"
	pb "github.com/iamvasanth07/showcase/common/protos/video"
	"github.com/iamvasanth07/showcase/video/config"
	"github.com/iamvasanth07/showcase/video/model"
	"github.com/iamvasanth07/showcase/video/repo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type IVideoService interface {
	CreateVideo(ctx context.Context, req *pb.CreateVideoRequest) (*pb.CreateVideoResponse, error)
	GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error)
	ListVideos(ctx context.Context, req *pb.ListVideosRequest) (*pb.ListVideosResponse, error)
	UpdateVideo(ctx context.Context, req *pb.UpdateVideoRequest) (*pb.UpdateVideoResponse, error)
	DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error)
}

type VideoServer struct {
	db       *repo.VideoRepo
	log      *log.Logger
	settings *config.Settings
	pb.UnimplementedVideoServiceServer
}

func NewVideoService(db *repo.VideoRepo, logger *log.Logger, settings *config.Settings) *VideoServer {
	return &VideoServer{
		db:       db,
		log:      logger,
		settings: settings,
	}
}

func (s *VideoServer) CreateVideo(ctx context.Context, req *pb.CreateVideoRequest) (*pb.CreateVideoResponse, error) {
	s.log.Println("Create video request received")
	video := model.Video{
		Title:       req.Video.Title,
		Description: req.Video.Description,
		Category:    req.Video.Category,
	}
	err := s.db.CreateVideo(&video)
	if err != nil {
		return nil, err
	}
	res := &pb.CreateVideoResponse{
		Video: &pb.Video{
			Id:          video.Id,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoServer) GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error) {
	s.log.Println("Get video request received")
	video, err := s.db.GetVideo(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.GetVideoResponse{
		Video: &pb.Video{
			Id:          video.Id,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoServer) ListVideos(ctx context.Context, req *pb.ListVideosRequest) (*pb.ListVideosResponse, error) {
	s.log.Println("List videos request received")
	videos, err := s.db.ListVideos(int(req.Page), int(req.Limit))
	if err != nil {
		return nil, err
	}
	var pbVideos []*pb.Video
	for _, video := range videos {
		pbVideos = append(pbVideos, &pb.Video{
			Id:          video.Id,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		})
	}
	res := &pb.ListVideosResponse{
		Videos: pbVideos,
	}
	return res, nil
}

func (s *VideoServer) UpdateVideo(ctx context.Context, req *pb.UpdateVideoRequest) (*pb.UpdateVideoResponse, error) {
	s.log.Println("Update video request received")
	video := model.Video{
		Id:          req.Video.Id,
		Title:       req.Video.Title,
		Description: req.Video.Description,
		Url:         req.Video.Url,
	}
	err := s.db.UpdateVideo(&video)
	if err != nil {
		return nil, err
	}
	res := &pb.UpdateVideoResponse{
		Video: &pb.Video{
			Id:          video.Id,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoServer) DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error) {
	s.log.Println("Delete video request received")
	err := s.db.DeleteVideo(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.DeleteVideoResponse{
		Video: &pb.Video{
			Id: req.Id,
		},
	}
	return res, nil
}

//func RunServer() {
//	logger := log.New(os.Stdout, "video-service ", log.LstdFlags)
//	settings := config.GetSettings()
//	logger.Println("Initializing user service with settings...")
//	logger.Printf("%v, %v, %v", settings.Database, settings.Server, settings.Logger)
//	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", settings.Database.Host, settings.Database.Port, settings.Database.User, settings.Database.Password, settings.Database.Name, settings.Database.SslMode)
//	conn, err := common.GetDBConnection(dsn)
//	if err != nil {
//		log.Fatalf("failed to connect to database: %v", err)
//	}
//	go func() {
//		repo := repo.NewVideoRepo(conn)
//		s := NewVideoService(repo, logger, settings)
//		s.log.Println("Video service started on port " + settings.Server.Port)
//		lis, err := net.Listen("tcp", "localhost:"+settings.Server.Port)
//		if err != nil {
//			log.Fatalf("failed to listen: %v", err)
//		}
//		srv := grpc.NewServer()
//		pb.RegisterVideoServiceServer(srv, s)
//		if err := srv.Serve(lis); err != nil {
//			log.Fatalf("failed to serve: %v", err)
//		}
//	}()
//
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt)
//
//	// Block until a signal is received.
//	<-c
//	logger.Println("Shutting down video service...")
//	os.Exit(0)
//
//}

func RunServer() {

	logger := log.New(os.Stdout, "video-service: ", log.LstdFlags)
	settings := config.GetSettings()
	logger.Println("Initializing video service with settings...")
	logger.Printf("%v, %v, %v", settings.Database, settings.Server, settings.Logger)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", settings.Database.Host, settings.Database.Port, settings.Database.User, settings.Database.Password, settings.Database.Name, settings.Database.SslMode)
	conn, err := common.GetDBConnection(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	logger.Println("Migration database...")
	err = migrateDB(conn)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	db := repo.NewVideoRepo(conn)

	// Starting HTTP server for gRPC gateway
	go runHTTPServer(settings, db, logger)
	// Starting gRPC server
	runGRPCServer(settings, db, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	logger.Println("Stopping the server")

	os.Exit(0)

}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Video{},
	)
}

func runGRPCServer(settings *config.Settings, db *repo.VideoRepo, logger *log.Logger) {
	videoServer := NewVideoService(db, logger, settings)
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.Server.GrpcHost, settings.Server.GrcpPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	pb.RegisterVideoServiceServer(s, videoServer)
	logger.Println("GRPC Server started on port: " + settings.Server.GrcpPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
}

func runHTTPServer(settings *config.Settings, db *repo.VideoRepo, logger *log.Logger) {
	videoServer := NewVideoService(db, logger, settings)
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterVideoServiceHandlerServer(ctx, grpcMux, videoServer)
	if err != nil {
		log.Fatalf("failed to register the handler to the server: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("/", grpcMux)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.Server.HTTPHost, settings.Server.HTTPPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger.Println("HTTP Server started on port: " + settings.Server.HTTPPort)
	if err := http.Serve(lis, mux); err != nil {
		log.Fatalf("failed to serv http server: %v", err)
	}
}
