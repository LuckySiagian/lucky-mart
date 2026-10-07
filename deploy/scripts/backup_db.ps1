# ==========================================
# PostgreSQL Automated Database Backup PowerShell Script
# ==========================================

$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$BackupDir = ".\backups"
$DbContainer = "sembako-postgres"
$DbUser = "sembako_usr"
$DbName = "sembako_db"
$BackupFile = "$BackupDir\backup_${DbName}_${Timestamp}.sql"

if (-not (Test-Path -Path $BackupDir)) {
    New-Item -ItemType Directory -Path $BackupDir | Out-Null
}

Write-Host "[INFO] Starting database backup for $DbName..." -ForegroundColor Cyan
docker exec -t $DbContainer pg_dump -U $DbUser $DbName > $BackupFile

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Database backup saved to $BackupFile" -ForegroundColor Green
} else {
    Write-Host "[ERROR] Database backup failed!" -ForegroundColor Red
}
