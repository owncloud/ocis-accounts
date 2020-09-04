package index

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestIndex(t *testing.T) {
	sut, _ := getNormalIdxSut(t)

	ids, err := sut.Lookup("Green")
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"goefe-789", "xadaf-189"}, ids)

	ids, err = sut.Lookup("White")
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"wefwe-456"}, ids)

	ids, err = sut.Lookup("Cyan")
	assert.Error(t, err)
	assert.EqualValues(t, []string{}, ids)

}

func TestIndexInit(t *testing.T) {
	dataDir := createTmpDir(t)
	indexRootDir := path.Join(dataDir, "index.disk")
	filesDir := path.Join(dataDir, "users")

	uniq := NewNormalIndex("User", "DisplayName", filesDir, indexRootDir)
	assert.Error(t, uniq.Init(), "Init should return an error about missing files-dir")

	if err := os.Mkdir(filesDir, 0777); err != nil {
		t.Fatalf("Could not create test data-dir %s", err)
	}

	assert.NoError(t, uniq.Init(), "Init shouldn't return an error")
	assert.DirExists(t, indexRootDir)
	assert.DirExists(t, path.Join(indexRootDir, "UserByDisplayName"))

	_ = os.RemoveAll(dataDir)
}

func getNormalIdxSut(t *testing.T) (sut Index, dataPath string) {
	dataPath = writeIndexTestData(t, testData, "Id")
	sut = NewNormalIndex("Pet", "Color", path.Join(dataPath, "pets"), path.Join(dataPath, "index.disk"))
	err := sut.Init()
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range testData["pets"] {
		pkVal := valueOf(u, "Id")
		idxByVal := valueOf(u, "Color")
		_, err := sut.Add(pkVal, idxByVal)
		if err != nil {
			t.Fatal(err)
		}
	}

	return
}
