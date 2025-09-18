#!/bin/bash

# MySQL到PostgreSQL数据迁移脚本
# 作者: VGO Framework
# 功能: 使用pgloader将MySQL数据库迁移到PostgreSQL
# 参数: 
#   $1: MySQL连接字符串 (可选)
#   $2: PostgreSQL连接字符串 (可选)
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

# 检查pgloader是否安装
check_pgloader() {
    if ! command -v pgloader &> /dev/null; then
        log_error "pgloader未安装，请先安装: brew install pgloader"
        exit 1
    fi
    log_info "pgloader版本: $(pgloader --version | head -1)"
}

# 设置默认环境变量
set_default_env() {
    export MYSQL_HOST=${MYSQL_HOST:-"localhost"}
    export MYSQL_PORT=${MYSQL_PORT:-"3306"}
    export MYSQL_USER=${MYSQL_USER:-"root"}
    export MYSQL_PASS=${MYSQL_PASS:-""}
    export MYSQL_DB=${MYSQL_DB:-""}
    
    export POSTGRES_HOST=${POSTGRES_HOST:-"localhost"}
    export POSTGRES_PORT=${POSTGRES_PORT:-"5432"}
    export POSTGRES_USER=${POSTGRES_USER:-"postgres"}
    export POSTGRES_PASS=${POSTGRES_PASS:-""}
    export POSTGRES_DB=${POSTGRES_DB:-""}
}

# 验证连接参数
validate_params() {
    if [[ -z "$MYSQL_DB" ]]; then
        log_error "请设置MYSQL_DB环境变量或提供MySQL数据库名"
        exit 1
    fi
    
    if [[ -z "$POSTGRES_DB" ]]; then
        log_error "请设置POSTGRES_DB环境变量或提供PostgreSQL数据库名"
        exit 1
    fi
}

# 测试MySQL连接
test_mysql_connection() {
    log_info "测试MySQL连接..."
    if command -v mysql &> /dev/null; then
        if mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" -e "USE $MYSQL_DB;" 2>/dev/null; then
            log_success "MySQL连接成功"
        else
            log_error "MySQL连接失败，请检查连接参数"
            exit 1
        fi
    else
        log_warning "mysql客户端未安装，跳过连接测试"
    fi
}

# 测试PostgreSQL连接
test_postgres_connection() {
    log_info "测试PostgreSQL连接..."
    if command -v psql &> /dev/null; then
        if PGPASSWORD="$POSTGRES_PASS" psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1;" &>/dev/null; then
            log_success "PostgreSQL连接成功"
        else
            log_error "PostgreSQL连接失败，请检查连接参数"
            exit 1
        fi
    else
        log_warning "psql客户端未安装，跳过连接测试"
    fi
}

# 备份PostgreSQL数据库
backup_postgres() {
    if [[ "${BACKUP_BEFORE_MIGRATION:-true}" == "true" ]]; then
        log_info "备份PostgreSQL数据库..."
        local backup_file="postgres_backup_$(date +%Y%m%d_%H%M%S).sql"
        if command -v pg_dump &> /dev/null; then
            PGPASSWORD="$POSTGRES_PASS" pg_dump -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" "$POSTGRES_DB" > "$backup_file"
            log_success "备份完成: $backup_file"
        else
            log_warning "pg_dump未安装，跳过备份"
        fi
    fi
}

# 执行迁移
run_migration() {
    log_info "开始数据迁移..."
    log_info "源数据库: MySQL - $MYSQL_HOST:$MYSQL_PORT/$MYSQL_DB"
    log_info "目标数据库: PostgreSQL - $POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB"
    
    local config_file="$(dirname "$0")/mysql_to_postgres.load"
    
    if [[ ! -f "$config_file" ]]; then
        log_error "配置文件不存在: $config_file"
        exit 1
    fi
    
    # 执行pgloader
    if pgloader "$config_file"; then
        log_success "数据迁移完成！"
    else
        log_error "数据迁移失败"
        exit 1
    fi
}

# 验证迁移结果
verify_migration() {
    log_info "验证迁移结果..."
    
    # 检查表数量
    if command -v mysql &> /dev/null && command -v psql &> /dev/null; then
        local mysql_tables=$(mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" -D"$MYSQL_DB" -e "SHOW TABLES;" -s -N 2>/dev/null | wc -l)
        local postgres_tables=$(PGPASSWORD="$POSTGRES_PASS" psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" 2>/dev/null | tr -d ' ')
        
        log_info "MySQL表数量: $mysql_tables"
        log_info "PostgreSQL表数量: $postgres_tables"
        
        if [[ "$mysql_tables" -eq "$postgres_tables" ]]; then
            log_success "表数量验证通过"
        else
            log_warning "表数量不匹配，请手动检查"
        fi
    else
        log_warning "缺少数据库客户端，跳过验证"
    fi
}

# 显示使用帮助
show_help() {
    echo "MySQL到PostgreSQL数据迁移脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [选项]"
    echo ""
    echo "环境变量:"
    echo "  MYSQL_HOST      MySQL主机地址 (默认: localhost)"
    echo "  MYSQL_PORT      MySQL端口 (默认: 3306)"
    echo "  MYSQL_USER      MySQL用户名 (默认: root)"
    echo "  MYSQL_PASS      MySQL密码"
    echo "  MYSQL_DB        MySQL数据库名 (必需)"
    echo ""
    echo "  POSTGRES_HOST   PostgreSQL主机地址 (默认: localhost)"
    echo "  POSTGRES_PORT   PostgreSQL端口 (默认: 5432)"
    echo "  POSTGRES_USER   PostgreSQL用户名 (默认: postgres)"
    echo "  POSTGRES_PASS   PostgreSQL密码"
    echo "  POSTGRES_DB     PostgreSQL数据库名 (必需)"
    echo ""
    echo "  BACKUP_BEFORE_MIGRATION  迁移前是否备份 (默认: true)"
    echo ""
    echo "示例:"
    echo "  MYSQL_DB=myapp POSTGRES_DB=myapp $0"
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
    
    log_info "开始MySQL到PostgreSQL数据迁移"
    
    check_pgloader
    set_default_env
    validate_params
    test_mysql_connection
    test_postgres_connection
    backup_postgres
    run_migration
    verify_migration
    
    log_success "迁移流程完成！"
}

# 执行主函数
main "$@"