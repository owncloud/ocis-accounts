package index

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Unique struct {
	indexBy  string
	typeName string
	filesDir string
	indexDir string
}

func (idx Unique) Init() error {
	indexDirName := fmt.Sprintf("%sBy%s", idx.TypeName(), idx.IndexBy())
	err := os.MkdirAll(path.Join(idx.indexDir, indexDirName), 0777)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(idx.filesDir), 0777)
	if err != nil {
		return err
	}

	return nil
}

func (idx Unique) IndexBy() string {
	return idx.indexBy
}

func (idx Unique) TypeName() string {
	return idx.typeName
}

func (idx Unique) Add(pk, v string) error {
	indexDirName := fmt.Sprintf("%sBy%s", idx.TypeName(), idx.IndexBy())
	oldName := path.Join(idx.filesDir, pk)
	newName := path.Join(idx.indexDir, indexDirName, v)

	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return nil
}

func (idx Unique) Lookup(v string) (string, error) {
	indexDirName := fmt.Sprintf("%sBy%s", idx.TypeName(), idx.IndexBy())
	searchPath := path.Join(idx.indexDir, indexDirName, v)
	info, err := os.Lstat(searchPath)
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		pk, err := os.Readlink(searchPath)
		if err != nil {
			return "", err
		}

		return path.Base(pk), nil
	}

	return "", errors.New("not found")

}

func (idx Unique) EntryRemoved(pk, v string) {

}
