#!/bin/bash

# ==========================================
# PostgreSQL Automated Database Restore Script
# ==========================================

if [ -z "$1" ]; then
  echo "Usage: ./restore_db.sh <path_to_backup_file.sql.gz>"
  exit 1
fi

BACKUP_FILE="$1"
DB_CONTAINER="sembako-postgres"
DB_USER="${DB_USER:-sembako_usr}"
DB_NAME="${DB_NAME:-sembako_db}"

echo "[INFO] Restoring database $DB_NAME from $BACKUP_FILE..."
gunzip -c "$BACKUP_FILE" | docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME

if [ $? -eq 0 ]; then
  echo "[SUCCESS] Database restoration completed successfully!"
else
  echo "[ERROR] Database restoration failed!"
  exit 1
fi
