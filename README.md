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
![Screenshot 2023-06-12 at 4.10.37 PM.png](..%2F..%2F..%2FDesktop%2FScreenshot%202023-06-12%20at%204.10.37%20PM.png)
- Error handling: Sending a text between 2 users where sender is invalid -> ```500 ```
![Screenshot 2023-06-12 at 4.14.42 PM.png](..%2F..%2F..%2F..%2F..%2Fvar%2Ffolders%2F9m%2Fcfh_sw491g92g55l_8_j2j1w0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_4ElihD%2FScreenshot%202023-06-12%20at%204.14.42%20PM.png)

2. ```Pull()``` API
- Getting texts between 2 users -> ```500 ```
![Screenshot 2023-06-12 at 8.27.13 PM.png](..%2F..%2F..%2F..%2F..%2Fvar%2Ffolders%2F9m%2Fcfh_sw491g92g55l_8_j2j1w0000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_dhrl6i%2FScreenshot%202023-06-12%20at%208.27.13%20PM.png)
- Error handling: When chatID does not exist -> ```500 ```
![Screenshot 2023-06-14 at 3.19.19 PM.png](..%2F..%2F..%2FDesktop%2FScreenshot%202023-06-14%20at%203.19.19%20PM.png)

### Concurrency Testing using Jmeter
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