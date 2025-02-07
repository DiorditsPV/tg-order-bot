package repository

import (
	"orderbot/internal/domain"
	"sync"
)

type Repository interface {
	Save(session *domain.Session)
	Get(chatID int64) *domain.Session
	CheckSession(chatID int64) bool
	SetState(chatID int64, state string)
	GetState(chatID int64) string
}

type SessionRepository struct {
	sessions map[int64]*domain.Session
	mu       sync.RWMutex
}

func NewRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[int64]*domain.Session),
	}
}

func (r *SessionRepository) Save(session *domain.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ChatID] = session
}

func (r *SessionRepository) Get(chatID int64) *domain.Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sessions[chatID]
}

func (r *SessionRepository) CheckSession(chatID int64) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.sessions[chatID]
	return ok
}

func (r *SessionRepository) SetState(chatID int64, state string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[chatID].State = state
}

func (r *SessionRepository) GetState(chatID int64) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sessions[chatID].State
}
