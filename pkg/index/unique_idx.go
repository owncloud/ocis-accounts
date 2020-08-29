package index

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Unique struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string
	backlinkDir  string
}

// NewUniqueIndex instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewUniqueIndex(typeName, indexBy, filesDir, indexBaseDir string) Unique {
	return Unique{
		indexBy:      indexBy,
		typeName:     typeName,
		filesDir:     filesDir,
		indexBaseDir: indexBaseDir,
		indexRootDir: path.Join(indexBaseDir, fmt.Sprintf("%sBy%s", typeName, indexBy)),
		backlinkDir:  path.Join(indexBaseDir, fmt.Sprintf("%sBacklinks", typeName)),
	}
}

func (idx Unique) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.backlinkDir, 0777); err != nil {
		return err
	}

	return nil
}

func (idx Unique) Add(pk, v string) error {
	oldName := path.Join(idx.filesDir, pk)
	newName := path.Join(idx.indexRootDir, v)
	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return &alreadyExistsErr{idx.typeName, idx.indexBy, v}
	}

	blPath := path.Join(idx.backlinkDir, pk)
	if err := os.MkdirAll(blPath, 0777); err != nil {
		return err
	}

	return os.Symlink(newName, path.Join(blPath, v))
}

func (idx Unique) Remove(v string) (err error) {
	searchPath := path.Join(idx.indexRootDir, v)
	if err = isValidSymlink(searchPath); err != nil {
		return
	}

	return os.Remove(searchPath)
}

func (idx Unique) Lookup(v string) (resultPath string, err error) {
	searchPath := path.Join(idx.indexRootDir, v)
	if err = isValidSymlink(searchPath); err != nil {
		if os.IsNotExist(err) {
			err = &notFoundErr{idx.typeName, idx.indexBy, v}
		}

		return
	}

	return os.Readlink(searchPath)
}

func (idx Unique) Update(oldV, newV string) (err error) {
	oldPath := path.Join(idx.indexRootDir, oldV)
	if err = isValidSymlink(oldPath); err != nil {
		if os.IsNotExist(err) {
			return &notFoundErr{idx.typeName, idx.indexBy, oldV}
		}

		return
	}

	newPath := path.Join(idx.indexRootDir, newV)
	if err = isValidSymlink(newPath); err == nil {
		return &alreadyExistsErr{idx.typeName, idx.indexBy, newV}
	}

	if os.IsNotExist(err) {
		err = os.Rename(oldPath, newPath)
	}

	return
}

func (idx Unique) IndexBy() string {
	return idx.indexBy
}

func (idx Unique) TypeName() string {
	return idx.typeName
}

func (idx Unique) FilesDir() string {
	return idx.filesDir
}

func (idx Unique) BacklinksDir() string {
	return idx.backlinkDir

}

func isValidSymlink(path string) (err error) {
	var symInfo os.FileInfo
	if symInfo, err = os.Lstat(path); err != nil {
		return
	}

	if symInfo.Mode()&os.ModeSymlink == 0 {
		err = fmt.Errorf("%s is not a valid symlink (bug/corruption?)", path)
		return
	}

	return

}
