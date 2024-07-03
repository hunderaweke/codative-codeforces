package types

type Response struct {
	Status string `json:"status,omitempty"`
	Result Result `json:"result,omitempty"`
}
type Result struct {
	Contest  Contest   `json:"contest,omitempty"`
	Problems []Problem `json:"problems,omitempty"`
}

type Contest struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Phase  string `json:"phase,omitempty"`
	Frozen bool   `json:"frozen,omitempty"`
}
type Problem struct {
	Index string `json:"index,omitempty"`
	Name  string `json:"name,omitempty"`
}
