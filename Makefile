.PHONY: all build frontend clean dev

all: build

frontend:
	cd frontend && npm install && npm run build

build: frontend
	go build -o gekka-dashboard .

clean:
	rm -rf static/ gekka-dashboard
	cd frontend && rm -rf node_modules

dev:
	cd frontend && npm run dev
