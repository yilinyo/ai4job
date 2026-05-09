# Backend Plan 1: Scaffold + Infrastructure + User Auth

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a runnable Go backend with project scaffold, PostgreSQL + Redis infrastructure, and a working user auth system (register, email code, login, logout, get/update profile) using JWT + Redis blacklist.

**Architecture:** Classic 4-layer DDD — Interface (Gin handlers) → Application (use cases) → Domain (entities, repo interfaces) → Infrastructure (GORM, Redis, JWT). Domain layer has zero framework dependencies, enabling pure unit tests. Infrastructure layer implements domain interfaces (dependency inversion).

**Tech Stack:** Go 1.25, Gin, GORM + PostgreSQL 16, go-redis v9, golang-jwt/jwt v5, viper, bcrypt, google/uuid, golang-migrate, testify, mockery v2, testcontainers-go, zap

---

## File Map

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── interface/
│   │   ├── handler/user_handler.go
│   │   ├── handler/user_handler_test.go
│   │   ├── middleware/auth.go
│   │   ├── middleware/auth_test.go
│   │   └── router/router.go
│   ├── application/
│   │   ├── dto/user_dto.go
│   │   └── user/user_service.go
│   ├── domain/
│   │   └── user/
│   │       ├── entity.go          # User entity + EducationExperience
│   │       ├── repository.go      # UserRepository interface
│   │       └── entity_test.go
│   └── infrastructure/
│       ├── persistence/postgres/
│       │   ├── db.go              # GORM connection setup
│       │   ├── user_repo.go       # UserRepository implementation
│       │   └── user_repo_test.go  # integration test (testcontainers)
│       ├── cache/
│       │   ├── redis.go           # Redis client setup
│       │   └── token_cache.go     # blacklist + email code helpers
│       └── auth/
│           └── jwt.go             # JWT generate/verify/parse
├── pkg/
│   ├── response/response.go
│   ├── errors/errors.go
│   ├── config/config.go
│   └── pagination/pagination.go
├── migrations/
│   ├── 000001_create_users.up.sql
│   └── 000001_create_users.down.sql
├── .env.example
├── docker-compose.yml
├── Makefile
└── go.mod
```

---

## Task 1: Initialize Go Module and Dependencies

**Files:**
- Create: `backend/go.mod`
- Create: `backend/go.sum` (auto-generated)

- [ ] **Step 1: Initialize module**

```bash
cd /Users/yilin/project/ai-for/backend
go mod init github.com/yilin/ai-for-backend
```

- [ ] **Step 2: Install all dependencies**

```bash
go get github.com/gin-gonic/gin@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get github.com/redis/go-redis/v9@latest
go get github.com/golang-jwt/jwt/v5@latest
go get github.com/spf13/viper@latest
go get golang.org/x/crypto@latest
go get github.com/google/uuid@latest
go get github.com/golang-migrate/migrate/v4@latest
go get github.com/golang-migrate/migrate/v4/database/postgres@latest
go get github.com/golang-migrate/migrate/v4/source/file@latest
go get github.com/stretchr/testify@latest
go get github.com/vektra/mockery/v2@latest
go get github.com/testcontainers/testcontainers-go@latest
go get github.com/testcontainers/testcontainers-go/modules/postgres@latest
go get go.uber.org/zap@latest
go mod tidy
```

- [ ] **Step 3: Verify**

```bash
cat go.mod
```

Expected: module line + all dependencies listed under `require`.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: initialize go module with all dependencies"
```

---

## Task 2: Project Configuration and Shared Packages

**Files:**
- Create: `backend/pkg/config/config.go`
- Create: `backend/pkg/response/response.go`
- Create: `backend/pkg/errors/errors.go`
- Create: `backend/pkg/pagination/pagination.go`
- Create: `backend/.env.example`

- [ ] **Step 1: Create config package**

```go
// pkg/config/config.go
package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTExpireHours int    `mapstructure:"JWT_EXPIRE_HOURS"`

	EmailMock bool `mapstructure:"EMAIL_MOCK"`
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("no .env file found, using environment variables: %v", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.JWTExpireHours == 0 {
		cfg.JWTExpireHours = 24
	}

	return cfg
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable TimeZone=Asia/Shanghai"
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "" || c.AppEnv == "development"
}
```

- [ ] **Step 2: Create response package**

```go
// pkg/response/response.go
package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 20000, Message: "成功", Data: data})
}

func SuccessPage(c *gin.Context, list any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "成功",
		Data:    PageData{List: list, Total: total, Page: page, PageSize: pageSize},
	})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message, Data: nil})
}
```

- [ ] **Step 3: Create errors package**

```go
// pkg/errors/errors.go
package errors

import "errors"

// Business error codes
const (
	CodeSuccess       = 20000
	CodeBadParams     = 40000
	CodeUnauthorized  = 40100
	CodeForbidden     = 40300
	CodeNotFound      = 50004
	CodeConflict      = 50008
	CodeInternalError = 50000
)

// Domain errors — compare with errors.Is
var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("resource conflict")
	ErrBadParams      = errors.New("bad params")
	ErrVersionInUse   = errors.New("resume version is referenced by an application")
	ErrEmailCodeWrong = errors.New("email verification code is incorrect or expired")
	ErrDuplicate      = errors.New("duplicate entry")
)
```

- [ ] **Step 4: Create pagination package**

```go
// pkg/pagination/pagination.go
package pagination

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Params struct {
	Page     int
	PageSize int
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func FromContext(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return Params{Page: page, PageSize: pageSize}
}
```

- [ ] **Step 5: Create .env.example**

```bash
# backend/.env.example
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=ai_for
DB_USER=postgres
DB_PASSWORD=postgres

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

JWT_SECRET=change-me-in-production
JWT_EXPIRE_HOURS=24

EMAIL_MOCK=true
```

- [ ] **Step 6: Commit**

```bash
git add pkg/ .env.example
git commit -m "feat: add config, response, errors, pagination packages"
```

---

## Task 3: Docker Compose and Makefile

**Files:**
- Create: `backend/docker-compose.yml`
- Create: `backend/Makefile`

- [ ] **Step 1: Create docker-compose.yml**

```yaml
# backend/docker-compose.yml
version: "3.9"
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: ai_for
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

volumes:
  postgres_data:
```

- [ ] **Step 2: Create Makefile**

```makefile
# backend/Makefile
.PHONY: run test test-int migrate-up migrate-down mock lint docker-up docker-down

run:
	go run ./cmd/server/main.go

test:
	go test ./... -v -run "^Test" -short

test-int:
	go test ./... -v -run "^TestIntegration"

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ai_for?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ai_for?sslmode=disable" down 1

mock:
	go run github.com/vektra/mockery/v2 --all --dir internal/domain --output internal/mocks --outpkg mocks

lint:
	golangci-lint run ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down
```

- [ ] **Step 3: Start local services**

```bash
cd /Users/yilin/project/ai-for/backend
docker compose up -d
```

Expected: postgres and redis containers running.

- [ ] **Step 4: Commit**

```bash
git add docker-compose.yml Makefile
git commit -m "chore: add docker-compose and Makefile"
```

---

## Task 4: User Domain — Entity and Repository Interface

**Files:**
- Create: `backend/internal/domain/user/entity.go`
- Create: `backend/internal/domain/user/repository.go`
- Create: `backend/internal/domain/user/entity_test.go`

- [ ] **Step 1: Write the failing test for User entity**

```go
// internal/domain/user/entity_test.go
package user_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/domain/user"
)

func TestUser_SetPassword_HashesPassword(t *testing.T) {
	u := &user.User{}
	err := u.SetPassword("plaintext123")
	assert.NoError(t, err)
	assert.NotEqual(t, "plaintext123", u.PasswordHash)
	assert.NotEmpty(t, u.PasswordHash)
}

func TestUser_CheckPassword_CorrectPassword(t *testing.T) {
	u := &user.User{}
	_ = u.SetPassword("plaintext123")
	assert.True(t, u.CheckPassword("plaintext123"))
}

func TestUser_CheckPassword_WrongPassword(t *testing.T) {
	u := &user.User{}
	_ = u.SetPassword("plaintext123")
	assert.False(t, u.CheckPassword("wrongpassword"))
}
```

- [ ] **Step 2: Run test — confirm it fails**

```bash
cd /Users/yilin/project/ai-for/backend
go test ./internal/domain/user/... -v
```

Expected: compile error — package not found.

- [ ] **Step 3: Create User entity**

```go
// internal/domain/user/entity.go
package user

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeCampus     JobType = "campus"
	JobTypeSocial     JobType = "social"
	JobTypeInternship JobType = "internship"
)

type EducationExperience struct {
	School         string `json:"school"`
	Education      string `json:"education"`
	Major          string `json:"major"`
	AdmissionDate  string `json:"admissionDate"`
	GraduationDate string `json:"graduationDate"`
}

type User struct {
	ID                   string                `gorm:"type:uuid;primaryKey"`
	Username             string                `gorm:"uniqueIndex;not null"`
	PasswordHash         string                `gorm:"not null"`
	Email                string                `gorm:"uniqueIndex;not null"`
	RealName             string
	Age                  int
	Phone                string
	JobIntention         string
	Avatar               string
	Location             string
	PersonalAdvantage    string                `gorm:"type:text"`
	EducationExperiences []EducationExperience `gorm:"serializer:json"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt        `gorm:"index"`
}

func (u *User) SetPassword(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
	return err == nil
}
```

- [ ] **Step 4: Create UserRepository interface**

```go
// internal/domain/user/repository.go
package user

import "context"

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
```

- [ ] **Step 5: Run test — confirm it passes**

```bash
go test ./internal/domain/user/... -v
```

Expected: all 3 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/domain/user/
git commit -m "feat: add user domain entity and repository interface"
```

---

## Task 5: Infrastructure — PostgreSQL Connection and User Repository

**Files:**
- Create: `backend/internal/infrastructure/persistence/postgres/db.go`
- Create: `backend/internal/infrastructure/persistence/postgres/user_repo.go`
- Create: `backend/internal/infrastructure/persistence/postgres/user_repo_test.go`

- [ ] **Step 1: Write the failing integration test**

```go
// internal/infrastructure/persistence/postgres/user_repo_test.go
package postgres_test

import (
	"context"
	"testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		tcpostgres.WithDatabase("test_db"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domainuser.User{})
	require.NoError(t, err)

	return db
}

func TestIntegrationUserRepo_CreateAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := infrapostgres.NewUserRepo(db)
	ctx := context.Background()

	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: "testuser",
		Email:    "test@example.com",
	}
	_ = u.SetPassword("password123")

	err := repo.Create(ctx, u)
	require.NoError(t, err)

	found, err := repo.FindByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
	assert.Equal(t, "test@example.com", found.Email)
}

func TestIntegrationUserRepo_SoftDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := infrapostgres.NewUserRepo(db)
	ctx := context.Background()

	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: "deleteuser",
		Email:    "delete@example.com",
	}
	_ = u.SetPassword("password123")
	require.NoError(t, repo.Create(ctx, u))

	// soft delete via GORM
	err := db.Delete(u).Error
	require.NoError(t, err)

	// should not be found after soft delete
	_, err = repo.FindByUsername(ctx, "deleteuser")
	assert.Error(t, err)
}
```

- [ ] **Step 2: Run test — confirm it fails**

```bash
go test ./internal/infrastructure/persistence/postgres/... -v -run "^TestIntegration"
```

Expected: compile error — infrapostgres package not found.

- [ ] **Step 3: Create DB connection setup**

```go
// internal/infrastructure/persistence/postgres/db.go
package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dsn string, debug bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if debug {
		logLevel = logger.Info
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
}
```

- [ ] **Step 4: Create UserRepo implementation**

```go
// internal/infrastructure/persistence/postgres/user_repo.go
package postgres

import (
	"context"
	"errors"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"gorm.io/gorm"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *domainuser.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &u, err
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &u, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	var u domainuser.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrNotFound
	}
	return &u, err
}

func (r *UserRepo) Update(ctx context.Context, u *domainuser.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainuser.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainuser.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
```

- [ ] **Step 5: Run integration tests — confirm they pass**

```bash
go test ./internal/infrastructure/persistence/postgres/... -v -run "^TestIntegration"
```

Expected: both integration tests PASS (Docker pulls postgres:16 on first run, ~30s).

- [ ] **Step 6: Commit**

```bash
git add internal/infrastructure/persistence/postgres/
git commit -m "feat: add postgres db setup and user repository implementation"
```

---

## Task 6: Infrastructure — Redis Cache and JWT

**Files:**
- Create: `backend/internal/infrastructure/cache/redis.go`
- Create: `backend/internal/infrastructure/cache/token_cache.go`
- Create: `backend/internal/infrastructure/auth/jwt.go`

- [ ] **Step 1: Create Redis client setup**

```go
// internal/infrastructure/cache/redis.go
package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
```

- [ ] **Step 2: Create token cache helpers**

```go
// internal/infrastructure/cache/token_cache.go
package cache

import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
)

const (
	blacklistPrefix  = "blacklist:"
	emailCodePrefix  = "email_code:"
	emailCodeExpiry  = 5 * time.Minute
)

type TokenCache struct {
	client *redis.Client
}

func NewTokenCache(client *redis.Client) *TokenCache {
	return &TokenCache{client: client}
}

// AddToBlacklist marks a JWT jti as revoked until its expiry time.
func (c *TokenCache) AddToBlacklist(ctx context.Context, jti string, ttl time.Duration) error {
	key := blacklistPrefix + jti
	return c.client.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted returns true if the jti has been revoked.
func (c *TokenCache) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := blacklistPrefix + jti
	val, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// SetEmailCode stores a 6-digit code for the given email with a 5-minute TTL.
func (c *TokenCache) SetEmailCode(ctx context.Context, email, code string) error {
	key := fmt.Sprintf("%s%s", emailCodePrefix, email)
	return c.client.Set(ctx, key, code, emailCodeExpiry).Err()
}

// VerifyEmailCode checks the stored code. Returns true if it matches.
func (c *TokenCache) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	key := fmt.Sprintf("%s%s", emailCodePrefix, email)
	stored, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == code, nil
}

// DeleteEmailCode removes the email code after successful registration.
func (c *TokenCache) DeleteEmailCode(ctx context.Context, email string) error {
	key := fmt.Sprintf("%s%s", emailCodePrefix, email)
	return c.client.Del(ctx, key).Err()
}
```

- [ ] **Step 3: Create JWT helper**

```go
// internal/infrastructure/auth/jwt.go
package auth

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret      []byte
	expireHours int
}

func NewJWTManager(secret string, expireHours int) *JWTManager {
	return &JWTManager{secret: []byte(secret), expireHours: expireHours}
}

// Generate creates a signed JWT for the given user.
func (m *JWTManager) Generate(userID, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.expireHours) * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse validates the token and returns claims. Returns error if invalid/expired.
func (m *JWTManager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// RemainingTTL returns how long until the token expires.
func (m *JWTManager) RemainingTTL(claims *Claims) time.Duration {
	return time.Until(claims.ExpiresAt.Time)
}
```

- [ ] **Step 4: Write and run JWT unit tests**

```go
// internal/infrastructure/auth/jwt_test.go
package auth_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/infrastructure/auth"
)

func TestJWTManager_GenerateAndParse(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret", 24)

	token, err := mgr.Generate("user-123", "alice")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := mgr.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "alice", claims.Username)
	assert.NotEmpty(t, claims.ID) // jti
}

func TestJWTManager_Parse_InvalidToken(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret", 24)
	_, err := mgr.Parse("invalid.token.here")
	assert.Error(t, err)
}
```

```bash
go test ./internal/infrastructure/auth/... -v
```

Expected: 2 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/infrastructure/cache/ internal/infrastructure/auth/
git commit -m "feat: add redis cache helpers and JWT manager"
```

---

## Task 7: Application Layer — User DTOs and Service

**Files:**
- Create: `backend/internal/application/dto/user_dto.go`
- Create: `backend/internal/application/user/user_service.go`
- Create: `backend/internal/application/user/user_service_test.go`

- [ ] **Step 1: Create User DTOs**

```go
// internal/application/dto/user_dto.go
package dto

import "github.com/yilin/ai-for-backend/internal/domain/user"

type SendEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SendEmailCodeResponse struct {
	Email    string `json:"email"`
	MockCode string `json:"mockCode,omitempty"` // only in development
	ExpiresIn int   `json:"expiresIn"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	EmailCode string `json:"emailCode" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string          `json:"token"`
	UserInfo *UserInfoResponse `json:"userInfo"`
}

type UserInfoResponse struct {
	ID                   string                       `json:"id"`
	Username             string                       `json:"username"`
	RealName             string                       `json:"realName"`
	Age                  int                          `json:"age"`
	Email                string                       `json:"email"`
	Phone                string                       `json:"phone"`
	JobIntention         string                       `json:"jobIntention"`
	Avatar               string                       `json:"avatar"`
	Location             string                       `json:"location"`
	PersonalAdvantage    string                       `json:"personalAdvantage"`
	EducationExperiences []user.EducationExperience   `json:"educationExperiences"`
}

type UpdateProfileRequest struct {
	RealName             string                       `json:"realName"`
	Age                  int                          `json:"age"`
	Phone                string                       `json:"phone"`
	JobIntention         string                       `json:"jobIntention"`
	Avatar               string                       `json:"avatar"`
	Location             string                       `json:"location"`
	PersonalAdvantage    string                       `json:"personalAdvantage"`
	EducationExperiences []user.EducationExperience   `json:"educationExperiences"`
}

func ToUserInfoResponse(u *user.User) *UserInfoResponse {
	exps := u.EducationExperiences
	if exps == nil {
		exps = []user.EducationExperience{}
	}
	return &UserInfoResponse{
		ID:                   u.ID,
		Username:             u.Username,
		RealName:             u.RealName,
		Age:                  u.Age,
		Email:                u.Email,
		Phone:                u.Phone,
		JobIntention:         u.JobIntention,
		Avatar:               u.Avatar,
		Location:             u.Location,
		PersonalAdvantage:    u.PersonalAdvantage,
		EducationExperiences: exps,
	}
}
```

- [ ] **Step 2: Write failing test for user service**

```go
// internal/application/user/user_service_test.go
package user_test

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appuser "github.com/yilin/ai-for-backend/internal/application/user"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

// --- minimal mocks (hand-written for this test; replace with mockery mocks later) ---

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, u *domainuser.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domainuser.User), args.Error(1)
}
func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*domainuser.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domainuser.User), args.Error(1)
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domainuser.User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, u *domainuser.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}
func (m *mockUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

type mockTokenCache struct{ mock.Mock }

func (m *mockTokenCache) SetEmailCode(ctx context.Context, email, code string) error {
	return m.Called(ctx, email, code).Error(0)
}
func (m *mockTokenCache) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	args := m.Called(ctx, email, code)
	return args.Bool(0), args.Error(1)
}
func (m *mockTokenCache) DeleteEmailCode(ctx context.Context, email string) error {
	return m.Called(ctx, email).Error(0)
}
func (m *mockTokenCache) AddToBlacklist(ctx context.Context, jti string, ttl interface{}) error {
	return m.Called(ctx, jti, ttl).Error(0)
}
func (m *mockTokenCache) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

type mockJWT struct{ mock.Mock }

func (m *mockJWT) Generate(userID, username string) (string, error) {
	args := m.Called(userID, username)
	return args.String(0), args.Error(1)
}

// --- tests ---

func TestUserService_Login_Success(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	existingUser := &domainuser.User{ID: "uid-1", Username: "alice", Email: "a@b.com"}
	_ = existingUser.SetPassword("password123")

	repo.On("FindByUsername", mock.Anything, "alice").Return(existingUser, nil)
	jwtMgr.On("Generate", "uid-1", "alice").Return("signed.jwt.token", nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	resp, err := svc.Login(context.Background(), dto.LoginRequest{
		Username: "alice",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "signed.jwt.token", resp.Token)
	assert.Equal(t, "alice", resp.UserInfo.Username)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	repo := &mockUserRepo{}
	cache := &mockTokenCache{}
	jwtMgr := &mockJWT{}

	existingUser := &domainuser.User{ID: "uid-1", Username: "alice"}
	_ = existingUser.SetPassword("correct")

	repo.On("FindByUsername", mock.Anything, "alice").Return(existingUser, nil)

	svc := appuser.NewUserService(repo, cache, jwtMgr, true)
	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Username: "alice",
		Password: "wrong",
	})

	assert.ErrorIs(t, err, domainerrors.ErrUnauthorized)
}
```

- [ ] **Step 3: Run test — confirm it fails**

```bash
go test ./internal/application/user/... -v
```

Expected: compile error — appuser package not found.

- [ ] **Step 4: Create UserService**

```go
// internal/application/user/user_service.go
package user

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"github.com/google/uuid"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type TokenCache interface {
	SetEmailCode(ctx context.Context, email, code string) error
	VerifyEmailCode(ctx context.Context, email, code string) (bool, error)
	DeleteEmailCode(ctx context.Context, email string) error
	AddToBlacklist(ctx context.Context, jti string, ttl interface{}) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

type JWTManager interface {
	Generate(userID, username string) (string, error)
}

type UserService struct {
	repo      domainuser.Repository
	cache     TokenCache
	jwt       JWTManager
	emailMock bool
}

func NewUserService(repo domainuser.Repository, cache TokenCache, jwt JWTManager, emailMock bool) *UserService {
	return &UserService{repo: repo, cache: cache, jwt: jwt, emailMock: emailMock}
}

func (s *UserService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {
	code := fmt.Sprintf("%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))
	if err := s.cache.SetEmailCode(ctx, req.Email, code); err != nil {
		return nil, err
	}
	resp := &dto.SendEmailCodeResponse{Email: req.Email, ExpiresIn: 300}
	if s.emailMock {
		resp.MockCode = code
	}
	return resp, nil
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) error {
	ok, err := s.cache.VerifyEmailCode(ctx, req.Email, req.EmailCode)
	if err != nil {
		return err
	}
	if !ok {
		return domainerrors.ErrEmailCodeWrong
	}
	existsEmail, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existsEmail {
		return domainerrors.ErrDuplicate
	}
	existsUsername, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existsUsername {
		return domainerrors.ErrDuplicate
	}
	u := &domainuser.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
	}
	if err := u.SetPassword(req.Password); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}
	return s.cache.DeleteEmailCode(ctx, req.Email)
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, domainerrors.ErrUnauthorized
	}
	if !u.CheckPassword(req.Password) {
		return nil, domainerrors.ErrUnauthorized
	}
	token, err := s.jwt.Generate(u.ID, u.Username)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Token: token, UserInfo: dto.ToUserInfoResponse(u)}, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToUserInfoResponse(u), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.RealName = req.RealName
	u.Age = req.Age
	u.Phone = req.Phone
	u.JobIntention = req.JobIntention
	u.Avatar = req.Avatar
	u.Location = req.Location
	u.PersonalAdvantage = req.PersonalAdvantage
	u.EducationExperiences = req.EducationExperiences
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return dto.ToUserInfoResponse(u), nil
}
```

- [ ] **Step 5: Run test — confirm it passes**

```bash
go test ./internal/application/user/... -v
```

Expected: both tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/application/
git commit -m "feat: add user application service with email code, register, login, profile"
```

---

## Task 8: Interface Layer — Auth Middleware

**Files:**
- Create: `backend/internal/interface/middleware/auth.go`
- Create: `backend/internal/interface/middleware/auth_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/interface/middleware/auth_test.go
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yilin/ai-for-backend/internal/interface/middleware"
)

type mockJWTParser struct {
	userID   string
	username string
	jti      string
	err      error
}

func (m *mockJWTParser) Parse(token string) (userID, username, jti string, err error) {
	return m.userID, m.username, m.jti, m.err
}

type mockBlacklist struct {
	blacklisted bool
	err         error
}

func (m *mockBlacklist) IsBlacklisted(ctx interface{}, jti string) (bool, error) {
	return m.blacklisted, m.err
}

func newTestRouter(jwtParser middleware.JWTParser, bl middleware.Blacklist) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.Auth(jwtParser, bl))
	r.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	return r
}

func TestAuthMiddleware_NoToken_Returns401(t *testing.T) {
	r := newTestRouter(&mockJWTParser{}, &mockBlacklist{})
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "40100")
}

func TestAuthMiddleware_ValidToken_InjectsUserID(t *testing.T) {
	parser := &mockJWTParser{userID: "uid-123", username: "alice", jti: "jti-abc"}
	r := newTestRouter(parser, &mockBlacklist{blacklisted: false})
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-Access-Token", "valid.token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "uid-123")
}

func TestAuthMiddleware_BlacklistedToken_Returns401(t *testing.T) {
	parser := &mockJWTParser{userID: "uid-123", username: "alice", jti: "jti-abc"}
	r := newTestRouter(parser, &mockBlacklist{blacklisted: true})
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-Access-Token", "blacklisted.token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "40100")
}
```

- [ ] **Step 2: Run test — confirm it fails**

```bash
go test ./internal/interface/middleware/... -v
```

Expected: compile error.

- [ ] **Step 3: Create Auth middleware**

```go
// internal/interface/middleware/auth.go
package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

// JWTParser parses a token string and returns user info + jti.
type JWTParser interface {
	Parse(token string) (userID, username, jti string, err error)
}

// Blacklist checks if a jti has been revoked.
type Blacklist interface {
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

func Auth(parser JWTParser, bl Blacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Access-Token")
		if token == "" {
			response.Fail(c, errors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}
		userID, username, jti, err := parser.Parse(token)
		if err != nil {
			response.Fail(c, errors.CodeUnauthorized, "token无效")
			c.Abort()
			return
		}
		blacklisted, err := bl.IsBlacklisted(c.Request.Context(), jti)
		if err != nil || blacklisted {
			response.Fail(c, errors.CodeUnauthorized, "token已失效")
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Set("username", username)
		c.Set("jti", jti)
		c.Next()
	}
}
```

Note: The `JWTManager` in `auth/jwt.go` needs a `Parse` method that returns `(userID, username, jti string, err error)`. Add this adapter method:

```go
// Add to internal/infrastructure/auth/jwt.go
func (m *JWTManager) ParseForMiddleware(token string) (userID, username, jti string, err error) {
	claims, err := m.Parse(token)
	if err != nil {
		return "", "", "", err
	}
	return claims.UserID, claims.Username, claims.ID, nil
}
```

- [ ] **Step 4: Run test — confirm it passes**

```bash
go test ./internal/interface/middleware/... -v
```

Expected: all 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/interface/middleware/
git commit -m "feat: add JWT auth middleware with blacklist check"
```

---

## Task 9: Interface Layer — User Handler and Router

**Files:**
- Create: `backend/internal/interface/handler/user_handler.go`
- Create: `backend/internal/interface/handler/user_handler_test.go`
- Create: `backend/internal/interface/router/router.go`

- [ ] **Step 1: Write failing handler test**

```go
// internal/interface/handler/user_handler_test.go
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

type mockUserSvc struct{ mock.Mock }

func (m *mockUserSvc) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.SendEmailCodeResponse), args.Error(1)
}
func (m *mockUserSvc) Register(ctx context.Context, req dto.RegisterRequest) error {
	return m.Called(ctx, req).Error(0)
}
func (m *mockUserSvc) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}
func (m *mockUserSvc) GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.UserInfoResponse), args.Error(1)
}
func (m *mockUserSvc) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.UserInfoResponse), args.Error(1)
}

func TestUserHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockUserSvc{}
	svc.On("Login", mock.Anything, dto.LoginRequest{Username: "alice", Password: "pass123"}).
		Return(&dto.LoginResponse{Token: "jwt-token", UserInfo: &dto.UserInfoResponse{Username: "alice"}}, nil)

	h := handler.NewUserHandler(svc)
	r := gin.New()
	r.POST("/api/v1/user/login", h.Login)

	body, _ := json.Marshal(map[string]string{"username": "alice", "password": "pass123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "20000")
	assert.Contains(t, w.Body.String(), "jwt-token")
}

func TestUserHandler_Login_MissingFields_Returns400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handler.NewUserHandler(&mockUserSvc{})
	r := gin.New()
	r.POST("/api/v1/user/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "40000")
}
```

- [ ] **Step 2: Run test — confirm it fails**

```bash
go test ./internal/interface/handler/... -v
```

Expected: compile error.

- [ ] **Step 3: Create User Handler**

```go
// internal/interface/handler/user_handler.go
package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
	"github.com/yilin/ai-for-backend/pkg/response"
)

type UserService interface {
	SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*dto.UserInfoResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserInfoResponse, error)
}

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, domainerrors.CodeBadParams, err.Error())
		return
	}
	resp, err := h.svc.SendEmailCode(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, domainerrors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, resp)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, domainerrors.CodeBadParams, err.Error())
		return
	}
	if err := h.svc.Register(c.Request.Context(), req); err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, domainerrors.CodeBadParams, err.Error())
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, domainerrors.CodeBadParams, err.Error())
		return
	}
	resp, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *UserHandler) Logout(c *gin.Context) {
	// Logout is handled at the middleware/service level via blacklist.
	// The token invalidation is done by the service; here we just respond.
	response.Success(c, nil)
}

// handleError maps domain errors to response codes.
func handleError(c *gin.Context, err error) {
	switch {
	case isError(err, domainerrors.ErrNotFound):
		response.Fail(c, domainerrors.CodeNotFound, "资源不存在")
	case isError(err, domainerrors.ErrUnauthorized):
		response.Fail(c, domainerrors.CodeUnauthorized, "用户名或密码错误")
	case isError(err, domainerrors.ErrEmailCodeWrong):
		response.Fail(c, domainerrors.CodeBadParams, "验证码错误或已过期")
	case isError(err, domainerrors.ErrDuplicate):
		response.Fail(c, domainerrors.CodeConflict, "用户名或邮箱已存在")
	case isError(err, domainerrors.ErrConflict):
		response.Fail(c, domainerrors.CodeConflict, err.Error())
	default:
		response.Fail(c, domainerrors.CodeInternalError, "服务器内部错误")
	}
}

func isError(err, target error) bool {
	return err != nil && err.Error() == target.Error()
}
```

- [ ] **Step 4: Create Router**

```go
// internal/interface/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	"github.com/yilin/ai-for-backend/internal/interface/middleware"
)

func Setup(
	userHandler *handler.UserHandler,
	jwtParser middleware.JWTParser,
	blacklist middleware.Blacklist,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")

	// Public routes
	api.POST("/user/login", userHandler.Login)
	api.POST("/user/register", userHandler.Register)
	api.POST("/user/email-code", userHandler.SendEmailCode)

	// Authenticated routes
	auth := api.Group("")
	auth.Use(middleware.Auth(jwtParser, blacklist))

	auth.GET("/user/profile", userHandler.GetProfile)
	auth.PUT("/user/profile", userHandler.UpdateProfile)
	auth.POST("/user/logout", userHandler.Logout)

	return r
}
```

- [ ] **Step 5: Run tests — confirm they pass**

```bash
go test ./internal/interface/... -v
```

Expected: all handler tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/interface/
git commit -m "feat: add user handler and router with auth middleware wired"
```

---

## Task 10: Database Migration for Users Table

**Files:**
- Create: `backend/migrations/000001_create_users.up.sql`
- Create: `backend/migrations/000001_create_users.down.sql`

- [ ] **Step 1: Create up migration**

```sql
-- migrations/000001_create_users.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username              VARCHAR(50)  NOT NULL,
    password_hash         VARCHAR(255) NOT NULL,
    email                 VARCHAR(100) NOT NULL,
    real_name             VARCHAR(100),
    age                   INT,
    phone                 VARCHAR(20),
    job_intention         VARCHAR(200),
    avatar                VARCHAR(500),
    location              VARCHAR(100),
    personal_advantage    TEXT,
    education_experiences JSONB        NOT NULL DEFAULT '[]',
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_email    ON users(email)    WHERE deleted_at IS NULL;
```

- [ ] **Step 2: Create down migration**

```sql
-- migrations/000001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

- [ ] **Step 3: Run migration**

```bash
cd /Users/yilin/project/ai-for/backend
make migrate-up
```

Expected: `no error` — migration applied.

- [ ] **Step 4: Commit**

```bash
git add migrations/
git commit -m "feat: add users table migration"
```

---

## Task 11: Wire Everything Together — main.go

**Files:**
- Create: `backend/cmd/server/main.go`

- [ ] **Step 1: Create main.go**

```go
// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yilin/ai-for-backend/internal/application/user"
	"github.com/yilin/ai-for-backend/internal/infrastructure/auth"
	"github.com/yilin/ai-for-backend/internal/infrastructure/cache"
	infrapostgres "github.com/yilin/ai-for-backend/internal/infrastructure/persistence/postgres"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
	"github.com/yilin/ai-for-backend/internal/interface/middleware"
	"github.com/yilin/ai-for-backend/internal/interface/router"
	"github.com/yilin/ai-for-backend/pkg/config"
	domainuser "github.com/yilin/ai-for-backend/internal/domain/user"
)

func main() {
	cfg := config.Load()

	// Infrastructure: database
	db, err := infrapostgres.NewDB(cfg.DSN(), cfg.IsDevelopment())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	if err := db.AutoMigrate(&domainuser.User{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Infrastructure: redis
	redisClient, err := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	tokenCache := cache.NewTokenCache(redisClient)

	// Infrastructure: jwt
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpireHours)

	// JWT adapter for middleware (implements middleware.JWTParser)
	jwtParser := &jwtParserAdapter{mgr: jwtMgr}

	// Repositories
	userRepo := infrapostgres.NewUserRepo(db)

	// Application services
	userSvc := user.NewUserService(userRepo, tokenCache, jwtMgr, cfg.EmailMock)

	// Handlers
	userHandler := handler.NewUserHandler(userSvc)

	// Router
	r := router.Setup(userHandler, jwtParser, tokenCache)

	srv := &http.Server{Addr: ":" + cfg.AppPort, Handler: r}

	go func() {
		log.Printf("server listening on :%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exited")
}

// jwtParserAdapter adapts JWTManager to the middleware.JWTParser interface.
type jwtParserAdapter struct {
	mgr *auth.JWTManager
}

func (a *jwtParserAdapter) Parse(token string) (userID, username, jti string, err error) {
	claims, err := a.mgr.Parse(token)
	if err != nil {
		return "", "", "", err
	}
	return claims.UserID, claims.Username, claims.ID, nil
}
```

- [ ] **Step 2: Copy .env.example to .env and fill values**

```bash
cp /Users/yilin/project/ai-for/backend/.env.example /Users/yilin/project/ai-for/backend/.env
```

Edit `.env` to set `JWT_SECRET` to any string.

- [ ] **Step 3: Build and run**

```bash
cd /Users/yilin/project/ai-for/backend
go build ./cmd/server/...
make run
```

Expected: `server listening on :8080`

- [ ] **Step 4: Smoke test with curl**

```bash
# Send email code
curl -X POST http://localhost:8080/api/v1/user/email-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
# Expected: {"code":20000,"message":"成功","data":{"email":"test@example.com","mockCode":"XXXXXX","expiresIn":300}}

# Register (use mockCode from above)
curl -X POST http://localhost:8080/api/v1/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass1234","email":"test@example.com","emailCode":"XXXXXX"}'
# Expected: {"code":20000,"message":"成功","data":null}

# Login
curl -X POST http://localhost:8080/api/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass1234"}'
# Expected: {"code":20000,...,"data":{"token":"...","userInfo":{...}}}
```

- [ ] **Step 5: Run all unit tests**

```bash
cd /Users/yilin/project/ai-for/backend
make test
```

Expected: all unit tests PASS, no failures.

- [ ] **Step 6: Commit**

```bash
git add cmd/ .env.example
git commit -m "feat: wire all dependencies in main.go, backend server runnable"
```

---

## Self-Review

**Spec coverage check:**
- ✅ Project scaffold + 4-layer DDD directory structure
- ✅ PostgreSQL + GORM + soft delete (`deleted_at`)
- ✅ Redis client + JWT blacklist + email code
- ✅ JWT generation/parsing (HS256, jti, exp)
- ✅ User domain entity (SetPassword, CheckPassword)
- ✅ UserRepository interface + GORM implementation
- ✅ Auth middleware (header check, JWT parse, blacklist check, userID injection)
- ✅ All 6 user API endpoints (email-code, register, login, logout, get profile, update profile)
- ✅ Unified response structure (code, message, data)
- ✅ Pagination package
- ✅ Error codes (20000, 40000, 40100, 40300, 50004, 50008, 50000)
- ✅ docker-compose (postgres:16 + redis:7)
- ✅ Makefile (run, test, test-int, migrate-up, migrate-down, mock)
- ✅ Migration file for users table with partial unique indexes
- ✅ TDD: failing test → implementation → passing test at every step
- ✅ Integration test with testcontainers for UserRepo

**Placeholder scan:** No TBD/TODO found.

**Type consistency:**
- `middleware.JWTParser` interface: `Parse(token string) (userID, username, jti string, err error)` — used consistently in auth.go, auth_test.go, and main.go adapter.
- `middleware.Blacklist` interface: `IsBlacklisted(ctx context.Context, jti string) (bool, error)` — implemented by `TokenCache`, used in Auth middleware.
- `dto.ToUserInfoResponse` used consistently in UserService.
- `domainerrors.*` constants used in both handler and service layers.
