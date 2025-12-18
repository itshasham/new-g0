package stats

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/dto"

	statsDto "sitecrawler/newgo/dto/stats"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) Details(ctx context.Context, req statsDto.PageDetailsRequest) (*dto.Response[statsDto.PageDetailsResponse], error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	page, err := s.pageDetailsRepo.GetPageByIDAndSKU(ctx, req.PageID, req.SearchKeywordURLID)
	if err != nil {
		if errors.Is(err, repository.ErrPageNotFound) {
			return dto.NewResponse[statsDto.PageDetailsResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[statsDto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	data := statsDto.PageDetailsData{}
	if page.ResponseCode >= 200 && page.ResponseCode <= 299 {
		images, err := s.pageDetailsRepo.GetPageImages(ctx, req.PageID, limit)
		if err != nil {
			return dto.NewResponse[statsDto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
		}
		broken, err := s.pageDetailsRepo.GetBrokenTargetsFrom(ctx, req.PageID, limit)
		if err != nil {
			return dto.NewResponse[statsDto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
		}
		data.PageImages = images
		data.BrokenPages = broken
	} else {
		referrers, err := s.pageDetailsRepo.GetReferrersToBroken(ctx, req.PageID, limit)
		if err != nil {
			return dto.NewResponse[statsDto.PageDetailsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
		}
		data.PagesLinkedToBrokenPage = referrers
	}

	return dto.NewSuccessResponse(statsDto.PageDetailsResponse{Data: data}, http.StatusOK), nil
}
