package database

import "testing"

func TestNew(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()
	_ = New()
}
