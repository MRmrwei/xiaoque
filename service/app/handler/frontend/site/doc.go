// Package site 提供官网前台页面渲染相关的 HTTP 控制器。
//
// 文件划分：
//   - page.go 页面渲染：按路径聚合页面、导航、站点配置与内容集合
//
// 控制器只负责参数解析与统一响应，业务逻辑在 app/services/site/impl，
// 数据访问在 app/repository/page、app/repository/menu、app/repository/setting。
package site
