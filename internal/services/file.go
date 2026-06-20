package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MartialM1nd/freefsm/internal/ent"
	"github.com/MartialM1nd/freefsm/internal/ent/file"
	"github.com/google/uuid"
)

var allowedMIMETypes = []string{
	"image/",
	"application/pdf",
	"text/",
	"application/msword",
	"application/vnd.ms-",
	"application/vnd.openxmlformats-",
	"application/zip",
	"application/json",
	"application/xml",
	"application/octet-stream",
}

type FileService struct {
	client     *ent.Client
	uploadDir  string
	maxSize    int64
}

func NewFileService(client *ent.Client, uploadDir string, maxSize int64) *FileService {
	return &FileService{client: client, uploadDir: uploadDir, maxSize: maxSize}
}

func (s *FileService) ValidateMIMEType(mimeType string) bool {
	for _, prefix := range allowedMIMETypes {
		if strings.HasPrefix(mimeType, prefix) {
			return true
		}
	}
	return false
}

func (s *FileService) MaxSize() int64 {
	return s.maxSize
}

func (s *FileService) UploadDir() string {
	return s.uploadDir
}

func (s *FileService) List(ctx context.Context, objectType string, objectID int64) ([]*ent.File, error) {
	return s.client.File.Query().
		Where(file.ObjectType(objectType), file.ObjectID(objectID)).
		Order(ent.Desc(file.FieldCreatedAt)).
		All(ctx)
}

func (s *FileService) GetByID(ctx context.Context, id int64) (*ent.File, error) {
	f, err := s.client.File.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get file %d: %w", id, err)
	}
	return f, nil
}

func (s *FileService) Create(ctx context.Context, objectType string, objectID int64, originalName string, mimeType string, fileSize int64, reader io.Reader, uploadedBy int64) (*ent.File, error) {
	if !s.ValidateMIMEType(mimeType) {
		return nil, fmt.Errorf("invalid MIME type: %s", mimeType)
	}
	if fileSize > s.maxSize {
		return nil, fmt.Errorf("file size %d exceeds maximum %d", fileSize, s.maxSize)
	}

	ext := strings.ToLower(filepath.Ext(originalName))
	storedName := uuid.New().String() + ext

	now := time.Now()
	subDir := filepath.Join(s.uploadDir, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	if err := os.MkdirAll(subDir, 0750); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}

	filePath := filepath.Join(subDir, storedName)
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file on disk: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return nil, fmt.Errorf("write file to disk: %w", err)
	}

	entFile, err := s.client.File.Create().
		SetObjectType(objectType).
		SetObjectID(objectID).
		SetOriginalName(originalName).
		SetStoredName(storedName).
		SetMimeType(mimeType).
		SetFileSize(fileSize).
		SetFilePath(filePath).
		SetUploadedBy(uploadedBy).
		Save(ctx)
	if err != nil {
		_ = os.Remove(filePath)
		return nil, fmt.Errorf("create file record: %w", err)
	}

	return entFile, nil
}

func (s *FileService) Delete(ctx context.Context, id int64) error {
	f, err := s.client.File.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("get file %d: %w", id, err)
	}

	if err := os.Remove(f.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove file from disk: %w", err)
	}

	if err := s.client.File.DeleteOneID(id).Exec(ctx); err != nil {
		return fmt.Errorf("delete file record: %w", err)
	}
	return nil
}

func (s *FileService) GetDiskPath(storedName string) string {
	return filepath.Join(s.uploadDir, storedName)
}

func IsInlineMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/") || mimeType == "application/pdf"
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
