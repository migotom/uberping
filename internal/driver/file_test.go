package driver

import (
	"testing"

	"github.com/migotom/uberping/internal/schema"
)

func TestInvalidFileLoadHosts(t *testing.T) {
	hosts, err := FileLoadHosts(trueParser, "inavalid_file_name")
	if err == nil {
		t.Error("fileLoadHosts doesn't return error on non existing file")
	}

	if hosts != nil {
		t.Error("fileLoadHosts returns hosts on reading invalid file")
	}
}

func TestValidFileLoadHosts(t *testing.T) {
	hosts, err := FileLoadHosts(trueParser, "../../hosts-file.example.txt")
	if err != nil {
		t.Error("fileLoadHosts returns error on reading existing file")
	}
	if hosts == nil {
		t.Error("fileLoadHosts doesn't return hosts on reading valid file")
	}

	var testHosts = []schema.Host{
		schema.Host{IP: "google.com", ID: 0},
		schema.Host{IP: "192.168.1.1", ID: 0},
		schema.Host{IP: "192.168.88.1/24", ID: 0}}
	for i := range testHosts {
		if hosts == nil || testHosts[i] != hosts[i] {
			t.Errorf("fileLoadHosts doesn't return valid host on id %d", i)
		}
	}
}

func TestValidFileFalseParserLoadHosts(t *testing.T) {
	hosts, err := FileLoadHosts(falseParser, "hosts-file.example.txt")
	if err == nil {
		t.Error("fileLoadHosts doesn't return error while parsing hosts with falseParser")
	}
	if hosts != nil {
		t.Error("fileLoadHosts returns hosts while parsing with falseParser")
	}
}
