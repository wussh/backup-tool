#!/bin/bash
################################################################################
# FLEXIBLE Database Backup Script - All Databases
# Supports: PostgreSQL, MySQL, MariaDB, MongoDB
# Methods: docker run, docker exec, kubectl exec
################################################################################

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Load environment variables
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    echo "Please create .env file with database credentials."
    echo "You can copy .env.example and modify it."
    exit 1
fi

source .env

# Set default backup method if not specified
BACKUP_METHOD=${BACKUP_METHOD:-docker-run}
BACKUP_TEMP_DIR=${BACKUP_TEMP_DIR:-/tmp/db-backups}

# Generate timestamp
DATE=$(date +%F_%H-%M-%S)

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Flexible Database Backup Automation${NC}"
echo -e "${BLUE}  Method: ${CYAN}$BACKUP_METHOD${NC}"
echo -e "${BLUE}  Timestamp: $DATE${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Create backup directories
echo -e "${YELLOW}Creating backup directories...${NC}"
mkdir -p backup/postgres backup/mysql backup/mariadb backup/mongodb
echo -e "${GREEN}✓ Directories created${NC}"
echo ""

################################################################################
# Helper Functions
################################################################################

# Function to backup PostgreSQL
backup_postgres() {
    if [ -z "$PG_HOST" ] || [ -z "$PG_DB" ]; then
        echo -e "${YELLOW}[PostgreSQL] Skipped (no configuration)${NC}"
        echo ""
        return
    fi

    echo -e "${BLUE}[PostgreSQL] Starting backup...${NC}"
    echo "  Method: $BACKUP_METHOD"
    echo "  Host: $PG_HOST"
    echo "  Database: $PG_DB"
    echo "  User: $PG_USER"
    
    BACKUP_FILE="backup/postgres/${PG_DB}_${DATE}.sql"
    
    case $BACKUP_METHOD in
        docker-run)
            # Original method: docker run
            docker run --rm \
              -e PGPASSWORD="$PG_PASS" \
              -v "$(pwd)/backup/postgres:/backup" \
              postgres:${PG_VERSION:-15} \
              pg_dump -h "$PG_HOST" -U "$PG_USER" "$PG_DB" \
              > "$BACKUP_FILE"
            ;;
            
        docker-exec)
            # Exec into existing container
            if [ -z "$PG_CONTAINER" ]; then
                echo -e "${RED}✗ PG_CONTAINER not set in .env${NC}"
                return 1
            fi
            echo "  Container: $PG_CONTAINER"
            docker exec "$PG_CONTAINER" \
              sh -c "PGPASSWORD='$PG_PASS' pg_dump -h localhost -U $PG_USER $PG_DB" \
              > "$BACKUP_FILE"
            ;;
            
        kubectl-exec)
            # Exec into Kubernetes pod
            if [ -z "$PG_POD" ]; then
                echo -e "${RED}✗ PG_POD not set in .env${NC}"
                return 1
            fi
            echo "  Pod: $PG_POD"
            echo "  Namespace: $K8S_NAMESPACE"
            kubectl exec -n "${K8S_NAMESPACE:-default}" "$PG_POD" -- \
              sh -c "PGPASSWORD='$PG_PASS' pg_dump -h localhost -U $PG_USER $PG_DB" \
              > "$BACKUP_FILE"
            ;;
            
        *)
            echo -e "${RED}✗ Unknown backup method: $BACKUP_METHOD${NC}"
            return 1
            ;;
    esac
    
    if [ -f "$BACKUP_FILE" ] && [ -s "$BACKUP_FILE" ]; then
        SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        echo -e "${GREEN}✓ PostgreSQL backup completed: $BACKUP_FILE ($SIZE)${NC}"
    else
        echo -e "${RED}✗ PostgreSQL backup failed or empty${NC}"
    fi
    echo ""
}

# Function to backup MySQL
backup_mysql() {
    if [ -z "$MYSQL_HOST" ] || [ -z "$MYSQL_DB" ]; then
        echo -e "${YELLOW}[MySQL] Skipped (no configuration)${NC}"
        echo ""
        return
    fi

    echo -e "${BLUE}[MySQL] Starting backup...${NC}"
    echo "  Method: $BACKUP_METHOD"
    echo "  Host: $MYSQL_HOST"
    echo "  Database: $MYSQL_DB"
    echo "  User: $MYSQL_USER"
    
    BACKUP_FILE="backup/mysql/${MYSQL_DB}_${DATE}.sql"
    
    case $BACKUP_METHOD in
        docker-run)
            docker run --rm \
              -v "$(pwd)/backup/mysql:/backup" \
              mysql:${MYSQL_VERSION:-8} \
              sh -c "mysqldump -h$MYSQL_HOST -u$MYSQL_USER -p$MYSQL_PASS $MYSQL_DB" \
              > "$BACKUP_FILE"
            ;;
            
        docker-exec)
            if [ -z "$MYSQL_CONTAINER" ]; then
                echo -e "${RED}✗ MYSQL_CONTAINER not set in .env${NC}"
                return 1
            fi
            echo "  Container: $MYSQL_CONTAINER"
            docker exec "$MYSQL_CONTAINER" \
              sh -c "mysqldump -h localhost -u$MYSQL_USER -p$MYSQL_PASS $MYSQL_DB" \
              > "$BACKUP_FILE"
            ;;
            
        kubectl-exec)
            if [ -z "$MYSQL_POD" ]; then
                echo -e "${RED}✗ MYSQL_POD not set in .env${NC}"
                return 1
            fi
            echo "  Pod: $MYSQL_POD"
            echo "  Namespace: $K8S_NAMESPACE"
            kubectl exec -n "${K8S_NAMESPACE:-default}" "$MYSQL_POD" -- \
              sh -c "mysqldump -h localhost -u$MYSQL_USER -p$MYSQL_PASS $MYSQL_DB" \
              > "$BACKUP_FILE"
            ;;
            
        *)
            echo -e "${RED}✗ Unknown backup method: $BACKUP_METHOD${NC}"
            return 1
            ;;
    esac
    
    if [ -f "$BACKUP_FILE" ] && [ -s "$BACKUP_FILE" ]; then
        SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        echo -e "${GREEN}✓ MySQL backup completed: $BACKUP_FILE ($SIZE)${NC}"
    else
        echo -e "${RED}✗ MySQL backup failed or empty${NC}"
    fi
    echo ""
}

# Function to backup MariaDB
backup_mariadb() {
    if [ -z "$MARIADB_HOST" ] || [ -z "$MARIADB_DB" ]; then
        echo -e "${YELLOW}[MariaDB] Skipped (no configuration)${NC}"
        echo ""
        return
    fi

    echo -e "${BLUE}[MariaDB] Starting backup...${NC}"
    echo "  Method: $BACKUP_METHOD"
    echo "  Host: $MARIADB_HOST"
    echo "  Database: $MARIADB_DB"
    echo "  User: $MARIADB_USER"
    
    BACKUP_FILE="backup/mariadb/${MARIADB_DB}_${DATE}.sql"
    
    case $BACKUP_METHOD in
        docker-run)
            docker run --rm \
              -v "$(pwd)/backup/mariadb:/backup" \
              mariadb:${MARIADB_VERSION:-11} \
              sh -c "mysqldump -h$MARIADB_HOST -u$MARIADB_USER -p$MARIADB_PASS $MARIADB_DB" \
              > "$BACKUP_FILE"
            ;;
            
        docker-exec)
            if [ -z "$MARIADB_CONTAINER" ]; then
                echo -e "${RED}✗ MARIADB_CONTAINER not set in .env${NC}"
                return 1
            fi
            echo "  Container: $MARIADB_CONTAINER"
            docker exec "$MARIADB_CONTAINER" \
              sh -c "mysqldump -h localhost -u$MARIADB_USER -p$MARIADB_PASS $MARIADB_DB" \
              > "$BACKUP_FILE"
            ;;
            
        kubectl-exec)
            if [ -z "$MARIADB_POD" ]; then
                echo -e "${RED}✗ MARIADB_POD not set in .env${NC}"
                return 1
            fi
            echo "  Pod: $MARIADB_POD"
            echo "  Namespace: $K8S_NAMESPACE"
            kubectl exec -n "${K8S_NAMESPACE:-default}" "$MARIADB_POD" -- \
              sh -c "mysqldump -h localhost -u$MARIADB_USER -p$MARIADB_PASS $MARIADB_DB" \
              > "$BACKUP_FILE"
            ;;
            
        *)
            echo -e "${RED}✗ Unknown backup method: $BACKUP_METHOD${NC}"
            return 1
            ;;
    esac
    
    if [ -f "$BACKUP_FILE" ] && [ -s "$BACKUP_FILE" ]; then
        SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        echo -e "${GREEN}✓ MariaDB backup completed: $BACKUP_FILE ($SIZE)${NC}"
    else
        echo -e "${RED}✗ MariaDB backup failed or empty${NC}"
    fi
    echo ""
}

# Function to backup MongoDB
backup_mongodb() {
    if [ -z "$MONGO_HOST" ] || [ -z "$MONGO_DB" ]; then
        echo -e "${YELLOW}[MongoDB] Skipped (no configuration)${NC}"
        echo ""
        return
    fi

    echo -e "${BLUE}[MongoDB] Starting backup...${NC}"
    echo "  Method: $BACKUP_METHOD"
    echo "  Host: $MONGO_HOST"
    echo "  Database: $MONGO_DB"
    
    BACKUP_DIR="backup/mongodb/${DATE}"
    
    case $BACKUP_METHOD in
        docker-run)
            docker run --rm \
              -v "$(pwd)/backup/mongodb:/backup" \
              mongo:${MONGO_VERSION:-7} \
              mongodump --host "$MONGO_HOST" --db "$MONGO_DB" --out "/backup/$DATE"
            ;;
            
        docker-exec)
            if [ -z "$MONGO_CONTAINER" ]; then
                echo -e "${RED}✗ MONGO_CONTAINER not set in .env${NC}"
                return 1
            fi
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
            if [ -z "$MONGO_POD" ]; then
                echo -e "${RED}✗ MONGO_POD not set in .env${NC}"
                return 1
            fi
            echo "  Pod: $MONGO_POD"
            echo "  Namespace: $K8S_NAMESPACE"
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
            return 1
            ;;
    esac
    
    if [ -d "$BACKUP_DIR" ]; then
        SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
        echo -e "${GREEN}✓ MongoDB backup completed: $BACKUP_DIR ($SIZE)${NC}"
    else
        echo -e "${RED}✗ MongoDB backup failed${NC}"
    fi
    echo ""
}

################################################################################
# Run Backups
################################################################################

backup_postgres
backup_mysql
backup_mariadb
backup_mongodb

################################################################################
# Summary
################################################################################
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Backup process completed!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Backup method: $BACKUP_METHOD"
echo "Backup location: $(pwd)/backup/"
echo ""
echo "Recent backups:"
find backup -type f -o -type d -name "*.sql" -o -name "*${DATE}*" 2>/dev/null | head -20
echo ""
