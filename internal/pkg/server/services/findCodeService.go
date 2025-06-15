package services

type FindCodeService interface {
}

type FindCodeServiceImpl struct {
}

func NewFindCodeService() FindCodeService {
	return &FindCodeServiceImpl{}
}
