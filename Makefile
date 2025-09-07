include .env

.PHONY: test

up:
	rm -f mongo.keyfile
	openssl rand -base64 24 > mongo.keyfile
	chmod 400 mongo.keyfile
	docker-compose up -d --build
	@echo ""
	@echo "👤 User Management API"
	@echo "    👉 Swagger UI:	http://localhost/user-management-api/swagger/index.html"
	@echo "    👉 gRPC UI:		http://localhost/user-management-api/grpcui/"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost/user-management-api/health"
	@echo ""
	@echo "🩺 Task Manager API"
	@echo "    👉 Swagger UI:	http://localhost/user-management-api/swagger/index.html"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost/task-manager-api/health"
	@echo ""
	@echo "📋 Health API"
	@echo "    👉 Swagger UI:	"
	@echo "    🔧 Command examples:"
	@echo "        curl "
	@echo ""
	@echo "🍃 Mongo Express:	http://localhost:${MONGO_EXPRESS_HOST_PORT} "
	@echo ""
down:
	docker-compose down

