package domain

type (
	AnalyzeResult struct {
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
		Internal     LinksInfo
		External     LinksInfo
		Inaccessible int
	}

	LinksInfo struct {
		Count int
	}
)
