package routes

// routes for the user microservice

import (
	"fmt"

	"github.com/gin-gonic/gin"
	pb "github.com/iamvasanth07/showcase/common/protos/user"
	"github.com/iamvasanth07/showcase/user/config"
	"google.golang.org/grpc"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreateRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}

type UserUpdateRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}

// UserRoutes struct
type UserRoutes struct {
	userClient pb.UserServiceClient
	config     *config.Settings
}

// NewUserRoutes returns a new user routes
func NewUserRoutes(config *config.Settings) *UserRoutes {

	client, err := grpc.Dial(fmt.Sprintf("%s:%s", config.Server.GrpcHost, config.Server.GrcpPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &UserRoutes{
		userClient: pb.NewUserServiceClient(client),
		config:     config,
	}
}

// RegisterRoutes registers the user routes
func (r *UserRoutes) RegisterUserSvcRoutes(router *gin.Engine) {

	router.GET("/user/:id", r.getUser)
	router.POST("/user", r.createUser)
	router.PUT("/user/:id", r.updateUser)
	router.DELETE("/user/:id", r.deleteUser)
	router.POST("/user/login", r.login)
}

// getUser call the user grpc service and returns a user
func (r *UserRoutes) getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := r.userClient.Get(c, &pb.GetUserRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error",
		})
		return
	}
	c.JSON(200, gin.H{user})
}

// createUser creates a user
func (r *UserRoutes) createUser(c *gin.Context) {
	body := &UserCreateRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	req := &pb.CreateUserRequest{
		User: &pb.User{
			Email:     body.Email,
			Password:  body.Password,
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Username:  body.Username,
			Phone:     body.Phone,
		},
	}
	res, err := r.userClient.Create(c, req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{res})
}

// updateUser updates a user
func (r *UserRoutes) updateUser(c *gin.Context) {
	body := &UserUpdateRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	req := &pb.UpdateUserRequest{
		User: &pb.User{
			Email:     body.Email,
			Password:  body.Password,
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Username:  body.Username,
			Phone:     body.Phone,
		},
	}
	res, err := r.userClient.Update(c, req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{res})
}

// deleteUser deletes a user
func (r *UserRoutes) deleteUser(c *gin.Context) {
	id := c.Param("id")
	res, err := r.userClient.Delete(c, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{res})
}

// login logs a user in
func (r *UserRoutes) login(c *gin.Context) {
	body := &LoginRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	req := &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	}
	res, err := r.userClient.Login(c, req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{res})
}
