# 数据库迁移工具使用说明

VGO框架提供了PostgreSQL和MySQL之间的双向数据迁移工具，支持开发和生产环境的数据库切换需求。

## 🛠️ 工具概览

### 1. MySQL → PostgreSQL 迁移
- **工具**: pgloader (推荐)
- **配置文件**: `mysql_to_postgres.load`
- **执行脚本**: `migrate_mysql_to_postgres.sh`

### 2. PostgreSQL → MySQL 迁移
- **工具**: pg_dump + 语法转换
- **执行脚本**: `migrate_postgres_to_mysql.sh`

## 📋 前置要求

### macOS 环境
```bash
# 安装pgloader
brew install pgloader

# 安装PostgreSQL客户端
brew install postgresql

# 安装MySQL客户端
brew install mysql
```

### Linux 环境
```bash
# Ubuntu/Debian
sudo apt-get install pgloader postgresql-client mysql-client

# CentOS/RHEL
sudo yum install pgloader postgresql mysql
```

## 🚀 使用方法

### MySQL → PostgreSQL

#### 方法1: 使用环境变量
```bash
# 设置环境变量
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASS=your_password
export MYSQL_DB=source_database

export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASS=your_password
export POSTGRES_DB=target_database

# 执行迁移
./scripts/migrate_mysql_to_postgres.sh
```

#### 方法2: 直接使用pgloader
```bash
# 修改配置文件中的连接信息
vim scripts/mysql_to_postgres.load

# 执行迁移
pgloader scripts/mysql_to_postgres.load
```

### PostgreSQL → MySQL

```bash
# 设置环境变量
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASS=your_password
export POSTGRES_DB=source_database

export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASS=your_password
export MYSQL_DB=target_database

# 执行迁移
./scripts/migrate_postgres_to_mysql.sh
```

## ⚙️ 配置选项

### 通用环境变量
| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `BACKUP_BEFORE_MIGRATION` | 迁移前是否备份目标数据库 | `true` |
| `CLEAN_TARGET_DB` | 是否清空目标数据库 | `false` |
| `KEEP_TEMP_FILES` | 是否保留临时文件 | `false` |

### MySQL连接参数
| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `MYSQL_HOST` | MySQL主机地址 | `localhost` |
| `MYSQL_PORT` | MySQL端口 | `3306` |
| `MYSQL_USER` | MySQL用户名 | `root` |
| `MYSQL_PASS` | MySQL密码 | (空) |
| `MYSQL_DB` | MySQL数据库名 | (必需) |

### PostgreSQL连接参数
| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `POSTGRES_HOST` | PostgreSQL主机地址 | `localhost` |
| `POSTGRES_PORT` | PostgreSQL端口 | `5432` |
| `POSTGRES_USER` | PostgreSQL用户名 | `postgres` |
| `POSTGRES_PASS` | PostgreSQL密码 | (空) |
| `POSTGRES_DB` | PostgreSQL数据库名 | (必需) |

## 📊 数据类型映射

### MySQL → PostgreSQL
| MySQL类型 | PostgreSQL类型 | 说明 |
|-----------|----------------|------|
| `TINYINT(1)` | `BOOLEAN` | 布尔值 |
| `AUTO_INCREMENT` | `SERIAL` | 自增字段 |
| `DATETIME` | `TIMESTAMP` | 时间戳 |
| `LONGTEXT` | `TEXT` | 长文本 |
| `JSON` | `JSONB` | JSON数据 |

### PostgreSQL → MySQL
| PostgreSQL类型 | MySQL类型 | 说明 |
|----------------|-----------|------|
| `SERIAL` | `INT AUTO_INCREMENT` | 自增字段 |
| `BIGSERIAL` | `BIGINT AUTO_INCREMENT` | 大整数自增 |
| `BOOLEAN` | `TINYINT(1)` | 布尔值 |
| `TEXT` | `LONGTEXT` | 长文本 |
| `JSONB` | `JSON` | JSON数据 |
| `UUID` | `VARCHAR(36)` | UUID字符串 |
| `TIMESTAMPTZ` | `DATETIME` | 时间戳 |

## ⚠️ 注意事项

### MySQL → PostgreSQL
1. **字符编码**: 确保两个数据库都使用UTF-8编码
2. **自增字段**: MySQL的AUTO_INCREMENT会自动转换为PostgreSQL的SERIAL
3. **索引**: 大部分索引会自动创建，但可能需要手动优化
4. **外键**: 外键约束会被保留

### PostgreSQL → MySQL
1. **数据类型**: 某些PostgreSQL特有类型需要手动调整
2. **序列**: SERIAL字段转换为AUTO_INCREMENT可能需要手动调整起始值
3. **函数和存储过程**: 需要手动重写
4. **索引**: 部分索引类型MySQL不支持，需要手动调整

## 🔍 故障排除

### 常见问题

#### 1. 连接失败
```bash
# 检查数据库服务是否运行
sudo systemctl status postgresql
sudo systemctl status mysql

# 检查防火墙设置
sudo ufw status
```

#### 2. 权限问题
```bash
# PostgreSQL权限
GRANT ALL PRIVILEGES ON DATABASE mydb TO myuser;

# MySQL权限
GRANT ALL PRIVILEGES ON mydb.* TO 'myuser'@'%';
FLUSH PRIVILEGES;
```

#### 3. 字符编码问题
```sql
-- MySQL设置UTF-8
ALTER DATABASE mydb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- PostgreSQL设置UTF-8
CREATE DATABASE mydb WITH ENCODING 'UTF8';
```

### 日志分析
迁移脚本会生成详细的日志信息，包括：
- 连接测试结果
- 数据导出/导入进度
- 错误信息和警告
- 验证结果

## 📝 最佳实践

1. **备份**: 迁移前务必备份源数据库和目标数据库
2. **测试**: 先在测试环境验证迁移脚本
3. **验证**: 迁移后检查数据完整性和业务逻辑
4. **监控**: 关注迁移过程中的性能和错误日志
5. **回滚**: 准备回滚方案以防迁移失败

## 🔧 自定义配置

### 修改pgloader配置
编辑 `mysql_to_postgres.load` 文件，可以自定义：
- 内存使用量
- 数据类型转换规则
- 迁移前后的SQL命令
- 包含/排除的表

### 扩展转换规则
在 `migrate_postgres_to_mysql.sh` 中的 `convert_sql_syntax` 函数中添加更多转换规则。

## 📞 技术支持

如果遇到问题，请：
1. 检查日志文件
2. 验证数据库连接参数
3. 确认数据库版本兼容性
4. 查看VGO框架文档: https://goframe.org/

---

**注意**: 数据库迁移是一个复杂的过程，建议在生产环境使用前充分测试。