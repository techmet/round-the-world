FRONTEND_PATH = $(PWD)/frontend
BACKEND_PATH = $(PWD)/backend

prepare-frontend:
	@if [ -d "$(FRONTEND_PATH)" ]; then cd $(FRONTEND_PATH) && yarn install; fi
prepare-backend:
	@if [ -d "$(BACKEND_PATH)" ]; then cd $(BACKEND_PATH) && go mod tidy; fi	
run-frontend:
	@if [ -d "$(FRONTEND_PATH)" ]; then cd $(FRONTEND_PATH) && yarn start; fi
run-backend:
	@if [ -d "$(BACKEND_PATH)" ]; then cd $(BACKEND_PATH) && go run cmd/main.go; fi	