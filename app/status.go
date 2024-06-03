package main

type status struct {
	code int
	text string
}

type Status interface {
	Code() int
	Text() string
}

var (
	Ok                  Status = &status{200, "OK"}
	Created             Status = &status{201, "Created"}
	NotFound            Status = &status{404, "Not Found"}
	InternalServerError Status = &status{500, "Internal Server Error"}
)

func (s *status) Code() int {
	return s.code
}

func (s *status) Text() string {
	return s.text
}
