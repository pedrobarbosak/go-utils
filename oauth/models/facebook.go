package models

type Facebook struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Picture Picture `json:"picture"`
}

type Picture struct {
	Data Data `json:"data"`
}

type Data struct {
	Height       uint   `json:"height"`
	Width        uint   `json:"width"`
	IsSilhouette bool   `json:"is_silhouette"`
	URL          string `json:"url"`
}
