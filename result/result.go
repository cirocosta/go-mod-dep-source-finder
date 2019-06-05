package result

type Result struct {
	Original           string `json:"original"`
	DiscoveredHomePage string `json:"discovered"`
	Problem            string `json:"problem,omitempty"`
}
