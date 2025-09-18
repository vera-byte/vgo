#!/bin/bash

# PostgreSQL到MySQL数据迁移脚本
# 作者: VGO Framework
# 功能: 将PostgreSQL数据库迁移到MySQL（使用pg_dump + 数据转换）
# 参数: 
#   $1: PostgreSQL连接字符串 (可选)
#   $2: MySQL连接字符串 (可选)
# 返回值: 0-成功, 1-失败

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要工具
check_tools() {
    local missing_tools=()
    
    if ! command -v pg_dump &> /dev/null; then
        missing_tools+=("pg_dump")
    fi
    
    if ! command -v mysql &> /dev/null; then
        missing_tools+=("mysql")
    fi
    
    if ! command -v sed &> /dev/null; then
        missing_tools+=("sed")
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_error "缺少必要工具: ${missing_tools[*]}"
        log_error "请安装PostgreSQL和MySQL客户端工具"
        exit 1
    fi
    
    log_info "工具检查完成"
}

# 设置默认环境变量
set_default_env() {
    export POSTGRES_HOST=${POSTGRES_HOST:-"localhost"}
    export POSTGRES_PORT=${POSTGRES_PORT:-"5432"}
    export POSTGRES_USER=${POSTGRES_USER:-"postgres"}
    export POSTGRES_PASS=${POSTGRES_PASS:-""}
    export POSTGRES_DB=${POSTGRES_DB:-""}
    
    export MYSQL_HOST=${MYSQL_HOST:-"localhost"}
    export MYSQL_PORT=${MYSQL_PORT:-"3306"}
    export MYSQL_USER=${MYSQL_USER:-"root"}
    export MYSQL_PASS=${MYSQL_PASS:-""}
    export MYSQL_DB=${MYSQL_DB:-""}
}

# 验证连接参数
validate_params() {
    if [[ -z "$POSTGRES_DB" ]]; then
        log_error "请设置POSTGRES_DB环境变量或提供PostgreSQL数据库名"
        exit 1
    fi
    
    if [[ -z "$MYSQL_DB" ]]; then
        log_error "请设置MYSQL_DB环境变量或提供MySQL数据库名"
        exit 1
    fi
}

# 测试PostgreSQL连接
test_postgres_connection() {
    log_info "测试PostgreSQL连接..."
    if PGPASSWORD="$POSTGRES_PASS" psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1;" &>/dev/null; then
        log_success "PostgreSQL连接成功"
    else
        log_error "PostgreSQL连接失败，请检查连接参数"
        exit 1
    fi
}

# 测试MySQL连接
test_mysql_connection() {
    log_info "测试MySQL连接..."
    if mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" -e "USE $MYSQL_DB;" 2>/dev/null; then
        log_success "MySQL连接成功"
    else
        log_error "MySQL连接失败，请检查连接参数"
        exit 1
    fi
}

# 备份MySQL数据库
backup_mysql() {
    if [[ "${BACKUP_BEFORE_MIGRATION:-true}" == "true" ]]; then
        log_info "备份MySQL数据库..."
        local backup_file="mysql_backup_$(date +%Y%m%d_%H%M%S).sql"
        if mysqldump -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" > "$backup_file"; then
            log_success "备份完成: $backup_file"
        else
            log_warning "备份失败，继续迁移"
        fi
    fi
}

# 导出PostgreSQL数据
export_postgres_data() {
    log_info "导出PostgreSQL数据..."
    local dump_file="postgres_dump_$(date +%Y%m%d_%H%M%S).sql"
    
    PGPASSWORD="$POSTGRES_PASS" pg_dump \
        -h "$POSTGRES_HOST" \
        -p "$POSTGRES_PORT" \
        -U "$POSTGRES_USER" \
        -d "$POSTGRES_DB" \
        --no-owner \
        --no-privileges \
        --clean \
        --if-exists \
        > "$dump_file"
    
    if [[ $? -eq 0 ]]; then
        log_success "PostgreSQL数据导出完成: $dump_file"
        echo "$dump_file"
    else
        log_error "PostgreSQL数据导出失败"
        exit 1
    fi
}

# 转换SQL语法
convert_sql_syntax() {
    local input_file="$1"
    local output_file="postgres_to_mysql_$(date +%Y%m%d_%H%M%S).sql"
    
    log_info "转换SQL语法..."
    
    # 创建转换后的SQL文件
    cat > "$output_file" << 'EOF'
-- PostgreSQL到MySQL转换后的SQL文件
-- 自动生成，请检查后使用

SET FOREIGN_KEY_CHECKS = 0;
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

EOF
    
    # 执行语法转换
    sed -E \
        -e 's/--.*$//' \
        -e '/^$/d' \
        -e 's/SERIAL/INT AUTO_INCREMENT/g' \
        -e 's/BIGSERIAL/BIGINT AUTO_INCREMENT/g' \
        -e 's/SMALLSERIAL/SMALLINT AUTO_INCREMENT/g' \
        -e 's/BOOLEAN/TINYINT(1)/g' \
        -e 's/TRUE/1/g' \
        -e 's/FALSE/0/g' \
        -e 's/BYTEA/LONGBLOB/g' \
        -e 's/TEXT/LONGTEXT/g' \
        -e 's/TIMESTAMP WITH TIME ZONE/DATETIME/g' \
        -e 's/TIMESTAMP WITHOUT TIME ZONE/DATETIME/g' \
        -e 's/TIMESTAMPTZ/DATETIME/g' \
        -e 's/UUID/VARCHAR(36)/g' \
        -e 's/JSONB/JSON/g' \
        -e 's/INET/VARCHAR(45)/g' \
        -e 's/CIDR/VARCHAR(45)/g' \
        -e 's/MACADDR/VARCHAR(17)/g' \
        -e 's/DOUBLE PRECISION/DOUBLE/g' \
        -e 's/REAL/FLOAT/g' \
        -e 's/SMALLINT/SMALLINT/g' \
        -e 's/BIGINT/BIGINT/g' \
        -e 's/INTEGER/INT/g' \
        -e 's/NUMERIC/DECIMAL/g' \
        -e 's/CHARACTER VARYING/VARCHAR/g' \
        -e 's/CHARACTER/CHAR/g' \
        -e 's/DROP SCHEMA IF EXISTS [^;]*;//g' \
        -e 's/CREATE SCHEMA [^;]*;//g' \
        -e 's/SET [^;]*;//g' \
        -e 's/SELECT pg_catalog[^;]*;//g' \
        -e 's/COMMENT ON [^;]*;//g' \
        -e 's/ALTER TABLE [^;]* OWNER TO [^;]*;//g' \
        -e 's/GRANT [^;]*;//g' \
        -e 's/REVOKE [^;]*;//g' \
        -e '/^\\./d' \
        -e 's/nextval\([^)]*\)/NULL/g' \
        -e 's/DEFAULT nextval\([^,)]*\)//g' \
        "$input_file" >> "$output_file"
    
    # 添加结尾
    cat >> "$output_file" << 'EOF'

SET FOREIGN_KEY_CHECKS = 1;
EOF
    
    log_success "SQL语法转换完成: $output_file"
    echo "$output_file"
}

# 导入MySQL数据
import_mysql_data() {
    local sql_file="$1"
    
    log_info "导入数据到MySQL..."
    
    # 清空目标数据库（可选）
    if [[ "${CLEAN_TARGET_DB:-false}" == "true" ]]; then
        log_warning "清空目标数据库..."
        mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" -e "DROP DATABASE IF EXISTS $MYSQL_DB; CREATE DATABASE $MYSQL_DB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    fi
    
    # 导入数据
    if mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" < "$sql_file"; then
        log_success "数据导入完成"
    else
        log_error "数据导入失败，请检查转换后的SQL文件"
        log_info "SQL文件位置: $sql_file"
        exit 1
    fi
}

# 验证迁移结果
verify_migration() {
    log_info "验证迁移结果..."
    
    # 检查表数量
    local postgres_tables=$(PGPASSWORD="$POSTGRES_PASS" psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" 2>/dev/null | tr -d ' ')
    local mysql_tables=$(mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" -D"$MYSQL_DB" -e "SHOW TABLES;" -s -N 2>/dev/null | wc -l)
    
    log_info "PostgreSQL表数量: $postgres_tables"
    log_info "MySQL表数量: $mysql_tables"
    
    if [[ "$postgres_tables" -eq "$mysql_tables" ]]; then
        log_success "表数量验证通过"
    else
        log_warning "表数量不匹配，请手动检查"
    fi
    
    log_warning "注意: PostgreSQL到MySQL的迁移可能需要手动调整以下内容:"
    log_warning "1. 序列(SERIAL)字段的AUTO_INCREMENT属性"
    log_warning "2. 索引和约束"
    log_warning "3. 存储过程和函数"
    log_warning "4. 特殊数据类型的数据"
}

# 清理临时文件
cleanup() {
    if [[ "${KEEP_TEMP_FILES:-false}" != "true" ]]; then
        log_info "清理临时文件..."
        rm -f postgres_dump_*.sql postgres_to_mysql_*.sql
        log_success "临时文件清理完成"
    else
        log_info "保留临时文件用于调试"
    fi
}

# 显示使用帮助
show_help() {
    echo "PostgreSQL到MySQL数据迁移脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [选项]"
    echo ""
    echo "环境变量:"
    echo "  POSTGRES_HOST   PostgreSQL主机地址 (默认: localhost)"
    echo "  POSTGRES_PORT   PostgreSQL端口 (默认: 5432)"
    echo "  POSTGRES_USER   PostgreSQL用户名 (默认: postgres)"
    echo "  POSTGRES_PASS   PostgreSQL密码"
    echo "  POSTGRES_DB     PostgreSQL数据库名 (必需)"
    echo ""
    echo "  MYSQL_HOST      MySQL主机地址 (默认: localhost)"
    echo "  MYSQL_PORT      MySQL端口 (默认: 3306)"
    echo "  MYSQL_USER      MySQL用户名 (默认: root)"
    echo "  MYSQL_PASS      MySQL密码"
    echo "  MYSQL_DB        MySQL数据库名 (必需)"
    echo ""
    echo "  BACKUP_BEFORE_MIGRATION  迁移前是否备份 (默认: true)"
    echo "  CLEAN_TARGET_DB         是否清空目标数据库 (默认: false)"
    echo "  KEEP_TEMP_FILES         是否保留临时文件 (默认: false)"
    echo ""
    echo "示例:"
    echo "  POSTGRES_DB=myapp MYSQL_DB=myapp $0"
    echo ""
    echo "选项:"
    echo "  -h, --help      显示此帮助信息"
}

# 主函数
main() {
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
    esac
    
    log_info "开始PostgreSQL到MySQL数据迁移"
    
    check_tools
    set_default_env
    validate_params
    test_postgres_connection
    test_mysql_connection
    backup_mysql
    
    # 执行迁移流程
    local dump_file=$(export_postgres_data)
    local converted_file=$(convert_sql_syntax "$dump_file")
    import_mysql_data "$converted_file"
    verify_migration
    
    # 清理
    cleanup
    
    log_success "迁移流程完成！"
    log_warning "请手动验证数据完整性和业务逻辑"
}

# 执行主函数
main "$@"