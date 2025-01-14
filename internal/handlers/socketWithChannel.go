package handlers

import (
	"log"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 客戶端管理結構
type ClientManager struct {
	clients    map[*websocket.Conn]bool // 活躍的客戶端連接
	Broadcast  chan []byte              // 廣播消息的通道
	register   chan *websocket.Conn     // 註冊新客戶端
	unregister chan *websocket.Conn     // 登出客戶端
	mu         sync.Mutex               // 用於保護 clients 的並發訪問
}

// 初始化客戶端管理器
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// 用戶ID對應ClientManager
var clientManagers = make(map[int]*ClientManager)

// 啓動客戶端管理器
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.mu.Lock()
			manager.clients[conn] = true
			manager.mu.Unlock()
			log.Println("New client connected:", conn.RemoteAddr())
		case conn := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				conn.Close()
				log.Println("Client disconnected:", conn.RemoteAddr())
			}
			manager.mu.Unlock()
		case message := <-manager.Broadcast:
			manager.mu.Lock()
			for conn := range manager.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error writing message:", err)
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mu.Unlock()
		}
	}
}

// 處理 WebSocket 連接
func (manager *ClientManager) HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	manager.register <- conn

	defer func() {
		manager.unregister <- conn
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Println("Received message:", string(message))
		// 將消息廣播到所有客戶端
		manager.Broadcast <- message
	}
}

func SocketHandlerWithManager(c *gin.Context) {
	channelIDStr := c.Param("channelID")
	channelID, err := strconv.Atoi(channelIDStr)
	token := c.Query("token")
	if err != nil {
		panic(err)
	}
	log.Println("channelID: " + channelIDStr + ", token: " + token)
	var manager *ClientManager
	if clientManagers[channelID] != nil {
		manager = clientManagers[channelID]
	} else {
		manager = NewClientManager()
		go manager.Start()
		clientManagers[channelID] = manager
	}
	manager.HandleConnection(c)
}
