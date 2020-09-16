package index

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestUniqueLookupSingleEntry(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t)
	filesDir := path.Join(dataDir, "users")

	t.Log("existing lookup")
	resultPath, err := uniq.Lookup("mikey@example.com")
	assert.NoError(t, err)

	assert.Equal(t, []string{path.Join(filesDir, "abcdefg-123")}, resultPath)

	t.Log("non-existing lookup")
	resultPath, err = uniq.Lookup("doesnotExists@example.com")
	assert.Error(t, err)
	assert.IsType(t, &notFoundErr{}, err)
	assert.Empty(t, resultPath)

	_ = os.RemoveAll(dataDir)

}

func TestUniqueUniqueConstraint(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t)

	_, err := uniq.Add("abcdefg-123", "mikey@example.com")
	assert.Error(t, err)
	assert.IsType(t, &alreadyExistsErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestUniqueRemove(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t)

	err := uniq.Remove("", "mikey@example.com")
	assert.NoError(t, err)

	_, err = uniq.Lookup("mikey@example.com")
	assert.Error(t, err)
	assert.IsType(t, &notFoundErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestUniqueUpdate(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t)

	t.Log("successful update")
	err := uniq.Update("", "mikey@example.com", "mikey2@example.com")
	assert.NoError(t, err)

	t.Log("failed update because already exists")
	err = uniq.Update("", "frank@example.com", "mikey2@example.com")
	assert.Error(t, err)
	assert.IsType(t, &alreadyExistsErr{}, err)

	t.Log("failed update because not found")
	err = uniq.Update("", "notexist@example.com", "something2@example.com")
	assert.Error(t, err)
	assert.IsType(t, &notFoundErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestUniqueInit(t *testing.T) {
	dataDir := createTmpDir(t)
	indexRootDir := path.Join(dataDir, "index.disk")
	filesDir := path.Join(dataDir, "users")

	uniq := NewUniqueIndex("User", "Email", filesDir, indexRootDir)
	assert.Error(t, uniq.Init(), "Init should return an error about missing files-dir")

	if err := os.Mkdir(filesDir, 0777); err != nil {
		t.Fatalf("Could not create test data-dir %s", err)
	}

	assert.NoError(t, uniq.Init(), "Init shouldn't return an error")
	assert.DirExists(t, indexRootDir)
	assert.DirExists(t, path.Join(indexRootDir, "UniqueUserByEmail"))

	_ = os.RemoveAll(dataDir)
}

func TestUniqueIndexSearch(t *testing.T) {
	sut, dataPath := getUniqueIdxSut(t)

	res, err := sut.Search("j*@example.com")

	assert.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t, "ewf4ofk-555", path.Base(res[0]))
	assert.Equal(t, "rulan54-777", path.Base(res[1]))

	res, err = sut.Search("does-not-exist@example.com")
	assert.Error(t, err)
	assert.IsType(t, &notFoundErr{}, err)

	_ = os.RemoveAll(dataPath)
}

func TestErrors(t *testing.T) {
	assert.True(t, IsAlreadyExistsErr(&alreadyExistsErr{}))
	assert.True(t, IsNotFoundErr(&notFoundErr{}))
}

func getUniqueIdxSut(t *testing.T) (sut Index, dataPath string) {
	dataPath = writeIndexTestData(t, testData, "Id")
	sut = NewUniqueIndex("User", "Email", path.Join(dataPath, "users"), path.Join(dataPath, "index.disk"))
	err := sut.Init()
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range testData["users"] {
		pkVal := valueOf(u, "Id")
		idxByVal := valueOf(u, "Email")
		_, err := sut.Add(pkVal, idxByVal)
		if err != nil {
			t.Fatal(err)
		}
	}

	return
}
