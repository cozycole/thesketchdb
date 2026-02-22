package models

type Repositories struct {
	Cast       CastModelInterface
	Categories CategoryInterface
	Characters CharacterModelInterface
	Creators   CreatorModelInterface
	Quotes     QuoteModelInterface
	People     PersonModelInterface
	Profile    ProfileModelInterface
	Recurring  RecurringModelInterface
	Shows      ShowModelInterface
	Series     SeriesModelInterface
	Tags       TagModelInterface
	Users      UserModelInterface
	Sketches   SketchModelInterface
}
