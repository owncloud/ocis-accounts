package index

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// primary is created for each document-type to track existing primary-key id's and the indices referencing them.
//
// On disk a folder for each entity is created which contains symlinks to other index-symlinks pointing at this entity.
// This allows for performant update and delete operations.
//
// 	/var/data/index.disk/UserPrimary/ewf4ofk-555/
//	├── jacky -> /var/data/index.disk/UserByUserName/jacky
//	└── jacky@example.com -> /var/data/index.disk/UserByEmail/jacky@example.com
//
type primary struct {
	typeName string
	indexDir string
	dataPath string
}

func (idx primary) init() error {
	if err := os.MkdirAll(idx.indexDir, 0777); err != nil {
		return err
	}

	return nil
}

func (idx primary) add(pk, v, createdPath string) error {
	priIdxPath := path.Join(idx.indexDir, pk)
	if err := os.MkdirAll(priIdxPath, 0777); err != nil {
		return err
	}

	return os.Symlink(createdPath, path.Join(priIdxPath, v))
}

func (idx primary) delete(pk string) error {
	entityIdxPath := path.Join(idx.indexDir, pk)
	fi, err := os.Stat(entityIdxPath)
	if os.IsNotExist(err) {
		return &notFoundErr{idx.typeName, "_PRIMARY_", pk}
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is supposed to be a directory (corruption/bug?)", fi.Name())
	}

	linksFi, err := ioutil.ReadDir(entityIdxPath)
	if err != nil {
		return err
	}

	for _, blInfo := range linksFi {
		blPath := path.Join(entityIdxPath, blInfo.Name())
		if err := isValidSymlink(blPath); err != nil {
			return err
		}

		origPath, err := os.Readlink(blPath)
		if err != nil {
			return err
		}

		if err := os.Remove(blPath); err != nil {
			return err
		}

		if err := os.Remove(origPath); err != nil {
			return err
		}
	}

	return os.RemoveAll(entityIdxPath)
}
