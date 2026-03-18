package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/VAGRAMCHIC/vless_reality_agent/internal/api"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/domain"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/repository"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repo   repository.UserRepository
	client api.Client
	log    *utils.Logger
	server string
	tag    string
}

func NewUserService(r repository.UserRepository, log *utils.Logger, server string, tag string) *UserService {
	client := api.New(server, tag)
	return &UserService{repo: r, client: *client, log: log}
}

func (s *UserService) AddUser(ctx context.Context, email string) error {
	shortID, _ := generateShortID()
	u := &domain.User{
		ID:      uuid.New(),
		ShortID: shortID,
		Email:   email,
		Tag:     s.client.Tag,
		Server:  s.client.Server,
	}
	err := s.repo.Create(ctx, u)
	if err != nil {
		s.log.Error(ctx, "repo_error", map[string]any{
			"error": err.Error(),
		})
		return err
	}
	users, err := s.repo.GetAll(ctx)
	cfg := BuildClients(users)
	SaveConfig("/etc/xray/config.json", cfg)
	ReloadXray("xray")
	if err != nil {
		s.log.Error(ctx, "client_error", map[string]any{
			"error": err.Error(),
		})
		s.repo.DeleteUser(ctx, &u.ID)
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.client.RemoveUser(id.String())
	if err != nil {
		s.log.Error(ctx, "client_error", map[string]any{
			"error": err.Error(),
		})
		return err
	}
	err = s.repo.DeleteUser(ctx, &id)
	if err != nil {
		s.log.Error(ctx, "repo_error", map[string]any{
			"error": err.Error(),
		})
		return err
	}

	return nil

}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		s.log.Error(ctx, "repo_error", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}
	return users, nil

}

func generateShortID() (string, error) {
	b := make([]byte, 8)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func BuildClients(users []domain.User) []domain.Client {
	clients := []domain.Client{}
	for _, u := range users {
		clients = append(clients, domain.Client{
			ID:    u.ID.String(),
			Email: u.Email,
		})
	}
	return clients
}

func SaveConfig(path string, cfg any) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ReloadXray(container string) error {
	cmd := exec.Command(
		"docker",
		"kill",
		"-s",
		"HUP",
		container,
	)
	return cmd.Run()
}
