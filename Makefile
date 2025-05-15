.PHONY: test coverage coverage-html clean proxymock-mock run build integration-test load-test

# Define proxymock environment variables
PROXYMOCK_ENV = http_proxy=socks5h://localhost:4140 \
                https_proxy=socks5h://localhost:4140 \
                SSL_CERT_FILE=~/.speedscale/certs/tls.crt

# Find first recording directory
PROXYMOCK_RECORDING := $(shell find ./proxymock -name "recorded-*" -type d | head -n 1)

test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build:
	go build -o outerspace-go main.go

run:
	go run main.go

clean:
	rm -f coverage.out coverage.html outerspace-go
	rm -rf logs
	rm -rf proxymock/mocked-* proxymock/replayed-*

integration-test: build proxymock-mock
	mkdir -p logs
	echo "Starting outerspace-go in background with proxymock..."
	$(PROXYMOCK_ENV) ./outerspace-go > logs/outerspace.log 2>&1 & echo $$! > logs/outerspace.pid
	echo "Waiting for outerspace-go to start..."
	sleep 2
	echo "Running integration tests with proxymock..."
	proxymock replay --in $(PROXYMOCK_RECORDING) --fail-if requests.response-pct!=100
	echo "Cleaning up..."
	pkill -f "outerspace-go" || true
	pkill -f "proxymock" || true
	echo "Integration tests completed. See logs in the logs directory."

load-test: build proxymock-mock
	mkdir -p logs
	echo "Starting outerspace-go in background with proxymock..."
	$(PROXYMOCK_ENV) ./outerspace-go > logs/outerspace.log 2>&1 & echo $$! > logs/outerspace.pid
	echo "Waiting for outerspace-go to start..."
	sleep 2
	echo "Running load tests with proxymock..."
	proxymock replay --in $(PROXYMOCK_RECORDING) --vus 10 --for 1m --fail-if "latency.p95 > 200"
	echo "Cleaning up..."
	pkill -f "outerspace-go" || true
	pkill -f "proxymock" || true
	echo "Load tests completed. See logs in the logs directory."

proxymock-mock:
	mkdir -p logs
	nohup proxymock mock --in $(PROXYMOCK_RECORDING) > logs/proxymock-mock.log 2>&1 & \
	sleep 2
	@if ! pgrep -f "proxymock mock" > /dev/null; then \
		echo "Error: Proxymock is NOT mocking!"; \
		cat logs/proxymock-mock.log; \
		exit 1; \
	fi
	@echo "Proxymock started successfully."
