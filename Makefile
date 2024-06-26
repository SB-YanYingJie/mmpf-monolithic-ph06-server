SHELL=/bin/bash

# --- run lint ---
# Support of generics golangci/golangci-lint#2649
run_lint:
	golangci-lint run --go=1.18 ./...

# --- testing ---
run_unit_test:
	go test -v ./...

run_e2e_test:
	go mod tidy
	go build -o main ./cmd/main.go
	./main &
	sleep 10
	go test /app/internal/e2e/ --run $(TEST) -tags=e2e_test
	pkill main

run_e2e_test__mono:
	@make run_e2e_test \
	TEST=Test_Slam_Mono* \
	NUMBER_OF_LENSES=mono \
	SEND_IMAGE_TYPE=mono \
	MMID=testMMID \
	REDIS_PUBSUB_CHANNEL_POSE=pose_testMMID \
	KD_CALIB_PATH=/app/lib/calibration_sample.ini \
	KD_MAP_PATH=/app/lib/sample_map.kdmp

run_e2e_test__stereo_merged:
	@make run_e2e_test \
	TEST=Test_Slam_Stereo_Merged_Success \
	NUMBER_OF_LENSES=stereo \
	SEND_IMAGE_TYPE=stereo_merged \
	MMID=testMMID \
	REDIS_PUBSUB_CHANNEL_POSE=pose_testMMID \
	TRIMMING_PARAMETER=0:0:640:400 \
	KD_CALIB_PATH=/app/lib/calib_LeadSense282_640x400.ini \
	KD_MAP_PATH=/app/slam-service/lib/hall.kdmp

run_e2e_test__stereo_separated:
	@make run_e2e_test \
	TEST=Test_Slam_Stereo_Separated_Success \
	NUMBER_OF_LENSES=stereo \
	SEND_IMAGE_TYPE=stereo_separated \
	MMID=testMMID \
	REDIS_PUBSUB_CHANNEL_POSE=pose_testMMID \
	KD_CALIB_PATH=/app/lib/calib_LeadSense282_640x400.ini \
	KD_MAP_PATH=/app/slam-service/lib/hall.kdmp

# --- for ci ----
DOCKER_COMPOSE_FILE=./docker-compose.ci.yml
DOCKERFILE=./Dockerfile.ci
CI_IMAGE=mmpf_monolithic_app_for_ci
WORK_DIR=/app

docker_compose_up:
	docker build --secret id=token,src=.token -f $(DOCKERFILE) -t $(CI_IMAGE) .
	docker-compose -f $(DOCKER_COMPOSE_FILE)  build
	docker-compose -f $(DOCKER_COMPOSE_FILE)  up -d

docker_compose_down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

run_unit_test_in_cicontainer:
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T -w $(WORK_DIR) app go mod tidy -compat=1.18
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T -w $(WORK_DIR) app go test -coverprofile coverage.out -covermode atomic ./...

run_lint_in_cicontainer:
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T -w $(WORK_DIR) app go mod tidy -compat=1.18
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T -w $(WORK_DIR) app golangci-lint run --go=1.18 ./...

# --- container image ---
IMAGE=ghcr.io/machinemapplatform/mmpf-monolithic
TAG=latest

git_url := $(shell git config --get remote.origin.url)

login:
	echo $(GITHUB_ACCESS_TOKEN) | docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin

build_image:
	@if [[ $(git_url) == git@* ]]; then\
		make build_image_ssh;\
	else\
		make build_image_https;\
	fi

build_image_ssh:
	docker build --ssh default=${HOME}/.ssh/id_rsa -t $(IMAGE) .;
	@make tag_image

build_image_https:
	docker build --build-arg con=https -t $(IMAGE) . --no-cache;
	@make tag_image


tag_image:
	docker tag $(IMAGE) $(IMAGE):$(TAG)

push_image:
	docker push $(IMAGE):$(TAG)
