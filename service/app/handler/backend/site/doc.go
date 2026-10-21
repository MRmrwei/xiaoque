// Package site 提供管理后台站点编排相关的 HTTP 控制器。
// 对应后台「站点编排」的四个 Tab：页面内容、顶部导航、底部导航、站点配置。
//
// 文件划分：
//   - page.go    页面内容：列表、新建、详情、保存草稿、删除、翻译、发布 / 下架、预览令牌
//   - menu.go    顶部 / 底部导航：菜单列表、菜单详情、整体保存菜单项
//   - setting.go 站点配置：读取与保存站点设置
//
// 控制器只负责参数解析与统一响应，业务逻辑在 app/services/site，
// 数据访问在 app/repository/page、app/repository/menu、app/repository/setting。
package site
