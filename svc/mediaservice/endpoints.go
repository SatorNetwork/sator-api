package mediaservice

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		AddImage        endpoint.Endpoint
		GetImageByID    endpoint.Endpoint
		GetImagesList   endpoint.Endpoint
		DeleteImageByID endpoint.Endpoint
	}

	service interface {
		AddImage(ctx context.Context, it Image, file io.ReadSeeker, fileHeader *multipart.FileHeader) (Image, error)
		GetImageByID(ctx context.Context, id uuid.UUID) (Image, error)
		GetImagesList(ctx context.Context, limit, offset int32) ([]Image, error)
		DeleteImageByID(ctx context.Context, id uuid.UUID) error
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omimagepty" validate:"number,gte=0"`
		ImagesPerPage int32 `json:"images_per_page,omimagepty" validate:"number,gte=0"`
	}

	// AddImageRequest struct
	AddImageRequest struct {
		File       io.ReadSeeker
		FileHeader *multipart.FileHeader
	}
)

// Limit of images
func (r PaginationRequest) Limit() int32 {
	if r.ImagesPerPage > 0 {
		return r.ImagesPerPage
	}

	return 20
}

// Offset images
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return (r.Page - 1) * r.Limit()
	}

	return 0
}

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		AddImage:        MakeAddImageEndpoint(s, validateFunc),
		GetImageByID:    MakeGetImageByIDEndpoint(s),
		GetImagesList:   MakeGetImagesListEndpoint(s, validateFunc),
		DeleteImageByID: MakeDeleteImageByIDEndpoint(s),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddImage = mdw(e.AddImage)
			e.GetImageByID = mdw(e.GetImageByID)
			e.GetImagesList = mdw(e.GetImagesList)
			e.DeleteImageByID = mdw(e.DeleteImageByID)
		}
	}

	return e
}

// MakeGetImagesListEndpoint ...
func MakeGetImagesListEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetImagesList(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddImageEndpoint ...
func MakeAddImageEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddImageRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddImage(ctx, Image{
			Filename: req.FileHeader.Filename,
		}, req.File, req.FileHeader)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetImageByIDEndpoint ...
func MakeGetImageByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w image id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetImageByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteImageByIDEndpoint ...
func MakeDeleteImageByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get image id: %w", err)
		}

		err = s.DeleteImageByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
