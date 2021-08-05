package mediaservice

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path"

	"github.com/SatorNetwork/sator-api/internal/storage"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	// Service struct
	Service struct {
		msr     mediaServiceRepository
		db      *sql.DB
		storage *storage.Interactor
	}

	Image struct {
		ID        uuid.UUID `json:"id"`
		Filename  string    `json:"filename"`
		Filepath  string    `json:"filepath"`
		FileUrl   string    `json:"file_url"`
		CreatedAt string    `json:"created_at"`
	}

	mediaServiceRepository interface {
		AddImage(ctx context.Context, arg repository.AddImageParams) (repository.Image, error)
		GetImageByID(ctx context.Context, id uuid.UUID) (repository.Image, error)
		GetImagesList(ctx context.Context, arg repository.GetImagesListParams) ([]repository.Image, error)
		DeleteImageByID(ctx context.Context, id uuid.UUID) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(msr mediaServiceRepository, db *sql.DB, storage *storage.Interactor) *Service {
	if msr == nil {
		log.Fatalln("media service repository is not set")
	}
	if db == nil {
		log.Fatalln("db is not set")
	}
	if storage == nil {
		log.Fatalln("storage interactor is not set")
	}

	return &Service{msr: msr, db: db, storage: storage}
}

// AddImage used to create new image.
func (s *Service) AddImage(ctx context.Context, it Image, file io.ReadSeeker, fileHeader *multipart.FileHeader) (Image, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Image{}, errors.Wrap(err, "begin db transaction")
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id.String(), path.Ext(fileHeader.Filename))
	ct := fileHeader.Header.Get("Content-Type")

	image, err := s.msr.AddImage(ctx, repository.AddImageParams{
		ID:       id,
		FileName: fileHeader.Filename,
		FilePath: s.storage.FilePath(fileName),
		FileUrl:  s.storage.FileURL(s.storage.FilePath(fileName)),
	})

	if err != nil {
		return Image{}, fmt.Errorf("could not add image with file name=%s: %w", it.Filename, err)
	}

	if err != nil {
		tx.Rollback()
		return Image{}, errors.Wrap(err, "store image to db")
	}
	if err := s.storage.Upload(file, s.storage.FilePath(fileName), storage.Public, ct); err != nil {
		tx.Rollback()
		return Image{}, errors.Wrap(err, "upload image")
	}

	if err := tx.Commit(); err != nil {
		return Image{}, errors.Wrap(err, "commit image")
	}

	return castToImage(image), nil
}

// DeleteImageByID used to delete Image by provided id.
func (s *Service) DeleteImageByID(ctx context.Context, id uuid.UUID) error {
	image, err := s.msr.GetImageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not get image with id=%s: %w", id, err)
	}

	//fileName := fmt.Sprintf("%s%s", image.ID, path.Ext(".png"))
	err = s.storage.Remove(image.FilePath)
	if err != nil {
		return errors.Wrap(err, "could not delete image from storage")
	}

	err = s.msr.DeleteImageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete image with id=%s:%w", id, err)
	}

	return nil
}

// GetImageByID returns image with provided id.
func (s *Service) GetImageByID(ctx context.Context, id uuid.UUID) (Image, error) {
	image, err := s.msr.GetImageByID(ctx, id)
	if err != nil {
		return Image{}, fmt.Errorf("could not get image with id=%s: %w", id, err)
	}

	return castToImage(image), nil
}

// GetImagesList returns list images.
func (s *Service) GetImagesList(ctx context.Context, limit, offset int32) ([]Image, error) {
	images, err := s.msr.GetImagesList(ctx, repository.GetImagesListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Image{}, fmt.Errorf("could not get images list: %w", err)
	}
	return castToListImages(images), nil
}

// Cast repository.Image to service Image structure
func castToImage(source repository.Image) Image {
	return Image{
		ID:        source.ID,
		Filename:  source.FileName,
		Filepath:  source.FilePath,
		FileUrl:   source.FileUrl,
		CreatedAt: source.CreatedAt.String(),
	}
}

// Cast repository.Image to service Image structure
func castToListImages(source []repository.Image) []Image {
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
