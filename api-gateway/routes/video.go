// routes for video microservice

package routes

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/iamvasanth07/showcase/common/protos/video"
	"github.com/iamvasanth07/showcase/video/config"
	"google.golang.org/grpc"
)

// VideoRoutes struct
type VideoRoutes struct {
	videoClient pb.VideoServiceClient
	config      *config.Settings
}

// NewVideoRoutes returns a new video routes
func NewVideoRoutes(config *config.Settings) *VideoRoutes {
	log.Println(config.Server)
	client, err := grpc.Dial(fmt.Sprintf("%s:%s", config.Server.GrpcHost, config.Server.GrcpPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &VideoRoutes{
		config:      config,
		videoClient: pb.NewVideoServiceClient(client),
	}
}

// RegisterRoutes registers the video routes
func (r *VideoRoutes) RegisterVideoSvcRoutes(router *gin.Engine) {
	routes := router.Group("/api/v1")
	routes.GET("/videos", r.GetVideos)
	routes.GET("/videos/:slug", r.GetVideo)
	routes.POST("/videos", r.CreateVideo)
	routes.PUT("/videos/:slug", r.UpdateVideo)
	routes.DELETE("/videos/:slug", r.DeleteVideo)
}

// GetVideos returns all the videos
func (r *VideoRoutes) GetVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	videos, err := r.videoClient.ListVideos(c, &pb.ListVideosRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(200, videos)
}

// GetVideo returns a video
func (r *VideoRoutes) GetVideo(c *gin.Context) {
	slug := c.Param("slug")
	video, err := r.videoClient.GetVideo(c, &pb.GetVideoRequest{
		Slug: slug,
	})
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(200, video)
}

// CreateVideo creates a video
func (r *VideoRoutes) CreateVideo(c *gin.Context) {
	body := &pb.CreateVideoRequest{}
	if err := c.BindJSON(body); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	video, err := r.videoClient.CreateVideo(c, body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(200, video)

}

// UpdateVideo updates a video
func (r *VideoRoutes) UpdateVideo(c *gin.Context) {
	slug := c.Param("slug")
	body := &pb.UpdateVideoRequest{
		Video: &pb.Video{
			Slug: slug,
		},
	}
	if err := c.BindJSON(body); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	video, err := r.videoClient.UpdateVideo(c, body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(200, video)
}

// DeleteVideo deletes a video
func (r *VideoRoutes) DeleteVideo(c *gin.Context) {
	slug := c.Param("slug")
	_, err := r.videoClient.DeleteVideo(c, &pb.DeleteVideoRequest{
		Slug: slug,
	})
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"message": "Video deleted successfully",
	})
}
