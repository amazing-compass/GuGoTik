package rpc

import (
	"GuGoTik/src/constant/config"
	"GuGoTik/src/rpc/publish"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"testing"
)

func TestListVideo(t *testing.T) {
	var Client publish.PublishServiceClient
	req := publish.ListVideoRequest{
		UserId:  123,
		ActorId: 123,
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.PublishRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = publish.NewPublishServiceClient(conn)
	//调用服务端方法
	res, err := Client.ListVideo(context.Background(), &req)
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestCountVideo(t *testing.T) {
	var Client publish.PublishServiceClient
	req := publish.CountVideoRequest{
		UserId: 1,
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.PublishRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = publish.NewPublishServiceClient(conn)
	res, err := Client.CountVideo(context.Background(), &req)
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestPublishVideo(t *testing.T) {
	var Client publish.PublishServiceClient
	reader, err := os.Open("/home/yangfeng/Repos/youthcamp/videos/upload_video_2_1080p.mp4")
	assert.Empty(t, err)
	bytes, err := io.ReadAll(reader)
	assert.Empty(t, err)
	req := publish.CreateVideoRequest{
		ActorId: 2,
		Data:    bytes,
		Title:   "原神，启动！",
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.PublishRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = publish.NewPublishServiceClient(conn)
	res, err := Client.CreateVideo(context.Background(), &req)
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}
