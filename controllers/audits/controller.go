package audits

import auditsvc "sitecrawler/newgo/internal/services/audits"

type Controller struct {
	service auditsvc.Service
}

func NewController(service auditsvc.Service) *Controller {
	if service == nil {
		panic("audit service required")
	}
	return &Controller{service: service}
}
