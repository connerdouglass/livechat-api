package services

import (
	"errors"
	"fmt"

	"github.com/godocompany/livechat-api/models"
	"github.com/godocompany/livechat-api/utils"
	socketio "github.com/googollee/go-socket.io"
)

type ChatUser struct {
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}

type SocketsService struct {
	Server      *socketio.Server
	ChatService *ChatService
	chatBuffers LiveChatBufferGroup
}

func socketRoomName(chatRoom *models.ChatRoom) string {
	return fmt.Sprintf("chatroom_%s", chatRoom.Identifier)
}

func (s *SocketsService) Setup() {

	// Add handlers to the socket server
	s.Server.OnConnect("/", func(conn socketio.Conn) error {
		fmt.Println("client connected: ", conn.RemoteAddr().String())
		return nil
	})

	// When a socket disconnects
	s.Server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		fmt.Println("client disconnected: ", conn.RemoteAddr().String())
		conn.LeaveAll()
	})

	// Register all of the event handlers
	s.Server.OnEvent("/", "chatroom.join", s.OnChatRoomJoin)
	s.Server.OnEvent("/", "chatroom.leave", s.OnChatRoomLeave)
	s.Server.OnEvent("/", "chatroom.message", s.OnChatRoomMessage)

}

// Broadcast broadcasts a message to every member of a room
func (s *SocketsService) Broadcast(room, event string, args ...interface{}) bool {
	return s.Server.BroadcastToRoom("/", room, event, args...)
}

//====================================================================================================
// chatroom.join event handler
// Called when a viewer joins a stream
//====================================================================================================

type ChatRoomJoinMsg struct {
	ChatRoomIdentifier string `json:"chat_room_identifier"`
}

func (s *SocketsService) OnChatRoomJoin(conn socketio.Conn, data ChatRoomJoinMsg) error {

	// Get the stream with the identifier
	chatRoom, err := s.ChatService.GetChatRoomByIdentifier(data.ChatRoomIdentifier)
	if err != nil {
		return err
	}
	if chatRoom == nil {
		return errors.New("chat room not found")
	}

	// Join the room for the event
	conn.Join(socketRoomName(chatRoom))

	// Emit all the buffered messages to the new viewer, so they don't open the page to
	// a completely empty live chat screen
	bufMsgs := s.chatBuffers.CopyMessages(chatRoom.ID)
	messagesSer := make([]map[string]interface{}, len(bufMsgs))
	for i, msg := range bufMsgs {
		messagesSer[i] = map[string]interface{}{
			"username":  msg.User.Username,
			"photo_url": msg.User.PhotoUrl,
			"message":   msg.Message,
		}
	}
	conn.Emit("chat.messages", messagesSer)

	fmt.Println("joined stream: ", chatRoom.Identifier, conn.RemoteAddr().String())

	return nil

}

//====================================================================================================
// chatroom.leave event handler
// Called when a viewer leaves a stream
//====================================================================================================

type ChatRoomLeaveMsg struct {
	ChatRoomIdentifier string `json:"chat_room_identifier"`
}

func (s *SocketsService) OnChatRoomLeave(conn socketio.Conn, data ChatRoomLeaveMsg) error {

	// Get the stream with the identifier
	chatRoom, err := s.ChatService.GetChatRoomByIdentifier(data.ChatRoomIdentifier)
	if err != nil {
		return err
	}
	if chatRoom == nil {
		return errors.New("chat room not found")
	}

	// Leave the room for the event
	conn.Leave(socketRoomName(chatRoom))

	fmt.Println("left stream: ", chatRoom.Identifier, conn.RemoteAddr().String())

	return nil

}

//====================================================================================================
// chatroom.message event handler
// Called when a viewer sends a message in the chat
//====================================================================================================

type ChatMsg struct {
	ChatRoomIdentifier string   `json:"chat_room_identifier"`
	Message            string   `json:"message"`
	User               ChatUser `json:"user"`
}

func (s *SocketsService) OnChatRoomMessage(conn socketio.Conn, data ChatMsg) error {

	// Get the stream with the identifier
	chatRoom, err := s.ChatService.GetChatRoomByIdentifier(data.ChatRoomIdentifier)
	if err != nil {
		return err
	}
	if chatRoom == nil {
		return errors.New("chat room not found")
	}

	// Check if the user is muted in chat
	muted, err := s.ChatService.IsUserMuted(
		chatRoom.OrganizationID,
		&ChatUserInfo{
			Username:  data.User.Username,
			IpAddress: utils.GetIpAddress(conn.RemoteHeader(), conn.RemoteAddr()),
		},
	)
	if err != nil {
		return err
	}
	if muted {
		return errors.New("user is muted in chat")
	}

	// Broadcast the message to the room
	go s.Broadcast(
		socketRoomName(chatRoom),
		"chat.messages",
		[]map[string]interface{}{
			{
				"username":  data.User.Username,
				"photo_url": data.User.PhotoUrl,
				"message":   data.Message,
			},
		},
	)

	// Push the chat message to the buffer
	// Do it in a goroutine because we don't care about the result and we don't want to block
	// the socket handler just to do this task
	go s.chatBuffers.PushMessage(chatRoom.ID, &data)

	return nil

}
