# Gost Panel

<div align="center">

**现代化 Gost v3 端口转发管理面板**

[![License](https://img.shields.io/github/license/code-gopher/gostPanel)](./LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://golang.org/)

</div>


---
# 推广 🚀 高速 Hysteria2 接入服务 | 极致性价比

基于最新的 **Hysteria2 (Hy2)** 协议构建，专为开发者、科研人员及跨境电商优化。晚高峰不降速，4K 视频秒开。

## ✨ 核心优势

* **极速协议**：全线采用 **Hysteria2** 协议，抗封锁性强，延迟极低。
* **节点丰富**：覆盖日本、香港、新加坡、美国等优质 IDC 线路，实测延迟低至 **70ms**。
* **极致性价比**：月付仅需 **￥3.00** 起，并提供 **不限时** 流量包套餐。
* **全方位解锁**：完美支持 ChatGPT、Netflix、Disney+、TikTok 等流媒体。
* **不限设备**：全计划均不限制连接设备数量，支持多端共用。
* 
## 🔗 立即接入

**官方入口地址**（手动去掉“去”字或直接复制）：

> **http://aip.acgf.vip:15432**
---

---

## 🚀 核心特性

- ⚡ **极简管理** - 统一管理多个 Gost 节点及其转发规则。
- 🔗 **多跳隧道** - 支持复杂的出口/中转隧道配置（入口 -> 出口 -> 目标）。
- 📊 **实时延迟** - 每 5 秒自动检测节点与转发链路的 Ping 值。
- 🛡️ **安全可靠** - 支持 API 认证，内置操作审计日志。
- 🐳 **快速部署** - 全量支持 Docker 一键运行及 GitHub Action 自动构建。

## 🛠️ 快速开始

### 1. 部署面板

**方式 A: Docker (推荐)**
确保已安装 Docker 和 Docker Compose，然后运行：
```bash
curl -sSL https://raw.githubusercontent.com/code-gopher/gostPanel/master/docker-compose.yml -o docker-compose.yml
docker-compose up -d
```

**方式 B: 二进制脚本 (原生部署)**
适用于没有 Docker 的 Linux 服务器（要求内存 > 128MB）：

**默认安装 (端口 39100):**
```bash
bash <(curl -sSL https://raw.githubusercontent.com/code-gopher/gostPanel/master/scripts/install_panel.sh)
```

**自定义端口 (例如 8080):**
```bash
bash <(curl -sSL https://raw.githubusercontent.com/code-gopher/gostPanel/master/scripts/install_panel.sh) 8080
```
或者
```bash
bash <(curl -sSL https://raw.githubusercontent.com/code-gopher/gostPanel/master/scripts/install_panel.sh) install --port 8080
```

**卸载面板:**
```bash
bash <(curl -sSL https://raw.githubusercontent.com/code-gopher/gostPanel/master/scripts/install_panel.sh) uninstall
```

**卸载节点:**
```bash
bash <(curl -sSL https://raw.githubusercontent.com/code-gopher/gostPanel/master/scripts/install_node.sh) uninstall
```

---

## 📸 界面预览

<table>
  <tr>
    <td align="center">
      <h3>仪表板</h3>
      <img src="./docs/screenshots/dash.png" alt="仪表板" width="400"/>
    </td>
    <td align="center">
      <h3>节点管理</h3>
      <img src="./docs/screenshots/node.png" alt="节点管理" width="400"/>
    </td>
  </tr>
  <tr>
    <td align="center">
      <h3>转发管理</h3>
      <img src="./docs/screenshots/forwards.png" alt="转发管理" width="400"/>
    </td>
    <td align="center">
      <h3>隧道管理</h3>
      <img src="./docs/screenshots/tunnels.png" alt="隧道管理" width="400"/>
    </td>
  </tr>
</table>

---

## 📡 添加转发节点

1. 在面板导航至 **节点管理**。
2. 点击已有节点的 **安装** 按钮，或者点击 **添加节点**。
3. 复制生成的安装命令，在目标服务器上运行即可完成自动化部署。

---

## 📦 预编译下载

项目支持多平台二进制发布，请访问 [Releases](https://github.com/code-gopher/gostPanel/releases) 下载：
- **Linux**: amd64, arm64
- **Windows**: amd64
- **macOS**: Apple Silicon, Intel

### 开发与编译
如果你想自行编译单二进制文件：

**使用 Makefile（推荐）**：
```bash
# 完整构建（前端+后端）
make build

# 只构建前端
make build-web

# 只构建后端
make build-server

# 清理构建产物
make clean

# 构建多平台发布版本
make release
```

---

## 🤝 声明
本项目仅供学习与科研使用，请在法律范围内使用。

**Star 也是一种支持！⭐**
