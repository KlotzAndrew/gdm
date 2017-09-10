# Game Data Store

Datastore for storing and retrieving game data

 - Go server using Goji
 - datastore using MongoDB

### Usage

```shell
# add
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"game_id":"1","player_id":"1","meta":"some stuff"}' \
  http://localhost:8080/games

# all
curl \
  -H "Content-Type: application/json" \
  http://localhost:8080/games

# find
curl \
  -H "Content-Type: application/json" \
  http://localhost:8080/games/1

# delete
curl -X POST \
  -H "Content-Type: application/json" \
  http://localhost:8080/games/1
```
