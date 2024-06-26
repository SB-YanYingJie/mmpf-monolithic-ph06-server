# --- utility ---
# command defined 'xxxx: guard_yyy' requires variable name 'yyy'
# refer to: https://dev.classmethod.jp/articles/lambda-deploy-with-make/
guard_%:
	@ if [ "$($(*))" = "" ]; then \
		echo "require value '$*' like 'make command $*=xxxx'"; \
		exit 1; \
	fi

login:
	echo $(GITHUB_ACCESS_TOKEN) | docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin

#----api test----#
run_mock_request:
	go clean -testcache
	go test apitest/api_test.go -v

#----mock-server---#
IMAGE=ghcr.io/machinemapplatform/mmpf-monolithic/mock-server
build_image_ssh:
	@docker build -t $(IMAGE) -f ./Dockerfile --build-arg default .
build_image_https:
	@docker build -t $(IMAGE) -f ./Dockerfile --build-arg con=https .
push_image: guard_tag
	@docker tag $(IMAGE) $(IMAGE):$(tag)
	@docker push $(IMAGE):$(tag)


#----ecs fagate-----#
CLUSTER_NAME=mmpf-mock-cluster
SERVICE_NAME=mmpf-mock-service
ecs_create_cluster: guard_env
	@aws ecs --profile $(env) create-cluster \
		--cluster-name $(CLUSTER_NAME) \
		--settings "name=containerInsights,value=disabled"

ecs_register_task: guard_env
	@aws ecs --profile $(env) register-task-definition \
		--cli-input-json file://deployment/definition/task.json \
	| jq '.taskDefinition.revision'

ecs_create_service: guard_env
	@aws ecs --profile $(env) create-service \
		--service-name $(SERVICE_NAME) \
		--launch-type FARGATE \
		--cli-input-json file://deployment/definition/service.json

ecs_update_service: guard_env
	@aws ecs --profile $(env) update-service \
		--service $(SERVICE_NAME) \
		--cli-input-json file://deployment/definition/service.json

ecs_delete_service: guard_env
	@aws ecs --profile $(env) delete-service \
		--cluster $(CLUSTER_NAME) \
		--service $(SERVICE_NAME) \
		--force

# --- iam role for ecs ---
TASK_ROLE_NAME=MockServerTaskRole
EXECUTION_POLICY_NAME=MockServerRestrictedExecution
EXECUTION_ROLE_NAME=MockServerExecutionRole
create_task_role: guard_env
	@aws iam --profile $(env) create-role \
		--path /$(env)/mmpf-monolithic/ \
		--role-name $(TASK_ROLE_NAME) \
		--assume-role-policy-document file://deployment/iam/assume-role-policy.json

create_execution_role: guard_env
	@aws iam --profile $(env) create-policy \
		--path /$(env)/mmpf-monolithic/ \
		--policy-name $(EXECUTION_POLICY_NAME) \
		--policy-document file://deployment/iam/execution-role.json
	@aws iam --profile $(env) create-role \
		--path /$(env)/mmpf-monolithic/ \
		--role-name $(EXECUTION_ROLE_NAME) \
		--assume-role-policy-document file://deployment/iam/assume-role-policy.json
	export ACCOUNT_ID=$$(aws sts get-caller-identity --profile $(env) --query Account --output text); \
	aws iam --profile $(env) attach-role-policy \
		--role-name $(EXECUTION_ROLE_NAME) \
		--policy-arn arn:aws:iam::$${ACCOUNT_ID}:policy/$(env)/mmpf-monolithic/$(EXECUTION_POLICY_NAME)

