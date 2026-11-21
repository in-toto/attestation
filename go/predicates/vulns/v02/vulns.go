/*
Validator APIs for vulns v0.2 protos.
*/
package v02

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrScannerRequired        = errors.New("scanner required")
	ErrScannerUriRequired     = errors.New("scanner.uri required")
	ErrDbRequired             = errors.New("scanner.db required")
	ErrDbLastUpdateRequired   = errors.New("scanner.db.lastUpdate required")
	ErrResultIdRequired       = errors.New("id required")
	ErrSeverityRequired       = errors.New("severity required")
	ErrSeverityMethodRequired = errors.New("method required")
	ErrSeverityScoreRequired  = errors.New("score required")
	ErrMetadataRequired       = errors.New("metadata required")
	ErrScanStartedOnRequired  = errors.New("metadata.scanStartedOn required")
	ErrScanFinishedOnRequired = errors.New("metadata.scanFinishedOn required")
)

func checkTimestamp(ts *timestamppb.Timestamp, name string, required error) error {
	if ts == nil {
		return required
	}
	if err := ts.CheckValid(); err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}
	return nil
}

func (s *Result_Severity) Validate() error {
	if s.GetMethod() == "" {
		return ErrSeverityMethodRequired
	}
	if s.GetScore() == "" {
		return ErrSeverityScoreRequired
	}
	return nil
}

func (r *Result) Validate() error {
	if r.GetId() == "" {
		return ErrResultIdRequired
	}

	severity := r.GetSeverity()
	if len(severity) == 0 {
		return ErrSeverityRequired
	}
	for i, s := range severity {
		if s == nil {
			return fmt.Errorf("severity[%d].%w", i, ErrSeverityMethodRequired)
		}
		if err := s.Validate(); err != nil {
			return fmt.Errorf("severity[%d].%w", i, err)
		}
	}

	return nil
}

func (d *VulnDatabase) Validate() error {
	return checkTimestamp(d.GetLastUpdate(), "scanner.db.lastUpdate", ErrDbLastUpdateRequired)
}

func (s *Scanner) Validate() error {
	if s.GetUri() == "" {
		return ErrScannerUriRequired
	}

	db := s.GetDb()
	if db == nil {
		return ErrDbRequired
	}
	if err := db.Validate(); err != nil {
		return err
	}

	// An empty result list is valid and means no vulnerabilities were found.
	for i, r := range s.GetResult() {
		if r == nil {
			return fmt.Errorf("scanner.result[%d].%w", i, ErrResultIdRequired)
		}
		if err := r.Validate(); err != nil {
			return fmt.Errorf("scanner.result[%d].%w", i, err)
		}
	}

	return nil
}

func (m *ScanMetadata) Validate() error {
	if err := checkTimestamp(m.GetScanStartedOn(), "metadata.scanStartedOn", ErrScanStartedOnRequired); err != nil {
		return err
	}
	return checkTimestamp(m.GetScanFinishedOn(), "metadata.scanFinishedOn", ErrScanFinishedOnRequired)
}

func (v *Vulns) Validate() error {
	scanner := v.GetScanner()
	if scanner == nil {
		return ErrScannerRequired
	}
	if err := scanner.Validate(); err != nil {
		return err
	}

	metadata := v.GetMetadata()
	if metadata == nil {
		return ErrMetadataRequired
	}
	return metadata.Validate()
}
