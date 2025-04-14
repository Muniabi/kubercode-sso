package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"kubercode-sso/internal/domain/auth/dto"
	"kubercode-sso/internal/domain/auth/values"
	sso "kubercode-sso/proto/pb/go"
	"sync"
)

type TokensRepresentation struct {
	tokens   *sso.PairTokens
	deviceId uuid.UUID
	userId   uuid.UUID
	isRevoke bool
}

func NewTokenRepresentation(tokens *sso.PairTokens, deviceId uuid.UUID, userId uuid.UUID, isRevoke bool) *TokensRepresentation {
	return &TokensRepresentation{
		tokens:   tokens,
		deviceId: deviceId,
		userId:   userId,
		isRevoke: isRevoke,
	}
}

// InMemoryAccountRepository - in-memory реализация репозитория аккаунтов
type InMemoryAccountRepository struct {
	log     *slog.Logger
	mu      sync.RWMutex
	storage map[values.Email]dto.UserDTO
	tokens  map[uuid.UUID]TokensRepresentation
}

// NewInMemoryAccountRepository - создает новый in-memory репозиторий аккаунтов
func NewInMemoryAccountRepository(log *slog.Logger) *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		storage: make(map[values.Email]dto.UserDTO),
		tokens:  make(map[uuid.UUID]TokensRepresentation),
		log:     log,
	}
}

// SaveTokens - сохраняет токены, в качестве ключа используется id агрегата
func (r *InMemoryAccountRepository) SaveTokens(ctx context.Context, tokens TokensRepresentation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[tokens.userId] = tokens
	return nil
}

// GetTokens - находит токены по id
func (r *InMemoryAccountRepository) GetTokens(ctx context.Context, userId uuid.UUID) (TokensRepresentation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tokens, ok := r.tokens[userId]
	if !ok {
		return TokensRepresentation{}, errors.New("token not found")
	}
	return tokens, nil
}

// RevokeTokens - отзывает токены у конкретного пользователя(сейчас просто удаляет, в разработке фича с отзывом токенов)
func (r *InMemoryAccountRepository) RevokeTokens(ctx context.Context, userId uuid.UUID) (TokensRepresentation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, userId)
	return TokensRepresentation{}, nil
}

// GetByEmail - находит аккаунт по ID
func (r *InMemoryAccountRepository) GetByEmail(ctx context.Context, email values.Email) (dto.UserDTO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.storage[email]
	if !exists {
		return dto.UserDTO{}, errors.New("user not found")
	}
	return user, nil
}

// Save - сохраняет аккаунт
func (r *InMemoryAccountRepository) Save(ctx context.Context, user dto.UserDTO) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[user.Email] = user
	return nil
}
