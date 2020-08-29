// Package index provides symlink-based index for on-disk document-directories.
package index

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"path"
	"reflect"
)

// Manager is a facade to configure and query over multiple indices.
type Manager struct {
	config         *ManagerConfig
	indices        indexMap
	typeToDir      map[string]diskMap
	primaryIndices map[string]primaryIndex
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
	Add(pk, v string) (string, error)
	Remove(v string) error
	Update(oldV, newV string) error
	IndexBy() string
	TypeName() string
	FilesDir() string
}

func NewManager(cfg *ManagerConfig) *Manager {
	return &Manager{
		config:         cfg,
		indices:        indexMap{},
		typeToDir:      map[string]diskMap{},
		primaryIndices: map[string]primaryIndex{},
	}
}

func (man Manager) AddPrimaryIndex(typeName, entityDirName string) error {
	fullDataPath := path.Join(man.config.DataDir, entityDirName)
	indexPath := path.Join(man.config.DataDir, man.config.IndexRootDirName)
	man.primaryIndices[typeName] = primaryIndex{
		typeName: typeName,
		indexDir: path.Join(indexPath, fmt.Sprintf("%sPrimary", typeName)),
		dataPath: fullDataPath,
	}

	return man.primaryIndices[typeName].init()

}

func (man Manager) AddUniqueIndex(typeName, indexBy, entityDirName string) error {
	fullDataPath := path.Join(man.config.DataDir, entityDirName)
	indexPath := path.Join(man.config.DataDir, man.config.IndexRootDirName)

	idx := NewUniqueIndex(typeName, indexBy, fullDataPath, indexPath)
	man.indices.addIndex(idx)

	return idx.Init()
}

func (man Manager) AddIndex(idx Index) error {
	man.indices.addIndex(idx)

	return idx.Init()
}

// Add a new entry to the index
func (man Manager) Add(primaryKey string, entity interface{}) error {
	t, err := getType(entity)
	if err != nil {
		return err
	}

	typeName := t.Type().Name()

	if typeIndices, ok := man.indices[typeName]; ok {
		for fieldName, fieldIndices := range typeIndices {
			for k := range fieldIndices {
				curIdx := fieldIndices[k]
				idxBy := curIdx.IndexBy()
				f := t.FieldByName(idxBy)
				if f.IsZero() {
					return fmt.Errorf("the indexBy-name of an index on %v does not exist", fieldName)
				}

				val := f.String()
				created, err := curIdx.Add(primaryKey, val)
				if err != nil {
					return err
				}

				if priIdx, ok := man.primaryIndices[typeName]; ok {
					err := priIdx.add(primaryKey, val, created)
					if err != nil {
						return err
					}
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
	if priIdx, ok := man.primaryIndices[typeName]; ok {
		err := priIdx.delete(pk)
		if err != nil {
			return err
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
