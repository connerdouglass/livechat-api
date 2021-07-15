package services

import (
	"errors"
	"fmt"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type SocketContext struct {
	User *TelegramUser
}

type SocketsService struct {
	Server               *socketio.Server
	TelegramService      *TelegramService
	ChatService          *ChatService
	streamChatBuffers    map[uint64]*LiveChatMessageBuffer
	streamChatBuffersMut sync.Mutex
}

func (s *SocketsService) Setup() {

	// Create the buffer
	s.streamChatBuffers = map[uint64]*LiveChatMessageBuffer{}

	// Add handlers to the socket server
	s.Server.OnConnect("/", func(conn socketio.Conn) error {
		fmt.Println("client connected: ", conn.RemoteAddr().String())
		conn.SetContext(SocketContext{})
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
	conn.Join(
		fmt.Sprintf("chatroom_%s", chatRoom.Identifier),
	)

	// Emit all the buffered messages to the new viewer, so they don't open the page to
	// a completely empty live chat screen
	bufMsgs := s.copyChatMsgBuffer(chatRoom.ID)
	for _, msg := range bufMsgs {
		conn.Emit(
			"chat.message",
			map[string]interface{}{
				"username":  msg.User.Username,
				"photo_url": msg.User.PhotoUrl,
				"message":   msg.Message,
			},
		)
	}

	fmt.Println("joined stream: ", chatRoom.Identifier, conn.RemoteAddr().String())

	// // Update the viewer count
	// go s.StreamsService.UpdateViewerCount(
	// 	chatRoom,
	// 	s.Server.RoomLen(
	// 		"/",
	// 		fmt.Sprintf("chatroom_%s", chatRoom.Identifier),
	// 	),
	// )

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
	conn.Leave(
		fmt.Sprintf("chatroom_%s", chatRoom.Identifier),
	)

	// // Update the viewer count
	// go s.StreamsService.UpdateViewerCount(
	// 	stream,
	// 	s.Server.RoomLen(
	// 		"/",
	// 		fmt.Sprintf("chatroom_%s", stream.Identifier),
	// 	),
	// )

	fmt.Println("left stream: ", chatRoom.Identifier, conn.RemoteAddr().String())

	return nil

}

//====================================================================================================
// chatroom.message event handler
// Called when a viewer sends a message in the chat
//====================================================================================================

type ChatMsg struct {
	ChatRoomIdentifier string       `json:"chat_room_identifier"`
	Message            string       `json:"message"`
	User               TelegramUser `json:"user"`
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

	// Validate the telegram user
	if !s.TelegramService.Verify(&data.User) {
		fmt.Println("Verification failed!")
		// return errors.New("invalid Telegram user hash")
	}

	// Check if the user is muted in chat
	muted, err := s.ChatService.IsUserMuted(
		chatRoom.OrganizationID,
		data.User.Username,
	)
	if err != nil {
		return err
	}
	if muted {
		return errors.New("user is muted in chat")
	}

	// Broadcast the message to the room
	s.Broadcast(
		fmt.Sprintf("chatroom_%s", chatRoom.Identifier),
		"chat.message",
		map[string]interface{}{
			"username":  data.User.Username,
			"photo_url": data.User.PhotoUrl,
			"message":   data.Message,
		},
	)

	// Push the chat message to the buffer
	// Do it in a goroutine because we don't care about the result and we don't want to block
	// the socket handler just to do this task
	go s.pushChatMsgToBuffer(chatRoom.ID, &data)

	return nil

}

func (s *SocketsService) pushChatMsgToBuffer(streamID uint64, msg *ChatMsg) {

	// Lock on the buffers
	s.streamChatBuffersMut.Lock()

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	if !ok {
		buf = &LiveChatMessageBuffer{
			MaxLength: 10,
		}
		s.streamChatBuffers[streamID] = buf
	}

	// Unlock the buffer mutex since we have a pointer to what we need now
	s.streamChatBuffersMut.Unlock()

	// Push the message
	buf.Push(msg)

}

func (s *SocketsService) copyChatMsgBuffer(streamID uint64) []*ChatMsg {

	// Lock on the buffers
	s.streamChatBuffersMut.Lock()

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	s.streamChatBuffersMut.Unlock()
	if !ok {
		return nil
	}

	// Copy the values from the buffer
	return buf.GetCopy()

}
