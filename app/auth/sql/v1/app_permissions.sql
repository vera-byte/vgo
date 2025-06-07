CREATE TABLE app_permissions (
    id           SERIAL PRIMARY KEY,
    app_id       VARCHAR(64) NOT NULL,
    method_id    INTEGER NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(app_id, method_id),
    FOREIGN KEY (app_id) REFERENCES apps(app_id) ON DELETE CASCADE,
    FOREIGN KEY (method_id) REFERENCES rpc_methods(id) ON DELETE CASCADE
);

-- 添加权限（为该 app 授权前两个接口）
INSERT INTO app_permissions (app_id, method_id)
SELECT 'com.demo.app', id FROM rpc_methods
WHERE full_method IN (
    '/grpc.auth.v1.Greeter/SayHello',
    '/grpc.auth.v1.UserService/GetUser'
);