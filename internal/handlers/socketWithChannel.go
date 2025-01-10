package handlers

import (
	"log"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	goChannelsMap = make(map[int][]chan string)
	mu            sync.Mutex
)

func SocketHandlerWithChannel(c *gin.Context) {
	wsChannelIDStr := c.Param("channelID")
	token := c.Query("token")
	log.Println("wsChannelIDStr: " + wsChannelIDStr)
	wsChannelID, err := strconv.Atoi(wsChannelIDStr)
	if err != nil {
		panic(err)
	}
	log.Println("wschannelID: " + wsChannelIDStr + ", token: " + token)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	addListenerWithChannel(wsChannelID, ws)

}

func addListenerWithChannel(wsChannelID int, ws *websocket.Conn) {
	ch := make(chan string)
	addListenerToWsChannel(wsChannelID, ch)

	done := make(chan bool)
	defer func() {
		ws.Close()
		mu.Lock()
		defer mu.Unlock()
		channelSlices := goChannelsMap[wsChannelID]
		for i, _ch := range channelSlices {
			if ch == _ch {
				channelSlices = append(channelSlices[:i], channelSlices[i+1:]...)
				goChannelsMap[wsChannelID] = channelSlices
				log.Println("remove listener: " + strconv.Itoa(wsChannelID) + " channel counts: " + strconv.Itoa(len(channelSlices)))
				// log.Println("channelSlices: ", channelSlices)
				break
			}
		}
		close(ch)
		close(done)
	}()
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("client websocket closed")
				done <- true
				break
			}
			// @todo Lock here?
			for _, ch := range goChannelsMap[wsChannelID] {
				ch <- string(msg)
			}
		}
	}()

	for {
		select {
		case msg := <-ch:
			log.Println("received message: " + msg)
			ws.WriteMessage(websocket.TextMessage, []byte("Server says: "+msg))

		case <-done:
			return
		}
	}
}

func addListenerToWsChannel(wsChannelID int, ch chan string) {
	mu.Lock()
	defer mu.Unlock()
	if goChannelsMap[wsChannelID] == nil {
		goChannelsMap[wsChannelID] = make([]chan string, 0)
	}
	channelSlices := goChannelsMap[wsChannelID]
	channelSlices = append(channelSlices, ch)
	goChannelsMap[wsChannelID] = channelSlices
	log.Println("add listener: " + strconv.Itoa(wsChannelID) + " channel counts: " + strconv.Itoa(len(channelSlices)))
}
