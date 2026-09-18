# LMS Database Automated Backup Utility
param (
    [string]$BackupDir = "./backups",
    [string]$DbPath = "./lms.db"
)

$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
if (!(Test-Path $BackupDir)) {
    New-Item -ItemType Directory -Path $BackupDir | Out-Null
}

$target = Join-Path $BackupDir "lms_backup_$timestamp.db"
if (Test-Path $DbPath) {
    Copy-Item -Path $DbPath -Destination $target
    Write-Output "[BACKUP] Successfully created database backup: $target"
} else {
    Write-Output "[WARNING] Database file $DbPath not found, backup skipped."
}
