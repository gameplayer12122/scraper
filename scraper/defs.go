package scraper

type jsonPosts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID   uint64 `json:"id"`
	File File   `json:"file"`
}

type File struct {
	URL       string `json:"url"`
	Extension string `json:"ext"`
	MD5       string `json:"md5"`
}

type r34JsonPost struct { // Transformed into a Post and a File eventually
	ID      uint64 `json:"id"`
	FileURL string `json:"file_url"`
	FileMD5 string `json:"hash"`
}

type gelBooruJsonPosts struct { // Gelbooru's API response is the samme as danbooru's response, except without extension.
	Posts []danBooruJsonPost `json:"post"`
}

type danBooruJsonPost struct { // Transformed into a Post and a File eventually
	ID      uint64 `json:"id"`
	FileExt string `json:"file_ext"`
	FileURL string `json:"file_url"`
	FileMD5 string `json:"md5"`
}
