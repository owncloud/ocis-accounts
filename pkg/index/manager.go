// Package index provides symlink-based index for on-disk document-directories.
package index

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path"
	"reflect"
)

// Manager is a facade to configure and query over multiple indices.
type Manager struct {
	config    *ManagerConfig
	indices   indexMap
	typeToDir map[string]diskMap
}

type ManagerConfig struct {
	DataDir          string
	IndexRootDirName string
	Log              zerolog.Logger
}

// Index can be implemented to create new index-strategies. See Unique for example.
// Each index implementation is bound to one data-column (IndexBy) and a data-type (TypeName)
type Index interface {
	Init() error
	Lookup(v string) (string, error)
	Add(pk, v string) error
	Remove(v string) error
	Update(oldV, newV string) error
	IndexBy() string
	TypeName() string
	FilesDir() string
	BacklinksDir() string
}

func NewManager(cfg *ManagerConfig) *Manager {
	return &Manager{
		config:    cfg,
		indices:   indexMap{},
		typeToDir: map[string]diskMap{},
	}
}

func (man Manager) AddUniqueIndex(typeName, indexBy, entityDirName string) error {
	fullDataPath := path.Join(man.config.DataDir, entityDirName)
	indexPath := path.Join(man.config.DataDir, man.config.IndexRootDirName)

	idx := NewUniqueIndex(typeName, indexBy, fullDataPath, indexPath)
	man.indices.addIndex(idx)
	man.typeToDir[idx.typeName] = diskMap{
		filesDirPath:     fullDataPath,
		backlinksDirPath: idx.backlinkDir,
	}

	return idx.Init()
}

func (man Manager) AddIndex(idx Index) error {
	man.indices.addIndex(idx)
	man.typeToDir[idx.TypeName()] = diskMap{idx.FilesDir(), idx.BacklinksDir()}

	return idx.Init()
}

// Add a new entry to the index
func (man Manager) Add(primaryKey string, entity interface{}) error {
	t, err := getType(entity)
	if err != nil {
		return err
	}

	if typeIndices, ok := man.indices[t.Type().Name()]; ok {
		for fieldName, fieldIndices := range typeIndices {
			for k := range fieldIndices {
				curIdx := fieldIndices[k]
				idxBy := curIdx.IndexBy()
				f := t.FieldByName(idxBy)
				if f.IsZero() {
					return fmt.Errorf("the indexBy-name of an index on %v does not exist", fieldName)
				}

				err := curIdx.Add(primaryKey, f.String())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Find a entry by type,field and value.
//  // Find a User type by email
//  man.Find("User", "Email", "foo@example.com")
func (man Manager) Find(typeName, key, value string) (pk string, err error) {
	if indices, ok := man.indices[typeName][key]; ok {
		for _, idx := range indices {
			if pk, err = idx.Lookup(value); IsNotFoundErr(err) {
				continue
			}

			if err != nil {
				return
			}
		}
	}

	if pk == "" {
		return
	}

	return path.Base(pk), err
}

func (man Manager) Delete(typeName, pk string) error {
	if dm, ok := man.typeToDir[typeName]; ok {
		entityBacklinksDir := path.Join(dm.backlinksDirPath, pk)
		fi, err := os.Stat(entityBacklinksDir)
		if os.IsNotExist(err) {
			return &notFoundErr{typeName, "_PRIMARY_", pk}
		}

		if !fi.IsDir() {
			return fmt.Errorf("%s is supposed to be a directory (corruption/bug?)", fi.Name())
		}

		blInfos, err := ioutil.ReadDir(entityBacklinksDir)
		if err != nil {
			return err
		}

		for _, blInfo := range blInfos {
			blPath := path.Join(entityBacklinksDir, blInfo.Name())
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
	}

	return nil
}

type diskMap struct {
	filesDirPath, backlinksDirPath string
}

// indexMap holds the index-configuration
type indexMap map[tName]map[indexByKey][]Index
type tName = string
type indexByKey = string

func (m indexMap) addIndex(idx Index) {
	typeName, indexBy := idx.TypeName(), idx.IndexBy()
	if _, ok := m[typeName]; !ok {
		m[typeName] = map[indexByKey][]Index{}
	}

	m[typeName][indexBy] = append(m[typeName][indexBy], idx)
}

func getType(v interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return reflect.Value{}, errors.New("error while detecting entity type for index-update")
	}

	return rv, nil
}

func getValueOf(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}
