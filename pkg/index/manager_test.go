package index

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"io/ioutil"
	"testing"
)

func TestIntegration(t *testing.T) {
	type User struct {
		Id, Name, Email string
	}

	mcfg := &ManagerConfig{
		DataDir:          "/tmp/data4",
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
		open:             false,
	}

	idxer, err := Start(mcfg, Unique{
		indexBy:  "Email",
		typeName: "User",
		filesDir: "/tmp/data4/users",
		indexDir: "/tmp/data4/index.disk/",
	},
		Unique{
			indexBy:  "Name",
			typeName: "User",
			filesDir: "/tmp/data4/users",
			indexDir: "/tmp/data4/index.disk/",
		})

	if err != nil {
		t.Fatal(err)
	}

	u := &User{
		Id:    "430599a9-09d5-4524-b9f4-2295e612a151",
		Name:  "Mike",
		Email: "mike@example.com",
	}
	data, _ := json.Marshal(u)

	err = ioutil.WriteFile("/tmp/data4/users/430599a9-09d5-4524-b9f4-2295e612a151", data, 0777)

	if err != nil {
		t.Fatal(err)
	}

	err = idxer.Add(u.Id, u)
	if err != nil {
		t.Fatal(err)
	}

	res, err := idxer.Find("User", "Email", "mike@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if res != "430599a9-09d5-4524-b9f4-2295e612a151" {
		t.Fatal("Expected primary-key was not returned by Find")
	}
}
