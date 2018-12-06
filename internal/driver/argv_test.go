package driver

import (
	"errors"
	"testing"
)

func trueParser(s string) (string, error) {
	return s, nil
}

func falseParser(s string) (string, error) {
	return "", errors.New("Invalid")
}

func TestLoadHostsValidData(t *testing.T) {
	data := []string{"192.168.1.1", "8.8.8.8"}
	res, err := ArgvLoadHosts(trueParser, data)
	if err != nil {
		t.Error(`argvLoadHosts() returns error`)
	}

	for i := range data {
		if res[i].IP != data[i] {
			t.Errorf(`argvLoadHosts() for data[%d] returns incorrect IP`, i)
		}
		if res[i].ID != 0 {
			t.Errorf(`argvLoadHosts() for data[%d] returns not empty ID`, i)
		}
	}
}

func TestLoadHostsInvalidData(t *testing.T) {
	data := []string{"some/invalid/host"}
	_, err := ArgvLoadHosts(falseParser, data)
	if err == nil {
		t.Error(`argvLoadHosts() doesn't return error`)
	}
}
