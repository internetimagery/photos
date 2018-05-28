// Media representation

package media

import (
	"path/filepath"

	"github.com/internetimagery/photos/utility"
)

// TODO: Consider file format: group.index.hash[tag tag].type
// TODO: stuck using [tag] format due to software

// thoughts...
// NewMediaFromFile(path)
// media object, path to associated file, name, path, id etc...
// create media object, parsing filename for metadata
// metadata can be changed to whatever. Media.Update() to rename it perhaps

// Typical format:
// GROUP_INDEX_HASH[TAGS TAGS].TYPE
type Media struct {
	path  string   // Path to file.
	index int      // Unique to folder. Sorting files.
	group string   // Group media is located in. Typically the name of the containing folder.
	_type string   // File type. Doubles as file extension.
	hash  string   // Hash fingerprint of files content.
	tags  []string // Extra metadata.
}

func (self *Media) GetPath() string {
	return self.path
}

func (self *Media) GetIndex() int {
	return self.index
}

func (self *Media) SetIndex(index int) {
	self.index = index
}

func (self *Media) GetGroup() string {
	return self.group
}

func (self *Media) SetGroup(group int) {
	// TODO: Sanitize input here
	self.group = group
}

func (self *Media) GetType() string {
	return self.group
}

func (self *Media) SetType(_type int) {
	// TODO: Sanitize input here
	self._type = _type
}

func (self *Media) GetHash() string {
	return self.hash
}

func (self *Media) GetTags() string {
	return self.tags
}

func (self *Media) SetTags(tags int) {
	// TODO: Sanitize input here
	self.tags = tags
}

// Get new media element from file.
func NewMedia(filename string) (*Media, error) {
	media := new(Media)
	media.path = filename
	media._type = filepath.Ext(filename)
	Fhash, err := utility.GetHashFromFile(filename, "md5")
	if err != nil {
		return nil, err
	}
	media := &Media{
		path: filename,
		_type: filepath.Ext(filename)
		hash: Fhash
	}
	// TODO: regex group, index, tags
	return media, nil
}
