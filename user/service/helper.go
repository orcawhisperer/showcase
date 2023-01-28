package service

import (
	pb "github.com/iamvasanth07/showcase/common/protos/user"
	"github.com/iamvasanth07/showcase/user/model"
)

// user model to user proto
func UserToProto(user *model.User) *pb.User {

	userProto := &pb.User{}

	if user == nil {
		return userProto
	}
	if user.UUID != "" {
		userProto.Id = user.UUID
	}
	if user.FirstName != "" {
		userProto.FirstName = user.FirstName
	}
	if user.LastName != "" {
		userProto.LastName = user.LastName
	}
	if user.Username != "" {
		userProto.Username = user.Username
	}
	if user.Email != "" {
		userProto.Email = user.Email
	}
	if user.Phone != "" {
		userProto.Phone = user.Phone
	}
	return userProto
}

// user proto to user model
func ProtoToUser(user *pb.User) *model.User {

	userModel := &model.User{}

	if user == nil {
		return userModel
	}
	if user.Id != "" {
		userModel.UUID = user.Id
	}
	if user.FirstName != "" {
		userModel.FirstName = user.FirstName
	}
	if user.LastName != "" {
		userModel.LastName = user.LastName
	}
	if user.Username != "" {
		userModel.Username = user.Username
	}
	if user.Email != "" {
		userModel.Email = user.Email
	}
	if user.Phone != "" {
		userModel.Phone = user.Phone
	}
	return userModel
}
