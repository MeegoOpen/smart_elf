**项目概览**
- 插件名称：`meego-app-robot-to-workitem-2.0`
- 目标：提供“小精灵配置”页面与控件展示，使用 Semi UI 组件与图标库，接入 Meego/Lark JS SDK 与后端接口。


**前置要求**
- 已安装 Node.js 16+ 与 Yarn。
- 安装 Lark Project CLI：
  - 全局安装：`npm i -g @lark-project/cli`
  - 或使用 npx：`npx @lark-project/cli <command>`

**快速开始**
- 安装依赖：`yarn`
- 开发启动：`yarn dev`
  - 脚本等价于：`LPM_ORIGIN_FOR_DEV_SERVER=10.92.190.41 node scripts.js start`
  - 依赖 `lpm` 命令，如报错 `lpm: command not found` 请安装 CLI。
- 构建：`yarn build`
- 发布：`node scripts.js deploy <token>`
  - 注意：`package.json` 中 `yarn deploy` 指向 `node scripts deploy`（缺少 `.js`）。建议直接使用上面的命令或修正脚本配置。

**配置文件说明**
- `plugin.config.json`（根目录）
  - 启动与构建直接读取该文件，不再从 `config/env-config.json` 按环境生成。
  - 必填字段：
    - `siteDomain`：插件所在站点域名（用于请求与控制台跳转）。
    - `pluginId`：插件 ID（用于构建/发布展示）。
  - 示例：
    ```json
    {
      "siteDomain": "https://project.feishu.cn",
      "pluginId": "your_plugin_id"
    }
    ```

**常用命令**
- `yarn dev`：本地开发，启动 `lpm start` 并注入 `LPM_ORIGIN_FOR_DEV_SERVER`。
- `yarn build`：构建产物，内部调用 `lpm build`。
- `node scripts.js deploy <token>`：发布版本，内部调用 `lpm release <token>`。

**项目结构**
- `src/page/config/`：配置页面（表单使用 Semi UI 的 `Form`、`Card`、`Skeleton`、`Banner`、`Toast` 等）。
- `src/page/control/`：控件展示页（使用 `Spin` 与 `@douyinfe/semi-icons`）。
- `src/api/`：接口请求封装（`axios` 拦截器、`services` 工具函数）。
- `src/constants/`：常量（从 `plugin.config.json` 读取 `pluginId` 与 `siteDomain`）。
- `src/utils/`：通用工具（`sdk` 适配、存储与错误处理）。
- `styles/semi-theme-overrides.scss`：Semi UI 主题细节覆盖。



**开发提示**
- 表单 API 通过 `Form` 的 `getFormApi={(api) => ...}` 获取，支持 `setValues`、`validate`、`setValue`、`getValue`。
- Select 组件使用 `optionList` 字段（Semi UI）而非 `options`。

