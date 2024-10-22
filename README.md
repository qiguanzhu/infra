# infra
infra of projects

/seele 基础功能接口定义 + 基础结构体定义 + 基础 func 定义

/nerv 基础功能的具体实现，可能同一类接口存在多种实现，包含mock。
如：xsql；fastsql；xstat 等

/nerv/magi 微小功能的实现集合入口，未拆分独立功能模块

/tokyo3 常规项目 service 类的常见接口定义。eg：User，Sso，Order 等
如果引用方需要增加新函数，可以用 interface 包含引用

/evangelion tokyo3 接口的具体实现方案，按照定义提供基础功能。对于复杂的逻辑采用 func 传入的方式注入相关逻辑。
主体流程保持原状。
如果需要自定义相关复杂流程，可以通过 interface 引用的方式自行实现。
主要用于实现基础方案 and 提供样例

/evangelion/repo db repo 具体实现

/evangelion/persistence db repo 默认链接，持久化调用
两种方案，一种是注册成 xdb 全局生效；另一种是注册到repo中，需要初始化 Newxx() 放到 persistence 路径下；service中引用的是 persistence 中持久化的数据库 repo

/evangelion/service service 层级的接口函数功能实现入口

エヴァ贰號機 出動！

/pkg 通用功能函数封装入口

/KommSüsserTod docs

/s2organ S2机关，跑通一个项目的基础，aka /cmd or init.go
go run 时 调用相关内容

/prototype 原型机
repo & service 的 mock 内容，all cases tested