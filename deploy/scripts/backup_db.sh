#!/bin/bash

# ==========================================
# PostgreSQL Automated Database Backup Script
# ==========================================

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_DIR="./backups"
DB_CONTAINER="sembako-postgres"
DB_USER="${DB_USER:-sembako_usr}"
DB_NAME="${DB_NAME:-sembako_db}"
BACKUP_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${TIMESTAMP}.sql.gz"

mkdir -p "$BACKUP_DIR"

echo "[INFO] Starting database backup for $DB_NAME..."
docker exec -t $DB_CONTAINER pg_dump -U $DB_USER $DB_NAME | gzip > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
  echo "[SUCCESS] Backup saved to $BACKUP_FILE"
else
  echo "[ERROR] Database backup failed!"
  exit 1
fi
