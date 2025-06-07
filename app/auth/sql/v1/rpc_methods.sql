CREATE TABLE rpc_methods (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(128) NOT NULL,                -- 接口名，如：SayHello
    full_method  VARCHAR(255) NOT NULL UNIQUE,         -- gRPC 方法路径，如：/package.Service/Method
    description  TEXT,                                 -- 接口说明（可选）
    enabled      BOOLEAN DEFAULT TRUE,                 -- 是否启用
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- 添加 RPC 接口
INSERT INTO rpc_methods (name, full_method, description, enabled) VALUES
('SayHello', '/grpc.auth.v1.Greeter/SayHello', '打招呼接口', true),
('GetUser', '/grpc.auth.v1.UserService/GetUser', '获取用户信息', true),
('DeleteUser', '/grpc.auth.v1.UserService/DeleteUser', '删除用户接口', false); -- 模拟禁用状态
