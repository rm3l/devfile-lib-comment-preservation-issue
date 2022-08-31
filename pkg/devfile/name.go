package devfile

import (
	"github.com/devfile/library/pkg/devfile/parser"
	"k8s.io/utils/pointer"
)

// UpdateName updates the 'metadata.name' field in the Devfile at the provided path and writes it to disk.
func UpdateName(p string, n string) error {
	d, err := parser.ParseDevfile(parser.ParserArgs{Path: p, FlattenedDevfile: pointer.Bool(false)})
	if err != nil {
		return err
	}
	err = d.SetMetadataName(n)
	if err != nil {
		return err
	}
	return d.WriteYamlDevfile()
}
