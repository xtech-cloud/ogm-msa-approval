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
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Healthy.Echo '{"msg":"hello"}'
	echo # -------------------------------------------------------------------------
	echo #  参数完整性测试
	echo # -------------------------------------------------------------------------
	echo # 创建工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make '{"name":"test"}'
	echo # 获取工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Get
	echo # 删除工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Remove
	echo # 删除工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Remove '{"uuid":"00000"}'
	echo # -------------------------------------------------------------------------
	echo #  代码覆盖测试 - 全票工作流
	echo # -------------------------------------------------------------------------
	echo # 创建工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make '{"name":"test1", "mode": 1}'
	echo # 获取工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Get '{"uuid":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Get '{"name":"test1"}'
	echo # 操作员加入
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.Join '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchJoin '{"operator":["002", "003"], "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchJoin '{"operator":["003", "004"], "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	echo # 列举操作员
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.List '{"workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	echo # 过滤工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.Filter '{"operator":"001"}'
	echo # 提交任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Submit '{"subject":"subject-1", "body": "boby-1", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	echo # 列举任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	echo # 通过任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"8a8c26af4df5ac2c617688fd84218899", "operator":"001"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"8a8c26af4df5ac2c617688fd84218899", "operator":"002"}'
	echo # 查询记录
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Action.Query '{"task":"8a8c26af4df5ac2c617688fd84218899"}'
	echo # 通过任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"8a8c26af4df5ac2c617688fd84218899", "operator":"003"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"8a8c26af4df5ac2c617688fd84218899", "operator":"004"}'
	echo # 查询记录
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Action.Query '{"task":"8a8c26af4df5ac2c617688fd84218899"}'
	echo # 列举任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	echo # 拒绝任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"8a8c26af4df5ac2c617688fd84218899", "operator":"001", "reason":"this is reason"}'
	echo # 查询记录
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Action.Query '{"task":"8a8c26af4df5ac2c617688fd84218899"}'
	echo # 列举任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	echo # 搜索任务
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Search '{"workflow":"test"}'
	echo # 操作员离开
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.Leave '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchLeave '{"operator":["003", "004"], "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchLeave '{"operator":["002", "003"], "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	echo # -------------------------------------------------------------------------
	echo #  代码覆盖测试 - 任一工作流
	echo # -------------------------------------------------------------------------
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make '{"name":"test2", "mode": 2}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.List
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchJoin '{"operator":["a", "b", "c"], "workflow":"ad0234829205b9033196ba818f7a872b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Submit '{"subject":"subject-2", "body": "boby-2", "workflow":"ad0234829205b9033196ba818f7a872b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"b6a3001a1408e614a2a6fe7bb7e208a5", "operator":"a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"b6a3001a1408e614a2a6fe7bb7e208a5", "operator":"b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"b6a3001a1408e614a2a6fe7bb7e208a5", "operator":"c"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"b6a3001a1408e614a2a6fe7bb7e208a5", "operator":"c"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"b6a3001a1408e614a2a6fe7bb7e208a5", "operator":"b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"b6a3001a1408e614a2a6fe7bb7e208a5", "operator":"a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	echo # -------------------------------------------------------------------------
	echo #  代码覆盖测试 - 半票工作流
	echo # -------------------------------------------------------------------------
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make '{"name":"test3", "mode": 3}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.List
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchJoin '{"operator":["a", "b", "c"], "workflow":"8ad8757baa8564dc136c1e07507f4a98"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Submit '{"subject":"subject-3", "body": "boby-3", "workflow":"8ad8757baa8564dc136c1e07507f4a98"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"245be4cc1f88180bb9bb7e86f4f9330b", "operator":"a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"245be4cc1f88180bb9bb7e86f4f9330b", "operator":"b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Accept '{"uuid":"245be4cc1f88180bb9bb7e86f4f9330b", "operator":"c"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"245be4cc1f88180bb9bb7e86f4f9330b", "operator":"c"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"245be4cc1f88180bb9bb7e86f4f9330b", "operator":"b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.Reject '{"uuid":"245be4cc1f88180bb9bb7e86f4f9330b", "operator":"a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Task.List 
	echo # -------------------------------------------------------------------------
	echo #  缺省参数测试
	echo # -------------------------------------------------------------------------
	# 列举工作流，无参数
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.List
	# 列举工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.List '{"offset":1, "count":1}'
	# 获取工作流，不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Get '{"uuid":"0000000"}'
	echo # -------------------------------------------------------------------------
	echo #  冲突测试
	echo # -------------------------------------------------------------------------
	# 创建工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make '{"name":"test1", "mode": 1}'
	# 操作员加入
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.Join '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	# 操作员离开, 不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.Leave '{"operator":"001", "workflow":"5a105e8b9d40e1329780d62ea2265d8a"}'
	# 删除工作流，存在
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Remove '{"uuid":"5a105e8b9d40e1329780d62ea2265d8a"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Remove '{"uuid":"ad0234829205b9033196ba818f7a872b"}'
	# @过半工作流
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Workflow.Make '{"name":"test2", "mode": 2}'
	MICRO_REGISTRY=consul micro call omo.api.msa.approval Operator.BatchJoin '{"operator":["A", "B"], "workflow":"ad0234829205b9033196ba818f7a872b"}'

.PHONY: post
post:
	curl -X POST -d '{"msg":"hello"}' 127.0.0.1:8080/msa/approval/Healthy/Echo

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
