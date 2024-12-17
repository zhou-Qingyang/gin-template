package request

type ApisParams struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
	ApiGroup    string `json:"api_group"`
	PageInfo
}

type ApiCommonID struct {
	ID int `json:"id"`
}

type ApiNewInfo struct {
	ID          int    `json:"id"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
	ApiGroup    string `json:"api_group"`
}
