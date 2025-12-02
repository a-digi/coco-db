package testmeta

type TestMeta struct {
	Indexes []struct {
		Name   string   `json:"name"`
		Fields []string `json:"fields"`
	} `json:"indexes"`
	TotalEntries int `json:"totalEntries"`
}

