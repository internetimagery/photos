// Media representation

package media

import (
	"path/filepath"

	"github.com/internetimagery/photos/utility"
)

// thoughts...
// NewMediaFromFile(path)
// media object, path to associated file, name, path, id etc...
// create media object, parsing filename for metadata
// metadata can be changed to whatever. Media.Update() to rename it perhaps

// Typical format:
// GROUP_INDEX_HASH[TAGS TAGS].TYPE
type Media struct {
	Path  string   // Path to file.
	Index int      // Unique to folder. Sorting files.
	Group string   // Group media is located in. Typically the name of the containing folder.
	Type  string   // File type. Doubles as file extension.
	Hash  string   // Hash fingerprint of files content.
	Tags  []string // Extra metadata
}

// Get new media element from file.
func NewMedia(filename string) (*Media, error) {
	media := new(Media)
	media.Path = filename
	media.Type = filepath.Ext(filename)
	Fhash, err := utility.GetHashFromFile(filename, "md5")
	if err != nil {
		return media, err
	}
	media.Hash = Fhash
	// TODO: regex group, index, tags
	return media, nil
}
