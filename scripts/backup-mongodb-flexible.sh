#!/bin/bash
################################################################################
# FLEXIBLE MongoDB Backup Script
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
echo -e "${BLUE}  MongoDB Flexible Backup${NC}"
echo -e "${BLUE}  Method: ${CYAN}$BACKUP_METHOD${NC}"
echo -e "${BLUE}  Timestamp: $DATE${NC}"
echo -e "${BLUE}=========================================${NC}"

mkdir -p backup/mongodb
BACKUP_DIR="backup/mongodb/${DATE}"

echo "Starting MongoDB backup..."
echo "  Host: $MONGO_HOST"
echo "  Database: $MONGO_DB"
echo ""

case $BACKUP_METHOD in
    docker-run)
        echo "Using docker run method..."
        docker run --rm \
          -v "$(pwd)/backup/mongodb:/backup" \
          mongo:${MONGO_VERSION:-7} \
          mongodump --host "$MONGO_HOST" --db "$MONGO_DB" --out "/backup/$DATE"
        ;;
        
    docker-exec)
        echo "Using docker exec method..."
        echo "  Container: $MONGO_CONTAINER"
        # Create backup inside container
        docker exec "$MONGO_CONTAINER" \
          mongodump --host localhost --db "$MONGO_DB" --out "$BACKUP_TEMP_DIR/$DATE"
        # Copy backup from container to host
        mkdir -p "$BACKUP_DIR"
        docker cp "$MONGO_CONTAINER:$BACKUP_TEMP_DIR/$DATE/$MONGO_DB" "$BACKUP_DIR/"
        # Cleanup inside container
        docker exec "$MONGO_CONTAINER" rm -rf "$BACKUP_TEMP_DIR/$DATE"
        ;;
        
    kubectl-exec)
        echo "Using kubectl exec method..."
        echo "  Pod: $MONGO_POD"
        echo "  Namespace: ${K8S_NAMESPACE:-default}"
        # Create backup inside pod
        kubectl exec -n "${K8S_NAMESPACE:-default}" "$MONGO_POD" -- \
          mongodump --host localhost --db "$MONGO_DB" --out "$BACKUP_TEMP_DIR/$DATE"
        # Copy backup from pod to host
        mkdir -p "$BACKUP_DIR"
        kubectl cp "${K8S_NAMESPACE:-default}/$MONGO_POD:$BACKUP_TEMP_DIR/$DATE/$MONGO_DB" "$BACKUP_DIR/"
        # Cleanup inside pod
        kubectl exec -n "${K8S_NAMESPACE:-default}" "$MONGO_POD" -- \
          rm -rf "$BACKUP_TEMP_DIR/$DATE"
        ;;
        
    *)
        echo -e "${RED}✗ Unknown backup method: $BACKUP_METHOD${NC}"
        echo "Valid methods: docker-run, docker-exec, kubectl-exec"
        exit 1
        ;;
esac

if [ -d "$BACKUP_DIR" ]; then
    SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
    echo ""
    echo -e "${GREEN}✓ Backup completed successfully!${NC}"
    echo "  Directory: $BACKUP_DIR"
    echo "  Size: $SIZE"
else
    echo -e "${RED}✗ Backup failed${NC}"
    exit 1
fi
