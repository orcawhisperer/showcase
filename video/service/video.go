// Video grpc service package

package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/iamvasanth07/showcase/common"
	"github.com/iamvasanth07/showcase/video/model"
	pb "github.com/iamvasanth07/showcase/video/proto"
	"github.com/iamvasanth07/showcase/video/repo"
	"google.golang.org/grpc"
)

type IVideoService interface {
	CreateVideo(ctx context.Context, req *pb.CreateVideoRequest) (*pb.CreateVideoResponse, error)
	GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error)
	ListVideos(ctx context.Context, req *pb.ListVideosRequest) (*pb.ListVideosResponse, error)
	UpdateVideo(ctx context.Context, req *pb.UpdateVideoRequest) (*pb.UpdateVideoResponse, error)
	DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error)
}

type VideoService struct {
	db     *repo.VideoRepo
	logger *log.Logger
}

func NewVideoService(db *repo.VideoRepo, logger *log.Logger) *VideoService {
	return &VideoService{db, logger}
}

func (s *VideoService) CreateVideo(ctx context.Context, req *pb.CreateVideoRequest) (*pb.CreateVideoResponse, error) {
	s.logger.Println("Create video request received")
	video := model.Video{
		Title:       req.Video.Title,
		Description: req.Video.Description,
		Url:         req.Video.Url,
	}
	err := s.db.CreateVideo(&video)
	if err != nil {
		return nil, err
	}
	res := &pb.CreateVideoResponse{
		Video: &pb.Video{
			Id:          video.ID,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoService) GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error) {
	s.logger.Println("Get video request received")
	video, err := s.db.GetVideo(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.GetVideoResponse{
		Video: &pb.Video{
			Id:          video.ID,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoService) ListVideos(ctx context.Context, req *pb.ListVideosRequest) (*pb.ListVideosResponse, error) {
	s.logger.Println("List videos request received")
	videos, err := s.db.ListVideos()
	if err != nil {
		return nil, err
	}
	var pbVideos []*pb.Video
	for _, video := range videos {
		pbVideos = append(pbVideos, &pb.Video{
			Id:          video.ID,
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

func (s *VideoService) UpdateVideo(ctx context.Context, req *pb.UpdateVideoRequest) (*pb.UpdateVideoResponse, error) {
	s.logger.Println("Update video request received")
	video := model.Video{
		ID:          req.Video.Id,
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
			Id:          video.ID,
			Title:       video.Title,
			Description: video.Description,
			Url:         video.Url,
		},
	}
	return res, nil
}

func (s *VideoService) DeleteVideo(ctx context.Context, req *pb.DeleteVideoRequest) (*pb.DeleteVideoResponse, error) {
	s.logger.Println("Delete video request received")
	err := s.db.DeleteVideo(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.DeleteVideoResponse{
		Success: true,
	}
	return res, nil
}

func RunServer() {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	conn, err := common.GetDBConnection(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	repo := repo.NewVideoRepo(conn)
	logger := log.New(os.Stdout, "video ", log.LstdFlags)
	s := NewVideoService(repo, logger)
	lis, err := net.Listen("tcp", ":"+os.Getenv("VIDEO_SERVICE_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterVideoServiceServer(srv, IVideoService)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("Video service started on port " + os.Getenv("VIDEO_SERVICE_PORT"))
}
