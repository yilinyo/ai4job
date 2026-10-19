package resume_attachment

import (
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	attachment "github.com/yilin/ai-for-backend/internal/domain/resume_attachment"
	attachmentrepo "github.com/yilin/ai-for-backend/internal/domain/resume_attachment/repository"
	domainversion "github.com/yilin/ai-for-backend/internal/domain/resume_version"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

const MaxAttachmentSize = 10 * 1024 * 1024

type VersionReader interface {
	Get(ctx context.Context, userID, id string) (*domainversion.ResumeVersion, error)
}

type Service struct {
	versions VersionReader
	repo     attachmentrepo.Repository
}

func NewService(versions VersionReader, repo attachmentrepo.Repository) *Service {
	return &Service{versions: versions, repo: repo}
}

func (s *Service) List(ctx context.Context, userID, versionID string) ([]attachment.ResumeAttachment, error) {
	if _, err := s.versions.Get(ctx, userID, versionID); err != nil {
		return nil, err
	}
	return s.repo.ListByVersionID(ctx, versionID)
}

func (s *Service) Upload(ctx context.Context, userID, versionID string, header *multipart.FileHeader) (*attachment.ResumeAttachment, error) {
	if _, err := s.versions.Get(ctx, userID, versionID); err != nil {
		return nil, err
	}
	if header == nil || header.Size <= 0 || header.Size > MaxAttachmentSize {
		return nil, domainerrors.ErrBadParams
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" && ext != ".md" {
		return nil, domainerrors.ErrBadParams
	}
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, MaxAttachmentSize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > MaxAttachmentSize {
		return nil, domainerrors.ErrBadParams
	}
	if ext == ".pdf" && !strings.HasPrefix(string(data), "%PDF-") {
		return nil, domainerrors.ErrBadParams
	}
	contentType := "text/markdown; charset=utf-8"
	if ext == ".pdf" {
		contentType = "application/pdf"
	}
	value := &attachment.ResumeAttachment{ID: uuid.NewString(), VersionID: versionID, FileName: filepath.Base(header.Filename), ContentType: contentType, FileSize: int64(len(data)), Data: data}
	if err := s.repo.Create(ctx, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (s *Service) Download(ctx context.Context, userID, versionID, id string) (*attachment.ResumeAttachment, error) {
	if _, err := s.versions.Get(ctx, userID, versionID); err != nil {
		return nil, err
	}
	return s.repo.FindByIDAndVersionID(ctx, id, versionID)
}

func (s *Service) Delete(ctx context.Context, userID, versionID, id string) error {
	if _, err := s.versions.Get(ctx, userID, versionID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id, versionID)
}
