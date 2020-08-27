package index

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path"
	"reflect"
)

type Manager struct {
	config  *ManagerConfig
	indices indexMap
}

type ManagerConfig struct {
	DataDir          string
	IndexRootDirName string
	Log              zerolog.Logger
	open             bool
}

type Index interface {
	Lookup(v string) (string, error)
	Add(pk, v string) error
	EntryRemoved(pk, v string)
	IndexBy() string
	TypeName() string
	Init() error
}

func Start(cfg *ManagerConfig, indices ...Index) (*Manager, error) {
	err := os.MkdirAll(path.Join(cfg.DataDir, cfg.IndexRootDirName), 0777)
	if err != nil {
		return nil, err
	}
	m := &Manager{config: cfg, indices: indexMap{}}
	for _, idx := range indices {
		m.indices.addIndex(idx)
		if err := idx.Init(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

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

func (man Manager) Find(typeName, field, value string) (string, error) {
	if indices, ok := man.indices[typeName][field]; ok {
		for _, idx := range indices {
			pk, err := idx.Lookup(value)
			if err != nil {
				return "", err
			}

			return pk, nil
		}
	}

	return "", errors.New("not found")

}

// indexMap holds the index-configuration
type indexMap map[typeName]map[indexByKey][]Index
type typeName = string
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
