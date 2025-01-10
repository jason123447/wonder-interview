package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Channels = make(map[int][]*websocket.Conn)
)

func SocketHandler(c *gin.Context) {
	channelIDStr := c.Param("channelID")
	token := c.Query("token")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		panic(err)
	}
	log.Println("channelID: " + channelIDStr + ", token: " + token)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	if err := addListener(channelID, ws); err != nil {
		panic(err)
	}

}

func addListener(channelID int, listener *websocket.Conn) error {
	if Channels[channelID] == nil {
		Channels[channelID] = make([]*websocket.Conn, 0)
		log.Println("add channel: " + strconv.Itoa(channelID))
	}
	// length := len(Channels[channelID])
	// if length >= connectionLimit {
	// 	return errors.New("channel connection limit reached")
	// }
	Channels[channelID] = append(Channels[channelID], listener)
	length := len(Channels[channelID])
	log.Println("add listener: "+strconv.Itoa(channelID), "listener counts: "+strconv.Itoa(length))
	listenChannel(channelID, listener)
	return nil
}

func removeListener(channelID int, listener *websocket.Conn) {
	conns, ok := Channels[channelID]
	if ok {
		for i, c := range conns {
			if c == listener {
				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		Channels[channelID] = conns
	}
	log.Println("remove listener: " + strconv.Itoa(channelID))
}

func listenChannel(channelID int, ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			// client 關閉連接
			removeListener(channelID, ws)
			// panic(err)
			break
		}
		broadcast(channelID, "Server says: "+string(message))
	}
}

func broadcast(channelID int, message string) {
	conns, ok := Channels[channelID]
	if ok {
		for _, c := range conns {
			log.Println("send message to channel: " + strconv.Itoa(channelID))
			c.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
