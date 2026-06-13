package services

import (
	"context"
	"fmt"

	"github.com/MartialM1nd/freefsm/internal/ent"
	"github.com/MartialM1nd/freefsm/internal/ent/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	client *ent.Client
}

func NewUserService(client *ent.Client) *UserService {
	return &UserService{client: client}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*ent.User, error) {
	u, err := s.client.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return u, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	u, err := s.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (s *UserService) ListAll(ctx context.Context) ([]*ent.User, error) {
	return s.client.User.Query().Order(ent.Asc(user.FieldName)).All(ctx)
}

type UserCreateParams struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func (s *UserService) Create(ctx context.Context, p UserCreateParams) (*ent.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	u, err := s.client.User.Create().
		SetName(p.Name).
		SetEmail(p.Email).
		SetPasswordHash(string(hash)).
		SetRole(p.Role).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

type UserUpdateParams struct {
	Name     *string
	Email    *string
	Role     *string
	Password *string
}

func (s *UserService) Update(ctx context.Context, id int64, p UserUpdateParams) (*ent.User, error) {
	u := s.client.User.UpdateOneID(id)
	if p.Name != nil { u.SetName(*p.Name) }
	if p.Email != nil { u.SetEmail(*p.Email) }
	if p.Role != nil { u.SetRole(*p.Role) }
	if p.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*p.Password), bcrypt.DefaultCost)
		if err != nil { return nil, fmt.Errorf("hash password: %w", err) }
		u.SetPasswordHash(string(hash))
	}
	ret, err := u.Save(ctx)
	if err != nil { return nil, fmt.Errorf("update user: %w", err) }
	return ret, nil
}

func (s *UserService) SetActive(ctx context.Context, id int64, active bool) error {
	return s.client.User.UpdateOneID(id).SetIsActive(active).Exec(ctx)
}

func (s *UserService) SetPassword(ctx context.Context, id int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return fmt.Errorf("hash password: %w", err) }
	return s.client.User.UpdateOneID(id).SetPasswordHash(string(hash)).Exec(ctx)
}
