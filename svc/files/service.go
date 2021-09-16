package files

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path"

	"github.com/SatorNetwork/sator-api/internal/storage"
	"github.com/SatorNetwork/sator-api/svc/files/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	resizerFunc func(f io.ReadCloser, w, h int) (io.ReadSeeker, error)

	// Service struct
	Service struct {
		msr     mediaServiceRepository
		db      *sql.DB
		storage *storage.Interactor
		resize  resizerFunc
	}

	Image struct {
		ID        uuid.UUID `json:"id"`
		Filename  string    `json:"filename"`
		Filepath  string    `json:"filepath"`
		FileUrl   string    `json:"file_url"`
		CreatedAt string    `json:"created_at"`
	}

	mediaServiceRepository interface {
		AddFile(ctx context.Context, arg repository.AddFileParams) (repository.File, error)
		GetFileByID(ctx context.Context, id uuid.UUID) (repository.File, error)
		GetFilesList(ctx context.Context, arg repository.GetFilesListParams) ([]repository.File, error)
		DeleteFileByID(ctx context.Context, id uuid.UUID) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(msr mediaServiceRepository, db *sql.DB, storage *storage.Interactor, resize resizerFunc) *Service {
	if msr == nil {
		log.Fatalln("media service repository is not set")
	}
	if db == nil {
		log.Fatalln("db is not set")
	}
	if storage == nil {
		log.Fatalln("storage interactor is not set")
	}

	return &Service{msr: msr, db: db, storage: storage, resize: resize}
}

// AddImage used to create new image.
func (s *Service) AddImage(ctx context.Context, it Image, file multipart.File, fileHeader *multipart.FileHeader, height, width int) (Image, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Image{}, errors.Wrap(err, "begin db transaction")
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id.String(), path.Ext(fileHeader.Filename))
	ct := fileHeader.Header.Get("Content-Type")

	image, err := s.msr.AddFile(ctx, repository.AddFileParams{
		ID:       id,
		FileName: fileHeader.Filename,
		FilePath: s.storage.FilePath(fileName),
		FileUrl:  s.storage.FileURL(s.storage.FilePath(fileName)),
	})
	if err != nil {
		return Image{}, fmt.Errorf("could not add image with file name=%s: %w", it.Filename, err)
	}

	resizedFile, err := s.resize(file, width, height)
	if err != nil {
		tx.Rollback()
		return Image{}, errors.Wrap(err, "resize image")
	}

	if err != nil {
		tx.Rollback()
		return Image{}, errors.Wrap(err, "store image to db")
	}
	if err := s.storage.Upload(resizedFile, s.storage.FilePath(fileName), storage.Public, ct); err != nil {
		tx.Rollback()
		return Image{}, errors.Wrap(err, "upload image")
	}

	if err := tx.Commit(); err != nil {
		return Image{}, errors.Wrap(err, "commit image")
	}

	return castToFile(image), nil
}

// DeleteImageByID used to delete File by provided id.
func (s *Service) DeleteImageByID(ctx context.Context, id uuid.UUID) error {
	image, err := s.msr.GetFileByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not get image with id=%s: %w", id, err)
	}

	//fileName := fmt.Sprintf("%s%s", image.ID, path.Ext(".png"))
	err = s.storage.Remove(image.FilePath)
	if err != nil {
		return errors.Wrap(err, "could not delete image from storage")
	}

	err = s.msr.DeleteFileByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete image with id=%s:%w", id, err)
	}

	return nil
}

// GetImageByID returns image with provided id.
func (s *Service) GetImageByID(ctx context.Context, id uuid.UUID) (Image, error) {
	image, err := s.msr.GetFileByID(ctx, id)
	if err != nil {
		return Image{}, fmt.Errorf("could not get image with id=%s: %w", id, err)
	}

	return castToFile(image), nil
}

// GetImagesList returns list images.
func (s *Service) GetImagesList(ctx context.Context, limit, offset int32) ([]Image, error) {
	images, err := s.msr.GetFilesList(ctx, repository.GetFilesListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get images list: %w", err)
	}
	return castToListFiles(images), nil
}

// Cast repository.File to service File structure
func castToFile(source repository.File) Image {
	return Image{
		ID:        source.ID,
		Filename:  source.FileName,
		Filepath:  source.FilePath,
		FileUrl:   source.FileUrl,
		CreatedAt: source.CreatedAt.String(),
	}
}

// Cast repository.File to service File structure
func castToListFiles(source []repository.File) []Image {
	result := make([]Image, 0, len(source))
	for _, s := range source {
		result = append(result, Image{
			ID:        s.ID,
			Filename:  s.FileName,
			Filepath:  s.FilePath,
			FileUrl:   s.FileUrl,
			CreatedAt: s.CreatedAt.String(),
		})
	}
	return result
}
