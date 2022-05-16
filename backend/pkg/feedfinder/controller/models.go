package controller

// FeedFinderRequest gives a website HTTP URL whose feeds should be extracted
type FeedFinderRequest struct {
	URL string `json:"url"`
}

// FeedFinderResponse contains the list of feeds we were able to find for the given URL
type FeedFinderResponse struct {
	OK           bool   `json:"ok"`
	ErrorMessage string `json:"error_message,omitempty"`
	URL          string `json:"url,omitempty"`
	Feeds        []Feed `json:"feeds,omitempty"`
}

type Feed struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
