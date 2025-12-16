package services

import (
	"context"
	"errors"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type pageDetailsService struct {
	repo repository.PageDetailsRepository
}

func NewPageDetailsService(repo repository.PageDetailsRepository) PageDetailsService {
	if repo == nil {
		panic("page details repository required")
	}
	return &pageDetailsService{repo: repo}
}

func (s *pageDetailsService) Details(ctx context.Context, req dto.PageDetailsRequest) (*dto.Response[dto.PageDetailsResponse], error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	page, err := s.repo.GetPageByIDAndSKU(ctx, req.PageID, req.SearchKeywordURLID)
	if err != nil {
		if errors.Is(err, repository.ErrPageNotFound) {
			return dto.NewResponse[dto.PageDetailsResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[dto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	data := dto.PageDetailsData{}
	if page.ResponseCode >= 200 && page.ResponseCode <= 299 {
		images, err := s.repo.GetPageImages(ctx, req.PageID, limit)
		if err != nil {
			return dto.NewResponse[dto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
		}
		broken, err := s.repo.GetBrokenTargetsFrom(ctx, req.PageID, limit)
		if err != nil {
			return dto.NewResponse[dto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
		}
		data.PageImages = images
		data.BrokenPages = broken
	} else {
		referrers, err := s.repo.GetReferrersToBroken(ctx, req.PageID, limit)
		if err != nil {
			return dto.NewResponse[dto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
		}
		data.PagesLinkedToBrokenPage = referrers
	}

	return dto.NewSuccessResponse(dto.PageDetailsResponse{Data: data}, http.StatusOK), nil
}
