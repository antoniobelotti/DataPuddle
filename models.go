package main

type SessionKeyReponse struct {
	Outcome string `json:"outcome"`
	Key     string `json:"key"`
}

type PWDResponse struct {
	Outcome string `json:"outcome"`
	Path    string `json:"path"`
}

type OutcomeResponse struct {
	Outcome string `json:"outcome"`
}
