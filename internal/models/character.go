package models

type Character struct {
	ID          *int
	Name        *string
	Image       *string
	Description *string
	Portrayal   *Person
}
