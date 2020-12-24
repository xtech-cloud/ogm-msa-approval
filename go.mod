module omo-msa-approval

go 1.13

require (
	github.com/micro/go-micro/v2 v2.3.0
	github.com/micro/go-plugins/config/source/consul/v2 v2.3.0
	github.com/micro/go-plugins/logger/logrus/v2 v2.3.0
	github.com/micro/go-plugins/registry/consul/v2 v2.3.0
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.3.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/xtech-cloud/omo-msp-approval v1.16.0
	github.com/xtech-cloud/omo-msp-notification v1.0.0
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/mysql v1.0.1
	gorm.io/driver/sqlite v1.1.1
	gorm.io/gorm v1.20.0
)
