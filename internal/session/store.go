package session

import "telegram-ai-bot/internal/llm"

type Store struct {
	data map[int64][]llm.Message
}

func New() *Store {
	return &Store{
		data: make(map[int64][]llm.Message),
	}
}

func (s *Store) Add(userID int64, role, text string) {
	s.data[userID] = append(s.data[userID], llm.Message{
		Role:    role,
		Content: text,
	})

	if len(s.data[userID]) > 10 {
		s.data[userID] = s.data[userID][len(s.data[userID])-10:]
	}
}

func (s *Store) Get(userID int64) []llm.Message {
	return s.data[userID]
}