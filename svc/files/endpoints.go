package files

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/utils"
	"github.com/SatorNetwork/sator-api/lib/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		AddFile         endpoint.Endpoint
		AddImage        endpoint.Endpoint
		GetImageByID    endpoint.Endpoint
		GetImagesList   endpoint.Endpoint
		DeleteImageByID endpoint.Endpoint
	}

	service interface {
		AddFile(ctx context.Context, it File, file []byte, fileHeader *multipart.FileHeader) (File, error)
		AddImage(ctx context.Context, it File, file io.ReadSeeker, fileHeader *multipart.FileHeader) (File, error)
		AddImageResize(ctx context.Context, it File, file multipart.File, fileHeader *multipart.FileHeader, maxHeight, maxWidth uint) (File, error)
		GetImageByID(ctx context.Context, id uuid.UUID) (File, error)
		GetImagesList(ctx context.Context, limit, offset int32) ([]File, error)
		DeleteImageByID(ctx context.Context, id uuid.UUID) error
	}

	// AddImageResizeRequest struct
	AddImageResizeRequest struct {
		File       multipart.File
		FileHeader *multipart.FileHeader
		MaxHeight  uint
		MaxWidth   uint
	}

	// AddFileRequest struct
	AddFileRequest struct {
		File       []byte
		FileHeader *multipart.FileHeader
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		AddFile:         MakeAddFileEndpoint(s, validateFunc),
		AddImage:        MakeAddImageResizeEndpoint(s, validateFunc),
		GetImageByID:    MakeGetImageByIDEndpoint(s),
		GetImagesList:   MakeGetImagesListEndpoint(s, validateFunc),
		DeleteImageByID: MakeDeleteImageByIDEndpoint(s),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddFile = mdw(e.AddFile)
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(utils.PaginationRequest)
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

// MakeAddImageResizeEndpoint ...
func MakeAddImageResizeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(AddImageResizeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddImageResize(ctx, File{
			Filename: req.FileHeader.Filename,
		}, req.File, req.FileHeader, req.MaxHeight, req.MaxWidth)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddFileEndpoint ...
func MakeAddFileEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(AddFileRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddFile(ctx, File{
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
