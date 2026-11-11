// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package gorm

import (
	"testing"

	"github.com/agntcy/dir/server/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Reproduction for agntcy/dir#1874: two scans of the same record by the same
// scanner should leave two scans to inspect. Today the second overwrites the
// first, so this test is expected to FAIL on main.
func TestRepro1874_ScanHistoryIsKept(t *testing.T) {
	db := setupScanDB(t)

	upsert(t, db, &ScanReport{
		RecordCID: scanTestCID, ScannerType: "MCP",
		Status: types.ScanStatusFailed, FailureReason: "network", MaxSeverity: "NONE",
	})
	upsert(t, db, &ScanReport{
		RecordCID: scanTestCID, ScannerType: "MCP",
		Status: types.ScanStatusCompleted, IsSafe: true, MaxSeverity: "NONE",
	})

	var rows []ScanReport
	require.NoError(t, db.gormDB.Where("record_cid = ? AND scanner_type = ?", scanTestCID, "MCP").Find(&rows).Error)

	for _, r := range rows {
		t.Logf("scan_reports row: status=%s failure_reason=%q consecutive_failures=%d", r.Status, r.FailureReason, r.ConsecutiveFailures)
	}

	assert.Len(t, rows, 2, "expected one row per scan (failed, then completed); scan_reports keeps only the latest")
}
