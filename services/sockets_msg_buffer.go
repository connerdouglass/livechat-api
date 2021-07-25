package services

import "sync"

type wrappedMsg struct {
	ID      string
	Message *ChatMsg
}

type LiveChatMessageBuffer struct {
	MaxLength int
	items     []*wrappedMsg
}

func (buf *LiveChatMessageBuffer) Push(msgID string, msg *ChatMsg) {

	// Create the wrapped message instance
	wmsg := &wrappedMsg{
		ID:      msgID,
		Message: msg,
	}

	// If there is still room under the max, add it
	if len(buf.items) < buf.MaxLength {
		buf.items = append(buf.items, wmsg)
		return
	}

	// Move everything over one space
	for i := 1; i < len(buf.items); i++ {
		buf.items[i-1] = buf.items[i]
	}

	// Insert the new message in the last slot
	buf.items[len(buf.items)-1] = wmsg

}

func (buf *LiveChatMessageBuffer) Revoke(msgID string) {

	// Create the new slice for the items
	items := []*wrappedMsg{}

	// Loop through the buffer
	for _, msg := range buf.items {

		// If the item matches
		if msg.ID == msgID {
			continue
		}

		// Add it to the new items
		items = append(items, msg)

	}

	// Update the items slice
	buf.items = items

}

func (buf *LiveChatMessageBuffer) GetCopy() []*wrappedMsg {

	// Create the new slice for elements
	items := make([]*wrappedMsg, len(buf.items))

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

func (s *LiveChatBufferGroup) PushMessage(streamID uint64, msgID string, msg *ChatMsg) {

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
	buf.Push(msgID, msg)

}

func (s *LiveChatBufferGroup) RevokeMessage(streamID uint64, msgID string) {

	// Lock on the buffers
	s.streamChatBuffersMut.Lock()
	defer s.streamChatBuffersMut.Unlock()

	// If the buffers map is nil, bail out now
	if s.streamChatBuffers == nil {
		return
	}

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	if !ok {
		return
	}

	// Revoke the message
	buf.Revoke(msgID)

}

func (s *LiveChatBufferGroup) CopyMessages(streamID uint64) []*wrappedMsg {

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
