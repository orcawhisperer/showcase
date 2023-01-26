// Video grpc service package

package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/iamvasanth07/showcase/common"
	pb "github.com/iamvasanth07/showcase/video/proto"
	"github.com/iamvasanth07/showcase/video/repo"
	"google.golang.org/grpc"
)

type IVideoService interface {
	CreateVideo(ctx context.Context, video *pb.Video) (*pb.Video, error)
	GetVideo(ctx context.Context, videoId *pb.VideoId) (*pb.Video, error)
	ListVideos(ctx context.Context, limit *pb.Limit) (*pb.Videos, error)
	UpdateVideo(ctx context.Context, video *pb.Video) (*pb.Video, error)
	DeleteVideo(ctx context.Context, videoId *pb.VideoId) (*pb.Video, error)
}

type VideoService struct {
	db     *repo.VideoRepo
	logger *log.Logger
}

func NewVideoService(db *repo.VideoRepo, logger *log.Logger) *VideoService {
	return &VideoService{db, logger}
}

func (v *VideoService) CreateVideo(ctx context.Context, video *pb.Video) (*pb.Video, error) {
	v.logger.Println("Creating video")
	err := v.db.CreateVideo(video)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (v *VideoService) GetVideo(ctx context.Context, videoId *pb.VideoId) (*pb.Video, error) {
	v.logger.Println("Getting video")
	video, err := v.db.GetVideo(videoId.Id)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (v *VideoService) ListVideos(ctx context.Context, limit *pb.Limit) (*pb.Videos, error) {
	v.logger.Println("Listing videos")
	videos, err := v.db.ListVideos(int(limit.Limit), int(limit.Offset))
	if err != nil {
		return nil, err
	}
	return &pb.Videos{Videos: videos}, nil
}

func (v *VideoService) UpdateVideo(ctx context.Context, video *pb.Video) (*pb.Video, error) {
	v.logger.Println("Updating video")
	err := v.db.UpdateVideo(video)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (v *VideoService) DeleteVideo(ctx context.Context, videoId *pb.VideoId) (*pb.Video, error) {
	v.logger.Println("Deleting video")
	video, err := v.db.GetVideo(videoId.Id)
	if err != nil {
		return nil, err
	}
	err = v.db.DeleteVideo(videoId.Id)
	if err != nil {
		return nil, err
	}
	return video, nil
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
	pb.RegisterVideoServiceServer(srv, s)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("Video service started on port " + os.Getenv("VIDEO_SERVICE_PORT"))
}
