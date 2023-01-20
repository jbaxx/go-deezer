package deezer

type Track struct {
	ID                    int       `json:"id,omitempty"`
	Readable              bool      `json:"readable,omitempty"`
	Title                 string    `json:"title,omitempty"`
	TitleShort            string    `json:"title_short,omitempty"`
	TitleVersion          string    `json:"title_version,omitempty"`
	Isrc                  string    `json:"isrc,omitempty"`
	Link                  string    `json:"link,omitempty"`
	Share                 string    `json:"share,omitempty"`
	Duration              int       `json:"duration,omitempty"`
	TrackPosition         int       `json:"track_position,omitempty"`
	DiskNumber            int       `json:"disk_number,omitempty"`
	Rank                  int       `json:"rank,omitempty"`
	ReleaseDate           string    `json:"release_date,omitempty"`
	ExplicitLyrics        bool      `json:"explicit_lyrics,omitempty"`
	ExplicitContentLyrics int       `json:"explicit_content_lyrics,omitempty"`
	ExplicitContentCover  int       `json:"explicit_content_cover,omitempty"`
	Preview               string    `json:"preview,omitempty"`
	Bpm                   float64   `json:"bpm,omitempty"`
	Gain                  float64   `json:"gain,omitempty"`
	AvailableCountries    []string  `json:"available_countries,omitempty"`
	Contributors          []*Artist `json:"contributors,omitempty"`
	Md5Image              string    `json:"md5_image,omitempty"`
	Artist                *Artist   `json:"artist,omitempty"`
	Album                 *Album    `json:"album,omitempty"`
	Type                  string    `json:"type,omitempty"`
}
