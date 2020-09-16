package index

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// NormalIndex is able to index an document by a key which might contain non-unique values
//
// /var/tmp/testfiles-395764020/index.disk/PetByColor/
// ├── Brown
// │   └── rebef-123 -> /var/tmp/testfiles-395764020/pets/rebef-123
// ├── Green
// │    ├── goefe-789 -> /var/tmp/testfiles-395764020/pets/goefe-789
// │    └── xadaf-189 -> /var/tmp/testfiles-395764020/pets/xadaf-189
// └── White
//     └── wefwe-456 -> /var/tmp/testfiles-395764020/pets/wefwe-456
type NormalIndex struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string
}

// NewNormalIndex instantiates a new NormalIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewNormalIndex(typeName, indexBy, filesDir, indexBaseDir string) NormalIndex {
	return NormalIndex{
		indexBy:      indexBy,
		typeName:     typeName,
		filesDir:     filesDir,
		indexBaseDir: indexBaseDir,
		indexRootDir: path.Join(indexBaseDir, fmt.Sprintf("%sBy%s", typeName, indexBy)),
	}
}

func (idx NormalIndex) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

func (idx NormalIndex) Lookup(v string) ([]string, error) {
	searchPath := path.Join(idx.indexRootDir, v)
	fi, err := ioutil.ReadDir(searchPath)
	if os.IsNotExist(err) {
		return []string{}, &notFoundErr{idx.typeName, idx.indexBy, v}
	}

	if err != nil {
		return []string{}, err
	}

	var ids []string = nil
	for _, f := range fi {
		ids = append(ids, f.Name())
	}

	if len(ids) == 0 {
		return []string{}, &notFoundErr{idx.typeName, idx.indexBy, v}
	}

	return ids, nil
}

func (idx NormalIndex) Add(id, v string) (string, error) {
	oldName := path.Join(idx.filesDir, id)
	newName := path.Join(idx.indexRootDir, v, id)

	if err := os.MkdirAll(path.Join(idx.indexRootDir, v), 0777); err != nil {
		return "", err
	}

	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &alreadyExistsErr{idx.typeName, idx.indexBy, v}
	}

	return newName, err

}

func (idx NormalIndex) Remove(id string, v string) error {
	res, err := filepath.Glob(path.Join(idx.indexRootDir, "/*/", id))
	if err != nil {
		return err
	}

	for _, p := range res {
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	return nil
}

func (idx NormalIndex) Update(id, oldV, newV string) (err error) {
	oldDir := path.Join(idx.indexRootDir, oldV)
	oldPath := path.Join(oldDir, id)
	newDir := path.Join(idx.indexRootDir, newV)
	newPath := path.Join(newDir, id)

	if _, err = os.Stat(oldPath); os.IsNotExist(err) {
		return &notFoundErr{idx.typeName, idx.indexBy, oldV}
	}

	if err != nil {
		return
	}

	if err = os.MkdirAll(newDir, 0777); err != nil {
		return
	}

	if err = os.Rename(oldPath, newPath); err != nil {
		return
	}

	di, err := ioutil.ReadDir(oldDir)
	if err != nil {
		return err
	}

	if len(di) == 0 {
		err = os.RemoveAll(oldDir)
		if err != nil {
			return
		}
	}

	return

}

func (idx NormalIndex) Search(pattern string) ([]string, error) {
	paths, err := filepath.Glob(path.Join(idx.indexRootDir, pattern, "*"))
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, &notFoundErr{idx.typeName, idx.indexBy, pattern}
	}

	return paths, nil
}

func (idx NormalIndex) IndexBy() string {
	return idx.indexBy
}

func (idx NormalIndex) TypeName() string {
	return idx.typeName
}

func (idx NormalIndex) FilesDir() string {
	return idx.filesDir
}
