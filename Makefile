include .env

.PHONY: test

up:
	cd health-api && go mod vendor
	cd task-manager-api && go mod vendor
	rm -f mongo.keyfile
	openssl rand -base64 24 > mongo.keyfile
	chmod 400 mongo.keyfile
	docker compose up -d --build
	@echo ""
	@echo "📋 Health Service"
	@echo "    👉 Swagger UI:	http://localhost:${NGINX_HOST_HTTP_PORT}/health-api/v1/swagger/index.html"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost:${NGINX_HOST_HTTP_PORT}/health-api/v1/health"
	@echo ""
	@echo "🩺 Task Manager API"
	@echo "    👉 Swagger UI:	http://localhost:${NGINX_HOST_HTTP_PORT}/task-manager-api/v1/swagger/index.html"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost:${NGINX_HOST_HTTP_PORT}/task-manager-api/v1/health"
	@echo ""
	@echo "👤 User Management API"
	@echo "    👉 Swagger UI:	http://localhost:${NGINX_HOST_HTTP_PORT}/user-management-api/v1/swagger/index.html"
	@echo "    👉 gRPC UI:		http://localhost:${NGINX_HOST_HTTP_PORT}/user-management-api/v1/grpcui/"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost:${NGINX_HOST_HTTP_PORT}/user-management-api/v1/health"
	@echo ""
	@echo "🍃 Mongo Express:	http://localhost:${MONGO_EXPRESS_HOST_PORT}"
	@echo ""
down:
	docker compose down
all-test-unit:
	$(MAKE) -C common test-unit & \
	$(MAKE) -C health-api test-unit & \
	$(MAKE) -C task-manager-api test-unit & \
	wait
test-e2e:
	cd test && go test -race e2e/*.go