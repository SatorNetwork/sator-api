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

	Item struct {
		ID        uuid.UUID `json:"id"`
		Filename  string    `json:"filename"`
		Filepath  string    `json:"filepath"`
		FileUrl   string    `json:"file_url"`
		CreatedAt string    `json:"created_at"`
	}

	mediaServiceRepository interface {
		AddItem(ctx context.Context, arg repository.AddItemParams) (repository.Item, error)
		GetItemByID(ctx context.Context, id uuid.UUID) (repository.Item, error)
		GetItemsList(ctx context.Context, arg repository.GetItemsListParams) ([]repository.Item, error)
		DeleteItemByID(ctx context.Context, id uuid.UUID) error
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

// AddItem used to create new item.
func (s *Service) AddItem(ctx context.Context, it Item, file io.ReadSeeker, fileHeader *multipart.FileHeader) (Item, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Item{}, errors.Wrap(err, "begin db transaction")
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%s%s", id.String(), path.Ext(fileHeader.Filename))
	ct := fileHeader.Header.Get("Content-Type")

	item, err := s.msr.AddItem(ctx, repository.AddItemParams{
		ID:       id,
		FileName: fileHeader.Filename,
		FilePath: s.storage.FilePath(fileName),
		FileUrl:  s.storage.FileURL(s.storage.FilePath(fileName)),
	})

	if err != nil {
		return Item{}, fmt.Errorf("could not add item with file name=%s: %w", it.Filename, err)
	}

	if err != nil {
		tx.Rollback()
		return Item{}, errors.Wrap(err, "store item to db")
	}
	if err := s.storage.Upload(file, s.storage.FilePath(fileName), storage.Public, ct); err != nil {
		tx.Rollback()
		return Item{}, errors.Wrap(err, "upload image")
	}

	if err := tx.Commit(); err != nil {
		return Item{}, errors.Wrap(err, "commit item")
	}

	return castToItem(item), nil
}

// DeleteItemByID used to delete Item by provided id.
func (s *Service) DeleteItemByID(ctx context.Context, id uuid.UUID) error {
	item, err := s.msr.GetItemByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not get item with id=%s: %w", id, err)
	}

	//fileName := fmt.Sprintf("%s%s", item.ID, path.Ext(".png"))
	err = s.storage.Remove(item.FilePath)
	if err != nil {
		return errors.Wrap(err, "could not delete item from storage")
	}

	err = s.msr.DeleteItemByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete item with id=%s:%w", id, err)
	}

	return nil
}

// GetItemByID returns item with provided id.
func (s *Service) GetItemByID(ctx context.Context, id uuid.UUID) (Item, error) {
	item, err := s.msr.GetItemByID(ctx, id)
	if err != nil {
		return Item{}, fmt.Errorf("could not get item with id=%s: %w", id, err)
	}

	return castToItem(item), nil
}

// GetItemsList returns list items.
func (s *Service) GetItemsList(ctx context.Context, limit, offset int32) ([]Item, error) {
	items, err := s.msr.GetItemsList(ctx, repository.GetItemsListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Item{}, fmt.Errorf("could not get items list: %w", err)
	}
	return castToListItems(items), nil
}

// Cast repository.Item to service Item structure
func castToItem(source repository.Item) Item {
	return Item{
		ID:        source.ID,
		Filename:  source.FileName,
		Filepath:  source.FilePath,
		FileUrl:   source.FileUrl,
		CreatedAt: source.CreatedAt.String(),
	}
}

// Cast repository.Item to service Item structure
func castToListItems(source []repository.Item) []Item {
	result := make([]Item, 0, len(source))
	for _, s := range source {
		result = append(result, Item{
			ID:        s.ID,
			Filename:  s.FileName,
			Filepath:  s.FilePath,
			FileUrl:   s.FileUrl,
			CreatedAt: s.CreatedAt.String(),
		})
	}
	return result
}
