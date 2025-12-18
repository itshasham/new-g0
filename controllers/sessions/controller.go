package sessions

import sessionsvc "sitecrawler/newgo/internal/services/sessions"

type Controller struct {
	service sessionsvc.Service
}

func NewController(service sessionsvc.Service) *Controller {
	if service == nil {
		panic("crawling session service required")
	}
	return &Controller{service: service}
}
