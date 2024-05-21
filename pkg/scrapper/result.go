package scrapper

type (
	ParseResult struct {
		HTMLVersion string
		Title       string
		Headings    Headings
		Links       []string
		Forms       []Form
	}

	ErrorResult struct {
		Error error
	}

	Headings struct {
		H1 []string
		H2 []string
		H3 []string
		H4 []string
		H5 []string
		H6 []string
	}

	Form struct {
		Action string
		Inputs []Input
	}

	Input struct {
		Type string
	}
)
