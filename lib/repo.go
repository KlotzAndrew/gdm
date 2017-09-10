package repo

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Game struct {
    GameId   string `json:"game_id"`
    PlayerId string `json:"player_id"`
    Meta     string `json:"meta"`
}

func session() *mgo.Session {
  session, err := mgo.Dial("localhost")
  if err != nil { panic(err) }

  session.SetMode(mgo.Monotonic, true)
  ensureIndex(session)

  return session
}

func ensureIndex(s *mgo.Session) {
  session := s.Copy()
  defer session.Close()

  c := session.DB("store").C("games")

  index := mgo.Index{
    Key:        []string{"game_id"},
    Unique:     true,
    DropDups:   true,
    Background: true,
    Sparse:     true,
  }
  err := c.EnsureIndex(index)
  if err != nil { panic(err) }
}

func AllGames() (games []Game, err error) {
  session := session()
  defer session.Close()

  c := session.DB("store").C("games")
  err = c.Find(bson.M{}).All(&games)

  return games, err
}

func AddGame(game Game) (err error) {
  session := session()
  defer session.Close()

  c := session.DB("store").C("games")
  err = c.Insert(game)

  return err
}

func GameById(GameId string) (game Game, err error) {
  session := session()
  defer session.Close()

  c := session.DB("store").C("games")
  err = c.Find(bson.M{"gameid": GameId}).One(&game)

  return game, err
}

func DeleteGame(GameId string) (err error) {
  session := session()
  defer session.Close()

  c := session.DB("store").C("games")
  err = c.Remove(bson.M{"gameid": GameId})

  return err
}
