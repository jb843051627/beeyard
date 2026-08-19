# rebuild_branches.ps1 - Recreate all 30 bug branches from amended main
$ErrorActionPreference = "Stop"
Set-Location "D:\Develop\Workspace-trae\go-workspace\beeyard"
$env:GOTOOLCHAIN = "local"

# Step 1: Amend main with comment fixes + neutral message
git add -A
git commit --amend -m "improve test data and code comments" 2>&1

# Step 2: Delete all old bug branches
for ($i=1; $i -le 10; $i++) {
    git branch -D "bug_base_$i" 2>$null
    git branch -D "gold_model_fix_$i" 2>$null
    git branch -D "test_model_fix_$i" 2>$null
}

# Helper: apply bug edit to a file
function Apply-BugEdit($file, $find, $replace) {
    $content = Get-Content $file -Raw
    $content = $content -replace [regex]::Escape($find), $replace
    Set-Content $file -Value $content -NoNewline
}

# Helper: create bug_base, commit, gold (restore from main), test
function New-BugBranches($n, $files, $edits) {
    git checkout -b "bug_base_$n" main 2>&1
    foreach ($e in $edits) {
        Apply-BugEdit $e.File $e.Find $e.Replace
    }
    # Verify build
    go build ./... 2>&1 | Out-Null
    git add -A
    git commit -m "update $($files -join ',')" 2>&1
    # Gold
    git checkout -b "gold_model_fix_$n" 2>&1
    foreach ($f in $files) { git checkout main -- $f 2>&1 }
    git add -A
    git commit -m "restore $($files -join ',')" 2>&1
    # Test
    git checkout -b "test_model_fix_$n" "bug_base_$n" 2>&1
}

Write-Host "=== Creating 30 branches ==="

# Bug-001: hive_store (nil,nil) + hive_service (ignore err)
New-BugBranches 1 @("internal/store/hive_store.go","internal/service/hive_service.go") @(
    @{File="internal/store/hive_store.go"; Find="return nil, model.ErrHiveNotFound"; Replace="return nil, nil"},
    @{File="internal/service/hive_service.go"; Find=`t"h, err := s.store.GetByID(ctx, id)`n`tif err != nil {`n`t`treturn nil, fmt.Errorf(`"hive service get %d: %w`", id, err)`n`t}`n`treturn h, nil"; Replace=`t"h, _ := s.store.GetByID(ctx, id)`n`treturn h, nil"},
    @{File="internal/service/hive_service.go"; Find=`t"context"`n`t"fmt`n`t"time""; Replace=`t"context"`n`t"time""}
)

# Bug-002: alert_store (no copy) + alert_service (sort in place)
New-BugBranches 2 @("internal/store/alert_store.go","internal/service/alert_service.go") @(
    @{File="internal/store/alert_store.go"; Find=`t"out := make([]*model.Alert, len(cached))`n`tcopy(out, cached)`n`treturn out, nil"; Replace=`t"return cached, nil"},
    @{File="internal/store/alert_store.go"; Find=`t"out := make([]*model.Alert, len(list))`n`tcopy(out, list)`n`treturn out, nil"; Replace=`t"return list, nil"},
    @{File="internal/service/alert_service.go"; Find=`t"sorted := make([]*model.Alert, len(alerts))`n`tcopy(sorted, alerts)`n`tsort.Slice(sorted, func(i, j int) bool {`n`t`treturn sorted[i].CreatedAt.After(sorted[j].CreatedAt)`n`t})`n`treturn sorted, nil"; Replace=`t"sort.Slice(alerts, func(i, j int) bool {`n`t`treturn alerts[i].CreatedAt.After(alerts[j].CreatedAt)`n`t})`n`treturn alerts, nil"}
)

# Bug-003: reading_cache (RLock + no Get lock)
New-BugBranches 3 @("internal/cache/reading_cache.go") @(
    @{File="internal/cache/reading_cache.go"; Find="c.mu.Lock()`n`tdefer c.mu.Unlock()"; Replace="c.mu.RLock()`n`tdefer c.mu.RUnlock()"},
    @{File="internal/cache/reading_cache.go"; Find="c.mu.RLock()`n`tdefer c.mu.RUnlock()`n`tif m, ok := c.latest[hiveID]; ok {"; Replace="if m, ok := c.latest[hiveID]; ok {"}
)

# Bug-004: alert_service (remove validation) + alert_handler (500)
New-BugBranches 4 @("internal/service/alert_service.go","internal/handler/alert_handler.go") @(
    @{File="internal/service/alert_service.go"; Find=`t"switch a.Level {`n`tcase model.AlertLevelInfo, model.AlertLevelWarning, model.AlertLevelCritical:`n`tdefault:`n`t`treturn 0, model.NewValidationError(`"level`", `"invalid alert level: `"+a.Level)`n`t}`n`t"; Replace=`t"},
    @{File="internal/handler/alert_handler.go"; Find="writeError(w, err)"; Replace="writeJSON(w, http.StatusInternalServerError, map[string]string{`"error`": err.Error()})"}
)

# Bug-005: reading_service (context.Background) + reading_handler (no WithTimeout)
New-BugBranches 5 @("internal/service/reading_service.go","internal/handler/reading_handler.go") @(
    @{File="internal/service/reading_service.go"; Find=`t"if err := ctx.Err(); err != nil {`n`t`treturn total, fmt.Errorf(`"batch ingest cancelled: %w`", err)`n`t}`n`tid, err := s.store.Record(ctx,"; Replace=`t"id, err := s.store.Record(context.Background(),"},
    @{File="internal/handler/reading_handler.go"; Find=`t"context`n`t`""; Replace=`t`""},
    @{File="internal/handler/reading_handler.go"; Find=`t"ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)`n`tdefer cancel()`n`ttotal, err := h.svc.BatchIngest(ctx, batch)"; Replace=`t"total, err := h.svc.BatchIngest(r.Context(), batch)"}
)

# Bug-006: queen_store (nil,nil) + transfer_service (ignore err)
New-BugBranches 6 @("internal/store/queen_store.go","internal/service/transfer_service.go") @(
    @{File="internal/store/queen_store.go"; Find="return nil, model.ErrQueenNotFound"; Replace="return nil, nil"},
    @{File="internal/service/transfer_service.go"; Find=`t"q, err := s.queens.GetByID(ctx, queenID)`n`tif err != nil {`n`t`treturn fmt.Errorf(`"get queen %d for assign: %w`", queenID, err)`n`t}`n`t_ = q`n`treturn s.store.UpdateStatus"; Replace=`t"_, _ = s.queens.GetByID(ctx, queenID)`n`treturn s.store.UpdateStatus"},
    @{File="internal/service/transfer_service.go"; Find=`t"context`n`t"fmt`n`t"time`""; Replace=`t"context`n`t"time`""}
)

# Bug-007: maintenance_service (own tx with bad defer)
New-BugBranches 7 @("internal/service/maintenance_service.go") @(
    @{File="internal/service/maintenance_service.go"; Find="return s.store.BatchCreate(ctx, tasks)"; Replace="tx, err := s.db.BeginTx(ctx, nil)`n`tif err != nil {`n`t`treturn 0, err`n`t}`n`tdefer tx.Rollback()`n`tvar total int64`n`tfor _, t := range tasks {`n`t`tres, err := tx.ExecContext(ctx, `n`t`t`t`INSERT INTO maintenance_tasks(hive_id, description, scheduled_for) VALUES(?,?,?)`,`n`t`t`tt.HiveID, t.Description, t.ScheduledFor)`n`t`tif err != nil {`n`t`t`tcontinue`n`t`t}`n`t`tn, _ := res.RowsAffected()`n`t`ttotal += n`n`t}`n`ttx.Commit()`n`treturn total, nil"}
)

# Bug-008: harvest_store (no copy) + harvest_service (sort in place)
New-BugBranches 8 @("internal/store/harvest_store.go","internal/service/harvest_service.go") @(
    @{File="internal/store/harvest_store.go"; Find=`t"out := make([]*model.Harvest, len(cached))`n`tcopy(out, cached)`n`treturn out, nil"; Replace=`t"return cached, nil"},
    @{File="internal/store/harvest_store.go"; Find=`t"out := make([]*model.Harvest, len(list))`n`tcopy(out, list)`n`treturn out, nil"; Replace=`t"return list, nil"},
    @{File="internal/service/harvest_service.go"; Find=`t"sorted := make([]*model.Harvest, len(list))`n`tcopy(sorted, list)`n`tsort.Slice(sorted, func(i, j int) bool {`n`t`treturn sorted[i].AmountKg > sorted[j].AmountKg`n`t})`n`treturn sorted, nil"; Replace=`t"sort.Slice(list, func(i, j int) bool {`n`t`treturn list[i].AmountKg > list[j].AmountKg`n`t})`n`treturn list, nil"}
)

# Bug-009: inspection_service (no validation) + inspection_handler (500)
New-BugBranches 9 @("internal/service/inspection_service.go","internal/handler/inspection_handler.go") @(
    @{File="internal/service/inspection_service.go"; Find=`t"ins, err := s.store.GetByID(ctx, id)`n`tif err != nil {`n`t`treturn fmt.Errorf(`"get inspection %d for complete: %w`", id, err)`n`t}`n`tif ins.Status == model.InspectionStatusCompleted {`n`t`treturn model.NewValidationError(`"status`", `"inspection already completed`")`n`t}`n`tif ins.Status != model.InspectionStatusPending && ins.Status != model.InspectionStatusOverdue {`n`t`treturn model.NewValidationError(`"status`", `"invalid status for completion: `"+ins.Status)`n`t}`n`treturn s.store.Complete"; Replace=`t"_, _ = s.store.GetByID(ctx, id)`n`treturn s.store.Complete"},
    @{File="internal/service/inspection_service.go"; Find=`t"context`n`t"fmt`n`t"time`""; Replace=`t"context`n`t"time`""},
    @{File="internal/handler/inspection_handler.go"; Find="writeError(w, err)`n`t`treturn"; Replace="writeJSON(w, http.StatusInternalServerError, map[string]string{`"error`": err.Error()})`n`t`treturn"}
)

# Bug-010: reading_service (UTC)
New-BugBranches 10 @("internal/service/reading_service.go") @(
    @{File="internal/service/reading_service.go"; Find=`t"sb.WriteString(`"id,hive_id,sensor_type,value,recorded_at`n`")`n`tfor _, r := range readings {`n`t`tsb.WriteString(fmt.Sprintf(`"%d,%d,%s,%.2f,%s`n`",`n`t`t`tr.ID, r.HiveID, r.SensorType, r.Value,`n`t`t`tr.RecordedAt.In(beeyardTZ).Format(`"2006-01-02 15:04:05`")))`n`t}"; Replace=`t"sb.WriteString(`"id,hive_id,sensor_type,value,timestamp`n`")`n`tfor _, r := range readings {`n`t`tts := r.RecordedAt.UTC().Format(`"2006-01-02 15:04:05`")`n`t`tsb.WriteString(fmt.Sprintf(`"%d,%d,%s,%.2f,%s`n`",`n`t`t`tr.ID, r.HiveID, r.SensorType, r.Value, ts))`n`t}"}
)

Write-Host "=== Verifying branches ==="
$branches = git branch
Write-Host "Branch count: $($branches.Count)"
git checkout main 2>&1
