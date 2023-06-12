# TikTok Backend Assignment

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

This is my submission for the backend assignment of 2023 TikTok Tech Immersion.

## Database used
- Redis
- Docker image: https://hub.docker.com/search?image_filter=official&type=image&q=
- Official documentation: https://redis.uptrace.dev/guide/go-redis.html#installation
- Commands used: ZADD, ZRANGE, ZREVRANGE
- key: room ID eg. userA:userB
  - function to sort order of names in alphabetical order so that `userA:userB === userB:userA`
- value: all messages exchanged b/w A and B

## Files changed
1. ```rpc-server/db.go```
   2. Initialised Redis database
   3. Create SaveMessage function to call in ```Send()``` to write to database
   4. Create GetMessagesByRoomId function to call in ```Pull()``` to read from database
3. ```rpc-server/main.go```
   4. Connect to the Redis server on "redis:6379" and no password
5. ```rpc-server/handler.go```
   6. Added API logic
7. ```scripts/kitex-regen.sh```
   8. To regenerate ```kitex_gen``` folder from ```idl_rpc.thrift```
9. ```scripts/api-regen.sh```
   10. To regenerate ```proto_gen``` folder from ```idl_http.proto```
9. ```docker-compose.yml```
   10. Define services (HTTP Server, RPC Server, Redis DB)