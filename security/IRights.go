package security

type ReadingRights interface {
	GetReadingRights() AuthorizationFunction
}
type WritingRights interface {
	GetWritingRights() AuthorizationFunction
}

func HasReadingRights(obj any) (ReadingRights, bool) {
	rr, ok := obj.(ReadingRights)
	return rr, ok
}

func HasWritingRights(obj any) (WritingRights, bool) {
	rr, ok := obj.(WritingRights)
	return rr, ok
}
