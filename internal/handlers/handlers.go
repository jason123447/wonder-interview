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
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		panic(err)
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	addListener(channelID, ws)

}

// func addChannel(channelID int, conn *websocket.Conn) {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	if Channels[channelID] == nil {
// 		Channels[channelID] = make([]*websocket.Conn, 0)
// 	}
// 	Channels[channelID] = append(Channels[channelID], conn)
// 	log.Println("add channel: " + strconv.Itoa(channelID))
// }

// func removeChannel(channelID int) {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	conns, ok := Channels[channelID]
// 	if ok {
// 		for _, c := range conns {
// 			c.Close()
// 		}
// 	}
// 	delete(Channels, channelID)
// 	log.Println("remove channel: " + strconv.Itoa(channelID))
// }

func addListener(channelID int, listener *websocket.Conn) {
	if Channels[channelID] == nil {
		Channels[channelID] = make([]*websocket.Conn, 0)
		log.Println("add channel: " + strconv.Itoa(channelID))
	}
	Channels[channelID] = append(Channels[channelID], listener)
	log.Println("add listener: " + strconv.Itoa(channelID))
	listenChannel(channelID, listener)
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
			removeListener(channelID, ws)
			panic(err)
			// break
		}
		Broadcast(channelID, "Server says: "+string(message))
	}

}

func Broadcast(channelID int, message string) {
	conns, ok := Channels[channelID]
	if ok {
		for _, c := range conns {
			log.Println("send message to channel: " + strconv.Itoa(channelID))
			c.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
