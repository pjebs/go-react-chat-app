package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/olahol/melody"
  "net/http"
//  "github.com/gorilla/websocket"
//  "github.com/garyburd/redigo/redis"
  "github.com/soveran/redisurl"
  "os"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "time"
)

type message struct {
  Handle string `json:"handle"`
  Text   string `json:"text"`
}

func main() {
  r := gin.Default()
  m := melody.New()

  db, _ := sql.Open("mysql", "abhishekpillai:@/sharedchat?charset=utf8")
  insert_stmt, _ := db.Prepare("INSERT messages SET user_id=?,message_type=?,content=?,timestamp=?")

  // Connect using os.Getenv("REDIS_URL").
  conn, err := redisurl.Connect()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer conn.Close()

  r.GET("/", func(c *gin.Context) {
    http.ServeFile(c.Writer, c.Request, "./index.html")
  })

  r.GET("/allMessages", func(c *gin.Context) {
    var (
      content string
      allMessages []string
    )

    //allMessages, _ := redis.Strings(conn.Do("LRANGE", "messages", 0, 1000))
    rows, _ := db.Query("SELECT content FROM messages ORDER BY timestamp ASC")
    defer rows.Close()
    for rows.Next() {
      err := rows.Scan(&content)
      if err != nil {
        fmt.Println(err)
      }
      allMessages = append(allMessages, content)
    }
    c.JSON(200, gin.H{
      "messages": allMessages,
    })
  })

  r.GET("/ws", func(c *gin.Context) {
    m.HandleRequest(c.Writer, c.Request)
  })

  m.HandleMessage(func(s *melody.Session, msg []byte) {
    // conn.Do("RPUSH", "messages", string(msg))
    _, err := insert_stmt.Exec(1, "text", string(msg), time.Now())
    if err != nil {
      fmt.Println(err)
    }
    m.Broadcast(msg)
  })

  r.Run(":5000")
}
