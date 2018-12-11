package schema

import (
	"errors"
	"testing"
)

func TestParseHost(t *testing.T) {
	cases := []struct {
		Input          string
		Response       string
		ExpectedErrStr string
	}{
		{
			"192.168.1.1",
			"192.168.1.1",
			"",
		},
		{
			"192.168.1.1/24",
			"192.168.1.1",
			"",
		},
		{
			"wp.pl",
			"wp.pl",
			"",
		},
		{
			"wp.pl/24",
			"",
			"Can't resolve host: wp.pl/24",
		},
		{
			"192.168.1.1.1.1",
			"",
			"Can't resolve host: 192.168.1.1.1.1",
		},
		{
			"192.168.1.1/abc",
			"",
			"Can't resolve host: 192.168.1.1/abc",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			var hosts Hosts
			_, _, err := hosts.parseHost(tc.Input)

			if err == nil && tc.ExpectedErrStr != "" ||
				err != nil && tc.ExpectedErrStr != err.Error() {
				t.Errorf("got: %v expected: %v", err, tc.ExpectedErrStr)
			}
		})
	}
}

func TestValidHostsSetGet(t *testing.T) {
	var hosts Hosts
	validHosts := []Host{{IP: "192.168.1.1", ID: 0}, {IP: "10.10.0.1", ID: 0}}

	loader := func(parser HostParser) ([]Host, error) {
		return validHosts, nil
	}

	if err := hosts.Add(loader); err != nil {
		t.Error("hosts.Set returns error on valid loader")
	}

	for i, host := range hosts.Get() {
		if host != validHosts[i] {
			t.Errorf("hosts.Get returns invalid host id %d", i)
		}
	}
}

func TestInvalidHostsSetGet(t *testing.T) {
	var hosts Hosts

	errorLoader := func(parser HostParser) ([]Host, error) {
		return nil, errors.New("error")
	}

	if err := hosts.Add(errorLoader); err == nil {
		t.Error("hosts.Set doesn't return error")
	}

	if hosts := hosts.Get(); hosts != nil {
		t.Error("hosts.Get returns hosts")
	}
}
