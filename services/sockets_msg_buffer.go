package services

import "sync"

type LiveChatMessageBuffer struct {
	MaxLength int
	items     []*ChatMsg
}

func (buf *LiveChatMessageBuffer) Push(msg *ChatMsg) {

	// If there is still room under the max, add it
	if len(buf.items) < buf.MaxLength {
		buf.items = append(buf.items, msg)
		return
	}

	// Move everything over one space
	for i := 1; i < len(buf.items); i++ {
		buf.items[i-1] = buf.items[i]
	}

	// Insert the new message in the last slot
	buf.items[len(buf.items)-1] = msg

}

func (buf *LiveChatMessageBuffer) GetCopy() []*ChatMsg {

	// Create the new slice for elements
	items := make([]*ChatMsg, len(buf.items))

	// Copy all the elements
	for i := range buf.items {
		items[i] = buf.items[i]
	}

	// Return the new slice
	return items

}

type LiveChatBufferGroup struct {
	streamChatBuffers    map[uint64]*LiveChatMessageBuffer
	streamChatBuffersMut sync.RWMutex
}

func (s *LiveChatBufferGroup) PushMessage(streamID uint64, msg *ChatMsg) {

	// Lock on the buffers
	s.streamChatBuffersMut.Lock()
	defer s.streamChatBuffersMut.Unlock()

	// If the buffers map is nil, create it
	if s.streamChatBuffers == nil {
		s.streamChatBuffers = map[uint64]*LiveChatMessageBuffer{}
	}

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	if !ok {
		buf = &LiveChatMessageBuffer{
			MaxLength: 25,
		}
		s.streamChatBuffers[streamID] = buf
	}

	// Push the message
	buf.Push(msg)

}

func (s *LiveChatBufferGroup) CopyMessages(streamID uint64) []*ChatMsg {

	// Lock on the buffers
	s.streamChatBuffersMut.RLock()
	defer s.streamChatBuffersMut.RUnlock()

	// If the buffers map is nil, return nil
	if s.streamChatBuffers == nil {
		return nil
	}

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	if !ok {
		return nil
	}

	// Copy the values from the buffer
	return buf.GetCopy()

}
