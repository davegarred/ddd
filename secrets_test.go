package ddd

import (
	"fmt"
	"testing"
)

func TestGetSecret(t *testing.T) {
	secrets,err := GetSecret("stc-admin-promary")
	if err != nil {
		panic(err)
	}
	fmt.Println(secrets)
	line := fmt.Sprintf("postgres://%s:%s@%s/stc?sslmode=disable", secrets["username"], secrets["password"], secrets["host"])
	if line != "postgres://postgres:kCEWyKGnXMqjgYFyGNTA@stc-admin.cluster-csiurdhhg3zk.us-west-2.rds.amazonaws.com/stc?sslmode=disable" {
		panic(line)
	}
}
