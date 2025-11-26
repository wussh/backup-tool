#!/bin/bash
################################################################################
# FLEXIBLE PostgreSQL Backup Script
# Supports: docker run, docker exec, kubectl exec
################################################################################

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Load environment variables
source .env

# Set defaults
BACKUP_METHOD=${BACKUP_METHOD:-docker-run}
BACKUP_TEMP_DIR=${BACKUP_TEMP_DIR:-/tmp/db-backups}
DATE=$(date +%F_%H-%M-%S)

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}  PostgreSQL Flexible Backup${NC}"
echo -e "${BLUE}  Method: ${CYAN}$BACKUP_METHOD${NC}"
echo -e "${BLUE}  Timestamp: $DATE${NC}"
echo -e "${BLUE}=========================================${NC}"

mkdir -p backup/postgres
BACKUP_FILE="backup/postgres/${PG_DB}_${DATE}.sql"

echo "Starting PostgreSQL backup..."
echo "  Host: $PG_HOST"
echo "  Database: $PG_DB"
echo "  User: $PG_USER"
echo ""

case $BACKUP_METHOD in
    docker-run)
        echo "Using docker run method..."
        docker run --rm \
          -e PGPASSWORD="$PG_PASS" \
          -v "$(pwd)/backup/postgres:/backup" \
          postgres:${PG_VERSION:-15} \
          pg_dump -h "$PG_HOST" -U "$PG_USER" "$PG_DB" \
          > "$BACKUP_FILE"
        ;;
        
    docker-exec)
        echo "Using docker exec method..."
        echo "  Container: $PG_CONTAINER"
        docker exec "$PG_CONTAINER" \
          sh -c "PGPASSWORD='$PG_PASS' pg_dump -h localhost -U $PG_USER $PG_DB" \
          > "$BACKUP_FILE"
        ;;
        
    kubectl-exec)
        echo "Using kubectl exec method..."
        echo "  Pod: $PG_POD"
        echo "  Namespace: ${K8S_NAMESPACE:-default}"
        kubectl exec -n "${K8S_NAMESPACE:-default}" "$PG_POD" -- \
          sh -c "PGPASSWORD='$PG_PASS' pg_dump -h localhost -U $PG_USER $PG_DB" \
          > "$BACKUP_FILE"
        ;;
        
    *)
        echo -e "${RED}✗ Unknown backup method: $BACKUP_METHOD${NC}"
        echo "Valid methods: docker-run, docker-exec, kubectl-exec"
        exit 1
        ;;
esac

if [ -f "$BACKUP_FILE" ] && [ -s "$BACKUP_FILE" ]; then
    SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo ""
    echo -e "${GREEN}✓ Backup completed successfully!${NC}"
    echo "  File: $BACKUP_FILE"
    echo "  Size: $SIZE"
else
    echo -e "${RED}✗ Backup failed or empty file created${NC}"
    exit 1
fi
