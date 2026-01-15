package repository

var _ ImageRepository = &ImageRepositoryMwLogger{}

type ImageRepositoryMwLogger struct {
	Next ImageRepository
}

func NewImageRepositoryMwLogger(next ImageRepository) *ImageRepositoryMwLogger {
	return &ImageRepositoryMwLogger{
		Next: next,
	}
}
