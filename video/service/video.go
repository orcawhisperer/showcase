// Video grpc service package

package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/iamvasanth07/showcase/common"
	pb "github.com/iamvasanth07/showcase/common/protos/video"
	"github.com/iamvasanth07/showcase/video/config"
	"github.com/iamvasanth07/showcase/video/model"
	"github.com/iamvasanth07/showcase/video/repo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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
			Id:          video.Uuid,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoServer) GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error) {
	s.log.Println("Get video request received")
	video, err := s.db.GetVideoBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	res := &pb.GetVideoResponse{
		Video: &pb.Video{
			Id:          video.Uuid,
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
			Id:          video.Uuid,
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
		Uuid:        req.Video.Id,
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
			Id:          video.Uuid,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoServer) DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error) {
	s.log.Println("Delete video request received")
	err := s.db.DeleteVideo(req.Slug)
	if err != nil {
		return nil, err
	}
	res := &pb.DeleteVideoResponse{
		Video: &pb.Video{
			Slug: req.Slug,
		},
	}
	return res, nil
}

func RunServer() {

	logger := log.New(os.Stdout, "video-service: ", log.LstdFlags)
	settings := config.GetSettings()
	logger.Println("Initializing video service with settings...")
	logger.Printf("%v, %v, %v", settings.Database, settings.Server, settings.Logger)
	conn, err := initDB(settings)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	logger.Println("Migration database...")
	err = migrateDB(conn)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}
	db := repo.NewVideoRepo(conn)
	runGRPCServer(settings, db, logger)

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

func initDB(settings *config.Settings) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", settings.Database.Host, settings.Database.Port, settings.Database.User, settings.Database.Password, settings.Database.Name, settings.Database.SslMode)
	conn, err := common.GetDBConnection(dsn)
	return conn, err
}
