.PHONY: help build up down restart logs clean prune ps stats health \
        backend-logs frontend-logs nginx-logs postgres-logs es-logs minio-logs \
        db-shell db-backup db-restore ssl-cert infra-up prod-up

help:
	@echo "Suipic Docker Management"
	@echo ""
	@echo "Service Management:"
	@echo "  build         - Build all Docker images"
	@echo "  up            - Start all services"
	@echo "  down          - Stop all services"
	@echo "  restart       - Restart all services"
	@echo "  ps            - Show running containers"
	@echo "  stats         - Show container resource usage"
	@echo "  health        - Check health of all services"
	@echo ""
	@echo "Logs:"
	@echo "  logs          - Show logs from all services"
	@echo "  backend-logs  - Show backend logs"
	@echo "  frontend-logs - Show frontend logs"
	@echo "  nginx-logs    - Show nginx logs"
	@echo "  postgres-logs - Show PostgreSQL logs"
	@echo "  es-logs       - Show Elasticsearch logs"
	@echo "  minio-logs    - Show MinIO logs"
	@echo ""
	@echo "Database:"
	@echo "  db-shell      - Open PostgreSQL shell"
	@echo "  db-backup     - Backup database to backup.sql"
	@echo "  db-restore    - Restore database from backup.sql"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean         - Stop and remove containers, networks"
	@echo "  prune         - Clean and remove volumes (WARNING: destroys data)"
	@echo ""
	@echo "SSL/Security:"
	@echo "  ssl-cert      - Generate self-signed SSL certificate"
	@echo ""
	@echo "Deployment:"
	@echo "  infra-up      - Start infrastructure services only (for local dev)"
	@echo "  prod-up       - Start with production configuration"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

ps:
	docker-compose ps

stats:
	docker stats

health:
	@echo "Checking service health..."
	@docker-compose ps
	@echo ""
	@echo "PostgreSQL:"
	@docker-compose exec -T postgres pg_isready -U suipic || echo "PostgreSQL not ready"
	@echo ""
	@echo "Elasticsearch:"
	@curl -s http://localhost:9200/_cluster/health?pretty || echo "Elasticsearch not ready"
	@echo ""
	@echo "MinIO:"
	@curl -s http://localhost:9000/minio/health/live || echo "MinIO not ready"

clean:
	docker-compose down

prune:
	docker-compose down -v

backend-logs:
	docker-compose logs -f backend

frontend-logs:
	docker-compose logs -f frontend

nginx-logs:
	docker-compose logs -f nginx

postgres-logs:
	docker-compose logs -f postgres

es-logs:
	docker-compose logs -f elasticsearch

minio-logs:
	docker-compose logs -f minio

db-shell:
	docker-compose exec postgres psql -U suipic -d suipic

db-backup:
	@echo "Backing up database to backup.sql..."
	@docker-compose exec -T postgres pg_dump -U suipic suipic > backup.sql
	@echo "Backup completed: backup.sql"

db-restore:
	@echo "Restoring database from backup.sql..."
	@docker-compose exec -T postgres psql -U suipic -d suipic < backup.sql
	@echo "Restore completed"

ssl-cert:
	@mkdir -p docker/nginx/ssl
	@openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout docker/nginx/ssl/key.pem \
		-out docker/nginx/ssl/cert.pem \
		-subj "/C=US/ST=State/L=City/O=Suipic/CN=localhost"
	@echo "SSL certificate generated in docker/nginx/ssl/"
	@echo "Update docker/nginx/conf.d/default.conf to enable HTTPS"

infra-up:
	docker-compose -f docker-compose.infra.yml up -d

prod-up:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
