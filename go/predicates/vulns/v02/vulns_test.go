/*
Tests for vulns v0.2 protos.
*/

package v02

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

const validVulns = `{
  "scanner": {
    "uri": "pkg:github/aquasecurity/trivy@244fd47e07d1004f0aed9",
    "version": "0.19.2",
    "db": {
      "uri": "pkg:github/aquasecurity/trivy-db/commit/4c76bb580b2736d67751410fa4ab66d2b6b9b27d",
      "version": "v1-2021080612",
      "lastUpdate": "2021-08-06T17:45:50.52Z"
    },
    "result": [
      {
        "id": "CVE-123",
        "severity": [
          {"method": "nvd", "score": "medium"},
          {"method": "cvss_score", "score": "5.2"}
        ]
      }
    ]
  },
  "metadata": {
    "scanStartedOn": "2021-08-06T17:45:50.52Z",
    "scanFinishedOn": "2021-08-06T17:50:50.52Z"
  }
}`

func parseVulns(t *testing.T, s string) *Vulns {
	t.Helper()
	v := &Vulns{}
	if err := protojson.Unmarshal([]byte(s), v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestValidVulns(t *testing.T) {
	v := parseVulns(t, validVulns)
	assert.NoError(t, v.Validate())
}

func TestValidVulnsNoResults(t *testing.T) {
	v := parseVulns(t, validVulns)
	v.Scanner.Result = nil
	assert.NoError(t, v.Validate(), "an empty result list means no vulnerabilities were found")
}

func TestBadVulns(t *testing.T) {
	tests := map[string]struct {
		mutate func(v *Vulns)
		err    error
		msg    string
	}{
		"no scanner": {
			mutate: func(v *Vulns) { v.Scanner = nil },
			err:    ErrScannerRequired,
		},
		"no scanner uri": {
			mutate: func(v *Vulns) { v.Scanner.Uri = "" },
			err:    ErrScannerUriRequired,
		},
		"no db": {
			mutate: func(v *Vulns) { v.Scanner.Db = nil },
			err:    ErrDbRequired,
		},
		"no db lastUpdate": {
			mutate: func(v *Vulns) { v.Scanner.Db.LastUpdate = nil },
			err:    ErrDbLastUpdateRequired,
		},
		"no result id": {
			mutate: func(v *Vulns) { v.Scanner.Result[0].Id = "" },
			err:    ErrResultIdRequired,
			msg:    "scanner.result[0].id required",
		},
		"no severity": {
			mutate: func(v *Vulns) { v.Scanner.Result[0].Severity = nil },
			err:    ErrSeverityRequired,
		},
		"no severity method": {
			mutate: func(v *Vulns) { v.Scanner.Result[0].Severity[1].Method = "" },
			err:    ErrSeverityMethodRequired,
			msg:    "scanner.result[0].severity[1].method required",
		},
		"no severity score": {
			mutate: func(v *Vulns) { v.Scanner.Result[0].Severity[0].Score = "" },
			err:    ErrSeverityScoreRequired,
		},
		"no metadata": {
			mutate: func(v *Vulns) { v.Metadata = nil },
			err:    ErrMetadataRequired,
		},
		"no scanStartedOn": {
			mutate: func(v *Vulns) { v.Metadata.ScanStartedOn = nil },
			err:    ErrScanStartedOnRequired,
		},
		"no scanFinishedOn": {
			mutate: func(v *Vulns) { v.Metadata.ScanFinishedOn = nil },
			err:    ErrScanFinishedOnRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			v := parseVulns(t, validVulns)
			tt.mutate(v)
			err := v.Validate()
			assert.ErrorIs(t, err, tt.err)
			if tt.msg != "" {
				assert.EqualError(t, err, tt.msg)
			}
		})
	}
}

func TestBadVulnsTimestamp(t *testing.T) {
	v := parseVulns(t, validVulns)
	v.Metadata.ScanStartedOn.Nanos = -1
	assert.ErrorContains(t, v.Validate(), "metadata.scanStartedOn error")
}
