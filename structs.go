package dca

// Base metadata struct
// 
// https://github.com/bwmarrin/dca/issues/5#issuecomment-189713886
type MetadataStruct struct {
    Dca             *DCAMetadata    `json:"dca"`
    SongInfo        *SongMetadata   `json:"info"`
    Origin          *OriginMetadata `json:"origin"`
    Opus            *OpusMetadata   `json:"opus"`

    ModifiedDate    int64           `json:"modified_date"`
    CreationDate    int64           `json:"creation_date"`
}

// DCA metadata struct
// 
// Contains the DCA version.
type DCAMetadata struct {
    Version int8                `json:"version"`
    Tool    *DCAToolMetadata    `json:"tool"`
}

// DCA tool metadata struct
// 
// Contains the Git revisions, commit author etc.
type DCAToolMetadata struct {
    Name        string  `json:"name"`
    Version     string  `json:"version"`
    Revision    string  `json:"rev"`
    Url         string  `json:"url"`
    Author      string  `json:"author"`
}

// Song Information metadata struct
// 
// Contains information about the song that was encoded.
type SongMetadata struct {
    Title       string  `json:"title"`
    Artist      string  `json:"artist"`
    Album       string  `json:"album"`
    Genre       string  `json:"genre"`
    Comments    string  `json:"comments"`
}

// Origin information metadata struct
// 
// Contains information about where the song came from,
// audio bitrate, channels and original encoding.
type OriginMetadata struct {
    Source      string  `json:"source"`
    Bitrate     int     `json:"bitrate"`
    Channels    int     `json:"channels"`
    Encoding    string  `json:"encoding"`
    Url         string  `json:"url"`
}

// Opus metadata struct
// 
// Contains information about how the file was encoded
// with Opus.
type OpusMetadata struct {
    SampleRate  int     `json:"sample_rate"`
    Application string  `json:"mode"`
    FrameSize   int     `json:"frame_size"`
    Channels    int     `json:"channels"`
}