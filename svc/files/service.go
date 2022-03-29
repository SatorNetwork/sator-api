package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/storage"
	"github.com/SatorNetwork/sator-api/svc/files/repository"
)

type (
	resizerFunc func(f io.ReadCloser, w, h uint, imageType string) (io.ReadSeeker, error)

	// Service struct
	Service struct {
		msr     mediaServiceRepository
		storage *storage.Interactor
		resize  resizerFunc
	}

	File struct {
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
		GetFilesByIDs(ctx context.Context, ids []uuid.UUID) ([]repository.File, error)
		DeleteFileByID(ctx context.Context, id uuid.UUID) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(msr mediaServiceRepository, storage *storage.Interactor, resize resizerFunc) *Service {
	if msr == nil {
		log.Fatalln("media service repository is not set")
	}
	if storage == nil {
		log.Fatalln("storage interactor is not set")
	}

	return &Service{msr: msr, storage: storage, resize: resize}
}

// AddImageResize used to create new resized image.
func (s *Service) AddImageResize(ctx context.Context, it File, file multipart.File, fileHeader *multipart.FileHeader, height, width uint) (File, error) {
	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id.String(), path.Ext(fileHeader.Filename))
	filePath := s.storage.FilePath(fileName)
	ct := fileHeader.Header.Get("Content-Type")

	resizedFile, err := s.resize(file, width, height, ct)
	if err != nil {
		return File{}, errors.Wrap(err, "resize image")
	}

	if err := s.storage.Upload(resizedFile, filePath, storage.Public, ct); err != nil {
		return File{}, errors.Wrap(err, "upload image")
	}

	image, err := s.msr.AddFile(ctx, repository.AddFileParams{
		ID:       id,
		FileName: fileHeader.Filename,
		FilePath: filePath,
		FileUrl:  s.storage.FileURL(filePath),
	})
	if err != nil {
		return File{}, fmt.Errorf("could not add image with file name=%s: %w", it.Filename, err)
	}

	return castToFile(image), nil
}

// AddImage used to create new image.
func (s *Service) AddImage(ctx context.Context, it File, file io.ReadSeeker, fileHeader *multipart.FileHeader) (File, error) {
	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id.String(), path.Ext(fileHeader.Filename))
	filePath := s.storage.FilePath(fileName)
	ct := fileHeader.Header.Get("Content-Type")

	if err := s.storage.Upload(file, filePath, storage.Public, ct); err != nil {
		return File{}, errors.Wrap(err, "upload image")
	}

	image, err := s.msr.AddFile(ctx, repository.AddFileParams{
		ID:       id,
		FileName: fileHeader.Filename,
		FilePath: filePath,
		FileUrl:  s.storage.FileURL(filePath),
	})
	if err != nil {
		return File{}, fmt.Errorf("could not add image with file name=%s: %w", it.Filename, err)
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
func (s *Service) GetImageByID(ctx context.Context, id uuid.UUID) (File, error) {
	image, err := s.msr.GetFileByID(ctx, id)
	if err != nil {
		return File{}, fmt.Errorf("could not get image with id=%s: %w", id, err)
	}

	return castToFile(image), nil
}

// GetImagesList returns list images.
func (s *Service) GetImagesList(ctx context.Context, limit, offset int32) ([]File, error) {
	images, err := s.msr.GetFilesList(ctx, repository.GetFilesListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get images list: %w", err)
	}

	return castToListFiles(images), nil
}

// GetImagesList returns list images.
func (s *Service) GetImagesListByIDs(ctx context.Context, ids []uuid.UUID) ([]File, error) {
	images, err := s.msr.GetFilesByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("could not get images list: %w", err)
	}

	return castToListFiles(images), nil
}

// AddFile used to create new file.
func (s *Service) AddFile(ctx context.Context, it File, file []byte, fileHeader *multipart.FileHeader) (File, error) {
	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id.String(), path.Ext(fileHeader.Filename))
	filePath := s.storage.FilePath(fileName)
	ct := fileHeader.Header.Get("Content-Type")

	r := bytes.NewReader(file)

	if err := s.storage.Upload(r, filePath, storage.Public, ct); err != nil {
		return File{}, errors.Wrap(err, "upload image")
	}

	image, err := s.msr.AddFile(ctx, repository.AddFileParams{
		ID:       id,
		FileName: fileHeader.Filename,
		FilePath: filePath,
		FileUrl:  s.storage.FileURL(filePath),
	})
	if err != nil {
		return File{}, fmt.Errorf("could not add image with file name=%s: %w", it.Filename, err)
	}

	return castToFile(image), nil
}

// Cast repository.File to service File structure
func castToFile(source repository.File) File {
	return File{
		ID:        source.ID,
		Filename:  source.FileName,
		Filepath:  source.FilePath,
		FileUrl:   source.FileUrl,
		CreatedAt: source.CreatedAt.String(),
	}
}

// Cast repository.File to service File structure
func castToListFiles(source []repository.File) []File {
	result := make([]File, 0, len(source))
	for _, s := range source {
		result = append(result, File{
			ID:        s.ID,
			Filename:  s.FileName,
			Filepath:  s.FilePath,
			FileUrl:   s.FileUrl,
			CreatedAt: s.CreatedAt.String(),
		})
	}
	return result
}

// Cast repository.File to service File structure
func castToImagesList(source []repository.File) []string {
	result := make([]string, 0, len(source))
	for _, s := range source {
		result = append(result, s.FileUrl)
	}
	return result
}
