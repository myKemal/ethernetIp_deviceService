package config

import (
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

func TestCreateConnectionInfo(t *testing.T) {

	type fields struct {
		Address     string
		CClimit     string
		Port        string
		TimeOut     string
		IdleTimeOut string
	}
	testsDatas := []struct {
		name      string
		fields    fields
		wantError bool
	}{
		{
			name:      "NOK - need protocol properties",
			fields:    fields{},
			wantError: true,
		},
		{
			name:      "NOK - invalid protocol address ",
			fields:    fields{Address: "A", CClimit: "50", Port: "0", TimeOut: "3000", IdleTimeOut: "1500"},
			wantError: true,
		},
		{
			name:      "NOK - invalid protocol datatype ",
			fields:    fields{Address: "A", CClimit: "50", Port: "0", TimeOut: "3000", IdleTimeOut: "1500"},
			wantError: true,
		},
		{
			name:      "NOK - invalid protocol configuration ",
			fields:    fields{CClimit: "50", Port: "0", TimeOut: "3000", IdleTimeOut: "1500"},
			wantError: true,
		},
		{
			name:      "OK - valid protocol configuretion",
			fields:    fields{Address: "0.0.0.0", CClimit: "50", Port: "0", TimeOut: "3000", IdleTimeOut: "1500"},
			wantError: false,
		},
	}

	for _, td := range testsDatas {
		t.Run(td.name, func(t *testing.T) {

			protocols := map[string]models.ProtocolProperties{
				"device-ethernetip-go": {
					"Address":                td.fields.Address,
					"ConcurrentCommandLimit": td.fields.CClimit,
					"Port":                   td.fields.Port,
					"Timeout":                td.fields.TimeOut,
					"IdleTimeout":            td.fields.IdleTimeOut,
				},
			}

			connectionInfo, err := CreateConnectionInfo(protocols)

			if err != nil && !td.wantError || err == nil && td.wantError {
				t.Errorf("Created ConnectionInfo --> %v, Error %v,Func Error %v", connectionInfo, td.wantError, err)
			}

		})
	}

}
