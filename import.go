// Package zipimport reads a zip archive and publishes the plain text files
// within it to a WriteFreely instance.
package zipimport

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/writeas/go-writeas"
	"github.com/writeas/web-core/posts"
)

func Parse(fname string) error {
	zf, err := zip.OpenReader(fname)
	if err != nil {
		return err
	}
	defer zf.Close()

	for _, f := range zf.File {
		// A trailing slash means this is an empty directory
		isEmptyDir := strings.HasSuffix(f.Name, "/")
		if isEmptyDir {
			// TODO: add this filename as a collection
			continue
		}

		// Get directory, slug, etc. from the filename
		fParts := strings.Split(f.Name, "/")
		var collAlias string
		var postFname string
		if len(fParts) == 1 {
			// This is a top-level file
			postFname = fParts[0]
		} else {
			// This is a collection post
			collAlias = fParts[0]
			postFname = fParts[1]
		}

		// Get file contents
		fc, err := f.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file failed: %s: %v", f.Name, err)
			continue
		}
		defer fc.Close()
		content, err := ioutil.ReadAll(fc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read file failed: %s: %v", f.Name, err)
			continue
		}

		// Build post parameters
		p := filenameToParams(postFname)
		p.Created = &f.Modified
		p.Title, p.Content = posts.ExtractTitle(string(content))

		fmt.Printf("%s - %s - %+v\n", f.Name, collAlias, p)
	}
	return nil
}

// filenameToParams returns PostParams with the ID and slug derived from the given filename.
func filenameToParams(fname string) *writeas.PostParams {
	baseParts := strings.Split(fname, ".")
	// This assumes there's at least one '.' in the filename, e.g. abc123.txt
	// TODO: handle the case where len(baseParts) != 2
	baseName := baseParts[0]

	p := &writeas.PostParams{}

	parts := strings.Split(baseName, "_")
	if len(parts) == 1 {
		// There's no slug -- only an ID
		p.ID = parts[0]
	} else {
		// len(parts) > 1
		p.Slug = parts[0]
		p.ID = parts[1]
	}
	return p
}
