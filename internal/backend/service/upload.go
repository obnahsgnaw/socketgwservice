package service

import (
	"github.com/obnahsgnaw/socketgwservice/internal/service/fileasset"
	commonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/common/v1"
)

func ToFileForm(item *commonv1.FileForm) *fileasset.FileForm {
	if item == nil {
		return nil
	}
	s := &fileasset.FileForm{
		SessionId: item.SessionId,
	}
	for _, i := range item.Items {
		s.Items = append(s.Items, &fileasset.FileItem{
			UploadId: i.UploadId,
			Name:     i.Name,
		})
	}
	return s
}

func ToFile(f *fileasset.File) *commonv1.File {
	return &commonv1.File{
		Name:    f.Name,
		ViewUrl: f.ViewUrl,
	}
}
