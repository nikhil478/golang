package models

type Fact struct {
	ID         string
	Type       string
	Value      string
	Confidence float64

	SourceFile string
	SourceLine int
}
