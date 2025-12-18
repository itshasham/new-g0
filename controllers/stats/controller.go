package stats

import statssvc "sitecrawler/newgo/internal/services/stats"

type Controller struct {
	service statssvc.Service
}

func NewController(service statssvc.Service) *Controller {
	if service == nil {
		panic("stats service required")
	}
	return &Controller{service: service}
}
