package domain

type (
	ParseResult struct {
		HTMLVersion       string
		Title             string
		Headings          HeadingsCounter
		Links             Links
		ContainsLoginForm bool
	}

	ErrorResult struct {
		Error error
	}

	HeadingsCounter struct {
		H1 int
		H2 int
		H3 int
		H4 int
		H5 int
		H6 int
	}

	Links struct {
		Internal LinksInfo
		External LinksInfo
	}

	LinksInfo struct {
		Count        int
		Accessible   AccessibleLinks
		Inaccessible InaccessibleLinks
	}

	AccessibleLinks struct {
		Links []string
		Count int
	}

	InaccessibleLinks struct {
		Links []string
		Count int
	}
)
