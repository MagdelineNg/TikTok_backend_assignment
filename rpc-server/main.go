package main

import (
	"context"
	"fmt"
	"log"

	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	rdb = &RedisClient{} // make the RedisClient with global visibility in the 'main' scope
)

func main() {
	ctx := context.Background() // https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go

	// connect to the Redis server on "redis:6379" and no password
	//cannot use "localhost:6379" here because it's inside the docker network.
	//Use the service name defined in docker-compose as the hostname.
	err := rdb.InitClient(ctx, "redis:6379", "")

	if err != nil {
		errMsg := fmt.Sprintf("Failed to init Redis client, err: %v", err)
		log.Fatal(errMsg)
	}

	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
