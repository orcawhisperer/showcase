// packgage repo for the video service

package repo

import (
	"github.com/iamvasanth07/showcase/video/model"

	"gorm.io/gorm"
)

type VideoRepo struct {
	db *gorm.DB
}

func NewVideoRepo(db *gorm.DB) *VideoRepo {
	return &VideoRepo{db}
}

func (v *VideoRepo) CreateVideo(video *model.Video) error {
	return v.db.Create(video).Error
}

func (v *VideoRepo) GetVideo(videoId string) (*model.Video, error) {
	var video model.Video
	err := v.db.First(&video, "id = ?", videoId).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (v *VideoRepo) ListVideos(limit int, offset int) ([]model.Video, error) {
	var videos []model.Video
	err := v.db.Limit(limit).Offset(offset).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (v *VideoRepo) UpdateVideo(video *model.Video) error {
	return v.db.Save(video).Error
}

func (v *VideoRepo) DeleteVideo(videoId string) error {
	return v.db.Delete(&model.Video{}, "id = ?", videoId).Error
}

func (v *VideoRepo) GetVideoBySlug(slug string) (*model.Video, error) {
	var video model.Video
	err := v.db.First(&video, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (v *VideoRepo) GetVideoByChannelId(channelId string) ([]model.Video, error) {
	var videos []model.Video
	err := v.db.Find(&videos, "channel_id = ?", channelId).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
