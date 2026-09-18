#!/bin/bash
BACKUP_DIR="${1:-./backups}"
DB_PATH="${2:-./lms.db}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

mkdir -p "$BACKUP_DIR"
if [ -f "$DB_PATH" ]; then
    cp "$DB_PATH" "$BACKUP_DIR/lms_backup_${TIMESTAMP}.db"
    echo "[BACKUP] Successfully created database backup: $BACKUP_DIR/lms_backup_${TIMESTAMP}.db"
else
    echo "[WARNING] Database file $DB_PATH not found, backup skipped."
fi
