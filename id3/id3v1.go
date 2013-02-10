package id3

import (
	"fmt"
	"strings"
)

type ID3v1Tag struct {
	Title   string
	Artist  string
	Album   string
	Year    string
	Comment string
	Track   string
	Genre   string
}

func trimString(data []byte) string {
	return strings.TrimRight(string(data), "\u0000")
}

func getGenre(i byte) string {
	if int(i) > len(ID3v1Genres)-1 {
		return "Unspecified"
	}
	return ID3v1Genres[i]
}

func ParseID3v1Tag(data []byte) (*ID3v1Tag, error) {
	if string(data[0:3]) != "TAG" {
		return nil, ErrNoTag
	}
	tag := new(ID3v1Tag)
	tag.Title = trimString(data[3:33])
	tag.Artist = trimString(data[33:63])
	tag.Album = trimString(data[63:93])
	tag.Year = trimString(data[93:97])
	if data[125] == 0 && data[126] != 0 {
		tag.Track = fmt.Sprint(data[126])
		tag.Comment = trimString(data[97:125])
	} else {
		tag.Comment = trimString(data[97:127])
	}
	tag.Genre = getGenre(data[127])
	return tag, nil
}
