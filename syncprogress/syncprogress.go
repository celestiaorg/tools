package syncprogress

import (
	"fmt"
	"time"
)

type SyncProgress struct {
	messages []Message
	stopChan chan bool
}

type Message struct {
	msg string
}

func New() *SyncProgress {
	return &SyncProgress{
		messages: []Message{},
		stopChan: make(chan bool),
	}
}

func (s *SyncProgress) AddNewMessage() *Message {
	s.messages = append(s.messages, Message{})
	return &s.messages[len(s.messages)-1] // we need a pointer to the latest entry here
}

func (s *SyncProgress) AddNewMessages(batch int) {
	for i := 0; i < batch; i++ {
		s.AddNewMessage()
	}
}

func (s *SyncProgress) GetMessageRef(id int) *Message {
	return &s.messages[id]
}

func (s *SyncProgress) Start() {
	go func(s *SyncProgress) {
		tick := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-s.stopChan:
				{
					fmt.Print("\033[H\033[2J")
					return
				}
			case <-tick.C:
				{
					fmt.Print("\033[H\033[2J") // Clearing up the screen (tested only on Linux)
					for i := range s.messages {
						fmt.Println(s.messages[i].msg)
					}
				}
			}
		}
	}(s)
}

func (s *SyncProgress) Stop() {
	s.stopChan <- true
}

func (m *Message) Update(msg string) {
	m.msg = msg
}

func (m *Message) Append(msg string) {
	m.msg += msg
}

func (m *Message) Clear() {
	m.Update("")
}
