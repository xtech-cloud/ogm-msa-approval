APP_NAME := omo-msa-approval
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

.PHONY: build
build:
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rf /tmp/msa-approval.db

.PHONY: call
call:
	# -------------------------------------------------------------------------
	# 创建工作流, 缺少参数
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Make
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Make '{"name":"test"}'
	# 创建工作流
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Make '{"name":"test1", "mode": 1}'
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Make '{"name":"test2", "mode": 1}'
	# 创建工作流，已存在
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Make '{"name":"test1", "mode": 1}'
	# 列举工作流，无参数
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.List
	# 列举工作流
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.List '{"offset":1, "count":1}'
	# 获取工作流，无参数
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Get
	# 获取工作流，不存在
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Get '{"uuid":"0000000"}'
	# 获取工作流
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Get '{"uuid":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Get '{"name":"test2"}'
	# 操作员加入
	MICRO_REGISTRY=consul micro call omo.msa.approval Operator.Join '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	# 操作员加入, 已存在
	MICRO_REGISTRY=consul micro call omo.msa.approval Operator.Join '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	# 操作员离开
	MICRO_REGISTRY=consul micro call omo.msa.approval Operator.Leave '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	# 操作员离开, 不存在
	MICRO_REGISTRY=consul micro call omo.msa.approval Operator.Leave '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	# 删除工作流，无参数
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Remove
	# 删除工作流，不存在
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Remove '{"uuid":"00000"}'
	# 删除工作流，存在
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Remove '{"uuid":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.msa.approval Workflow.Remove '{"uuid":"ad0234829205b9033196ba818f7a872b"}'

.PHONY: tcall
tcall:
	mkdir -p ./bin
	go build -o ./bin/ ./tester
	./bin/tester

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo-msa-startkit:latest
