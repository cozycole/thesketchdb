package views

import (
	"sketchdb.cozycole.net/internal/models"
)

type PaginationItem struct {
	Page       int  // the page number (real, or 0 for ellipsis)
	IsCurrent  bool // whether this is the current page
	IsEllipsis bool
	URL        string // the link to use if applicable
}

// Creates the view data necessary to render pagination for any result page
func buildPagination(current, total int, baseUrl string, filter *models.Filter) ([]*PaginationItem, error) {
	var pageItems []*PaginationItem

	pages := paginate(current, total)

	var err error
	for _, page := range pages {
		item := &PaginationItem{}
		if page == -1 {
			item.IsEllipsis = true
		} else if page == current {
			item.IsCurrent = true
			item.Page = page
		} else {
			item.Page = page
			item.URL, err = BuildURL(baseUrl, page, filter)
		}

		pageItems = append(pageItems, item)
	}

	return pageItems, err
}

func paginate(currentPage, totalPages int) []int {
	var pages []int

	// Show the current page and two pages before and after
	start := currentPage - 1
	if start < 1 {
		start = 1
	}

	end := currentPage + 1
	if end > totalPages {
		end = totalPages
	}

	// Add the main range
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}

	// Add ellipsis and the last page if necessary
	if end < totalPages {
		if end+1 < totalPages {
			pages = append(pages, -1) // -1 represents "..."
		}
		pages = append(pages, totalPages)
	}

	// Add the first page and ellipsis if necessary
	if start > 1 {
		if start > 2 {
			pages = append([]int{-1}, pages...)
		}
		pages = append([]int{1}, pages...)
	}

	return pages
}
