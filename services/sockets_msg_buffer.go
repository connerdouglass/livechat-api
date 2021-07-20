package services

type LiveChatMessageBuffer struct {
	MaxLength int
	items     []*ChatMsg
	// mut       sync.RWMutex
}

func (buf *LiveChatMessageBuffer) Push(msg *ChatMsg) {

	// // Lock on the mutex with write access
	// buf.mut.Lock()
	// defer buf.mut.Unlock()

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

	// // Lock on the mutex with readonly access
	// buf.mut.RLock()
	// defer buf.mut.RUnlock()

	// Create the new slice for elements
	items := make([]*ChatMsg, len(buf.items))

	// Copy all the elements
	for i := range buf.items {
		items[i] = buf.items[i]
	}

	// Return the new slice
	return items

}
