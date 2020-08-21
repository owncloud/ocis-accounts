package storage

type TypeName string
type KeyAttribute string

type indexMap map[TypeName]map[KeyAttribute][]Index

type Indexer struct {
	dataFolder string
	indices    indexMap
}

func NewIndexer(dataFolder string, indices ...Index) *Indexer {
	indexer := &Indexer{
		dataFolder: dataFolder,
		indices:    make(indexMap),
	}

	for k := range indices {
		indexer.AddIndex(indices[k])
	}

	return indexer
}

func (ind Indexer) AddIndex(idx Index) {
	typeName := idx.TypeName()
	keyAttr := idx.KeyAttribute()

	if _, ok := ind.indices[typeName]; !ok {
		ind.indices[idx.TypeName()] = make(map[KeyAttribute][]Index)
	}

	ind.indices[typeName][keyAttr] = append(ind.indices[typeName][keyAttr], idx)

}

type Index interface {
	EntryAdded(v string, pk string)
	EntryRemoved(v string, pk string)
	KeyAttribute() KeyAttribute
	TypeName() TypeName
	RootFolder() string
}

type UniqueIndex struct {
	key        KeyAttribute
	typeName   TypeName
	rootFolder string
}

func (uniq UniqueIndex) KeyAttribute() KeyAttribute {
	return uniq.key
}

func (uniq UniqueIndex) TypeName() TypeName {
	return uniq.typeName
}

func (uniq UniqueIndex) RootFolder() string {
	return uniq.rootFolder
}

func (uniq UniqueIndex) EntryAdded(v string, pk string) {

}

func (uniq UniqueIndex) EntryRemoved(v string, pk string) {

}
