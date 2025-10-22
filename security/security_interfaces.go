package security

// ReadingRights is an interface for objects that have reading rights
// it should return an AuthorizationFunction
type ReadingRights interface {
	GetReadingRights() AuthorizationFunction
}

// WritingRights is an interface for objects that have writing rights
// it should return an AuthorizationFunction
type WritingRights interface {
	GetWritingRights() AuthorizationFunction
}

// HasReadingRights checks if the object implements the ReadingRights interface.
func HasReadingRights(obj any) (ReadingRights, bool) {
	rr, ok := obj.(ReadingRights)
	return rr, ok
}

// HasWritingRights checks if the object implements the WritingRights interface.
func HasWritingRights(obj any) (WritingRights, bool) {
	rr, ok := obj.(WritingRights)
	return rr, ok
}
