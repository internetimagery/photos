// Format functionality
package format

import "regexp"

type Media struct {
	Name      string
	Formatted bool
	Index     int
	Tags      []string
}

//
// func ParseFiles(parent string, files []string) ([]Media, error){
//   reg, err := GetRegex(parent)
//   result := []Media{}
//   for i := 0; i < len(files); i++ {
//     result.append(Media{})
//   }
//   return result, err
// }
func getRegex(dir string) (*regexp.Regexp, error) {
	// Apply prefix and suffix. Return compiled regex
	suffix := "_(\\d+)(?:\\[(.+?)])?\\.(\\w+)"
	return regexp.Compile(regexp.QuoteMeta(dir) + suffix)
}

func Match(dir string, names []string) ([]Media, error) {
	media := []Media{}
	reg, err := getRegex(dir)
	if err != nil {
		return media, err
	}
	for _, n := range names {
		TODO
	}
	return media, nil
}
