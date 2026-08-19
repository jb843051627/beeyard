# Beeyard — 蜂场蜂箱监测管理系统

一个面向中小型蜂场的垂直领域后端服务，用于管理蜂场、蜂箱、蜂王、传感器读数、异常告警、巡检计划、采蜜记录、转场调度与维护任务。数据落盘 SQLite（纯 Go 驱动 `modernc.org/sqlite`），HTTP API 基于 Go 1.22 标准库 `net/http` 增强路由，附带轻量前端仪表盘。

## 业务领域

蜂场管理存在大量真实业务实体与跨模块流转：

- **蜂场 (Apiary)**：顶级组织单位，记录位置与建场时间。
- **蜂箱 (Hive)**：隶属于蜂场，有状态机（active / maintenance / inactive），安装后可挂载蜂王与传感器。
- **蜂王 (Queen)**：绑定到蜂箱，记录品种与标记时间，状态流转（active / superseded / dead）。
- **传感器读数 (SensorReading)**：温度 / 湿度 / 重量 三类，按蜂箱入库，提供批量导入与最新值缓存。
- **异常告警 (Alert)**：按蜂场+蜂箱产生，分级（info / warning / critical），状态机（active / acknowledged）。
- **巡检 (Inspection)**：按蜂箱排期，状态机（pending / completed / overdue），完成后回填记录。
- **维护任务 (MaintenanceTask)**：按蜂箱派发，pending / completed，可改期。
- **采蜜记录 (Harvest)**：按蜂场+蜂箱登记采蜜量，用于产量报表。
- **转场 (Transfer)**：蜂箱跨蜂场迁移，状态机（pending / in_transit / completed）。

## 架构

单体单进程，分层：

```
main.go                     # HTTP server bootstrap + 静态前端
internal/
  model/   # 领域实体与哨兵错误
  store/   # SQLite 持久化（每实体一个 store）
  cache/   # 读数缓存（并发安全）
  service/ # 业务逻辑（校验/状态机/聚合/批量导入/导出）
  handler/ # HTTP 路由与 JSON 编解码
web/                       # 前端仪表盘（HTML/JS/CSS，不计入 Go 代码量）
```

## HTTP 接口（5–20）

蜂箱 CRUD、读数上报/批量导入/导出、告警创建/确认/查询、巡检排期/完成/逾期列表、采蜜登记/报表、转场创建/分配/完成/日期范围查询、维护任务创建/完成/改期。

## 持久化

SQLite 文件模式（`*.db`），`database/sql` + `modernc.org/sqlite`，进程重启数据仍在。测试用 `t.TempDir()` 临时文件库。

## 运行

```bash
go build ./...
go run .            # 默认 :8080，自动建表
go test ./...
```
