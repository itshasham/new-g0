package views

import viewsvc "sitecrawler/newgo/internal/services/views"

type Controller struct {
	service viewsvc.Service
}

func NewController(service viewsvc.Service) *Controller {
	if service == nil {
		panic("view service required")
	}
	return &Controller{service: service}
}
