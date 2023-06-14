# TikTok Backend Assignment

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

This is my submission for the backend assignment of 2023 TikTok Tech Immersion.

## Run
```
docker-compose up -d
```

If ```idl_rpc.thrift``` is changed,
```
bash ./scripts/kitex-regen.sh
```

If ```idl_rpc.proto``` is changed,
```
bash ./scripts/api-regen.sh
```

## API Testing using Postman
1. ```Send()``` API
- Sending a text between 2 users -> ```200 OK```
![send-success.png](assets%2Fimages%2Fsend-success.png)
- Error handling: Sending a text between 2 users where sender is invalid -> ```500 ```
![send-fail.png](assets%2Fimages%2Fsend-fail.png)

2. ```Pull()``` API
- Getting texts between 2 users -> ```500 ```
![pull-success.png](assets%2Fimages%2Fpull-success.png)
- Error handling: When chatID does not exist -> ```500 ```
![pull-fail.png](assets%2Fimages%2Fpull-fail.png)

## Concurrency Testing using Jmeter
- Tested using 20, 500 and 1000 threads
- All but 2 threads have HTTP Status 200 OK
- Results in ```concurrency-testing-res.csv```

## Database used
- Redis
- Docker image: https://hub.docker.com/search?image_filter=official&type=image&q=
- Official documentation: https://redis.uptrace.dev/guide/go-redis.html#installation
- key: room ID eg. userA:userB
  - Used function to sort order of names in alphabetical order so that `userA:userB === userB:userA`
- value: all messages exchanged b/w A and B

## Files changed
1. ```rpc-server/db.go```
- Initialised Redis database
- Create SaveMessage function to call in ```Send()``` to write to database
- Create GetMessagesByRoomId function to call in ```Pull()``` to read from database
2. ```rpc-server/main.go```
- Connect to the Redis server on "redis:6379" and no password
3. ```rpc-server/handler.go```
- Added API business logic
4. ```scripts/kitex-regen.sh```
- To regenerate ```kitex_gen``` folder from ```idl_rpc.thrift```
5. ```scripts/api-regen.sh```
- To regenerate ```proto_gen``` folder from ```idl_http.proto```
6. ```docker-compose.yml```
- Define services (HTTP Server, RPC Server, Redis DB)


## Further improvements to be made
- Scaling the application horizontally using Kubernates cluster and increasing replica pods