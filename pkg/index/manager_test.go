package index

/*
func TestManagerQueryMultipleIndices(t *testing.T) {
	dataDir := writeIndexTestData(t, testData, "Id")
	man := NewManager(&ManagerConfig{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	err := man.AddPrimaryIndex("User", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("User", "Email", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("User", "UserName", "users")
	assert.NoError(t, err)

	err = man.AddPrimaryIndex("TestPet", "pets")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("TestPet", "Color", "pets")
	assert.NoError(t, err)

	for path := range testData {
		for _, entity := range testData[path] {
			err := man.Add(valueOf(entity, "Id"), entity)
			assert.NoError(t, err)
		}
	}

	type test struct {
		typeName, key, value, wantRes string
		wantErr                       error
	}

	tests := []test{
		{typeName: "User", key: "Email", value: "jacky@example.com", wantRes: "ewf4ofk-555"},
		{typeName: "User", key: "UserName", value: "jacky", wantRes: "ewf4ofk-555"},
		{typeName: "TestPet", key: "Color", value: "Brown", wantRes: "rebef-123"},
		{typeName: "TestPet", key: "Color", value: "Cyan", wantRes: "", wantErr: &notFoundErr{}},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("Query%sBy%s=%s", tc.typeName, tc.key, tc.value)
		t.Run(name, func(t *testing.T) {
			pk, err := man.Find(tc.typeName, tc.key, tc.value)
			assert.Equal(t, tc.wantRes, pk)
			assert.IsType(t, tc.wantErr, err)
		})
	}

	_ = os.RemoveAll(dataDir)
}

func TestManagerDelete(t *testing.T) {
	dataDir := writeIndexTestData(t, testData, "Id")
	man := NewManager(&ManagerConfig{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	err := man.AddPrimaryIndex("User", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("User", "Email", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("User", "UserName", "users")
	assert.NoError(t, err)

	err = man.AddPrimaryIndex("TestPet", "pets")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("TestPet", "Color", "pets")
	assert.NoError(t, err)

	for path := range testData {
		for _, entity := range testData[path] {
			err := man.Add(valueOf(entity, "Id"), entity)
			assert.NoError(t, err)
		}
	}

	err = man.Delete("User", "hijklmn-456")
	_ = os.RemoveAll(dataDir)

}

*/
