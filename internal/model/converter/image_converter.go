package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func ModelUploadImageRequestToEntityImage(req *model.UploadImageRequest, image *entity.Image) {
	if req == nil || image == nil {
		return
	}
}

func EntityImageToModelImageUploadedEvent(image *entity.Image, event *model.ImageUploadedEvent) {
	if image == nil || event == nil {
		return
	}
}
