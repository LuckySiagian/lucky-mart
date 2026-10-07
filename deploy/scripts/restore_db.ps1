# ==========================================
# PostgreSQL Automated Database Restore PowerShell Script
# ==========================================

param (
    [Parameter(Mandatory=$true)]
    [string]$BackupFile
)

$DbContainer = "sembako-postgres"
$DbUser = "sembako_usr"
$DbName = "sembako_db"

if (-not (Test-Path -Path $BackupFile)) {
    Write-Host "[ERROR] Backup file not found: $BackupFile" -ForegroundColor Red
    exit 1
}

Write-Host "[INFO] Restoring database $DbName from $BackupFile..." -ForegroundColor Cyan
Get-Content $BackupFile | docker exec -i $DbContainer psql -U $DbUser -d $DbName

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Database restoration completed successfully!" -ForegroundColor Green
} else {
    Write-Host "[ERROR] Database restoration failed!" -ForegroundColor Red
}
