package main

import (
	"./websocket" //gorilla websocket implementation
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

//############ CHAT ROOM TYPE AND METHODS

type ChatRoom struct {
	clients      map[string]Client
	mutex        sync.Mutex
	messageQueue chan string
}

//Initializing the chat room
func (chatRoom *ChatRoom) Initialize() {
	chatRoom.messageQueue = make(chan string, 5)
	chatRoom.clients = make(map[string]Client)

	go func() {
		for {
			chatRoom.BroadCast()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

//Registering a new client
//Returns pointer to a Client, or Nil, if the name is already taken
func (chatRoom *ChatRoom) Join(name string, connection *websocket.Conn) *Client {
	defer chatRoom.mutex.Unlock()

	chatRoom.mutex.Lock() //Preventing simultaneous access to the `clients` map
	if _, exists := chatRoom.clients[name]; exists {
		return nil
	}
	client := Client{
		name:       name,
		connection: connection,
		chatRoom:   chatRoom,
	}
	chatRoom.clients[name] = client

	chatRoom.AddMessage("<B>" + name + "</B> has joined the chat.")
	return &client
}

//Leaving the chat room
func (chatRoom *ChatRoom) Leave(name string) {
	chatRoom.mutex.Lock() //preventing simultaneous access to the `clients` map
	delete(chatRoom.clients, name)
	chatRoom.mutex.Unlock()
	chatRoom.AddMessage("<B>" + name + "</B> has left the chat.")
}

//Adding message to queue
func (chatRoom *ChatRoom) AddMessage(msg string) {
	chatRoom.messageQueue <- msg
}

//Broadcasting all the messages in the queue in one block
func (chatRoom *ChatRoom) BroadCast() {
	messageBlock := ""
infiniteLoop:
	for {
		select {
		case message := <-chatRoom.messageQueue:
			messageBlock += message + "<BR>"
		default:
			break infiniteLoop
		}
	}
	if len(messageBlock) > 0 {
		for _, client := range chatRoom.clients {
			client.Send(messageBlock)
		}
	}
}

//################CLIENT TYPE AND METHODS

type Client struct {
	name       string
	connection *websocket.Conn
	chatRoom   *ChatRoom
}

//Client has a new message to broadcast
func (client *Client) NewMessage(message string) {
	client.chatRoom.AddMessage("<B>" + client.name + ":</B> " + message)
}

//Exiting out
func (client *Client) Exit() {
	client.chatRoom.Leave(client.name)
}

//Sending message block to the client
func (client *Client) Send(messages string) {
	client.connection.WriteMessage(websocket.TextMessage, []byte(messages))
}

//Global variable for handling all chat traffic
var chat ChatRoom

//##############SERVING STATIC FILES
func staticFiles(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./static/"+request.URL.Path)
}

//##############HANDLING THE WEBSOCKET
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //not checking origin
}

//Handler for joining to the chat
func wsHandler(writer http.ResponseWriter, request *http.Request) {

	connection, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	go func() {
		//First message has to be the name
		_, message, err := connection.ReadMessage()
		client := chat.Join(string(message), connection)
		if client == nil || err != nil {
			connection.Close()
			return
		}

		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				client.Exit()
				return
			}
			client.NewMessage(string(message))
		}

	}()
}

//Printing out the ways the server can be reached by the clients
func printClientConnectionInfo() {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println("Chat clients can connect at the following addresses:")

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("http://" + ipnet.IP.String() + ":1717")
			}
		}

	}
}

func main() {
	printClientConnectionInfo()
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", staticFiles)
	chat.Initialize()
	http.ListenAndServe(":1717", nil)
}
