package pkg

type FormField struct {
	Field           string  `json:"field"`
	Name            string  `json:"name"`
	Level           string  `json:"level"`
	TableShow       bool    `json:"tableShow"`
	SortOrder       float32 `json:"sortOrder"`
	Sortable        bool    `json:"sortable"`
	Editable        bool    `json:"editable"`
	Type            string  `json:"type"`
	Validate        bool    `json:"validate"`
	SearchType      string  `json:"searchType"`
	SearchLevel     string  `json:"searchLevel"`
	Searchable      bool    `json:"searchable"`
	DefaultSort     bool    `json:"defaultSort"`
	DefaultSortType string  `json:"defaultSortType"`
	DictType        string  `json:"dictType"`
	CustomUrl       string  `json:"customUrl"`
	SearchDictType  string  `json:"searchDictType"`
	SearchCustomUrl string  `json:"searchCustomUrl"`
}
