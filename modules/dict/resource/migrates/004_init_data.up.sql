-- 初始化种子数据（由 resource/initjson 生成，go run main.go migrate 自动入库）

INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (21,'2026-11-02 17:52:45','2026-11-02 17:52:45',19,'COOL','cool',1) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (22,'2026-11-02 17:52:45','2026-11-02 17:52:45',19,'闪酷','https://show.cool-admin.com/api/public/uploads/20230308/c731b0cba84046268b10edbbcf36f948_315c243a448e1369fa145c5ea3f020da.gif',2) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (23,'2026-11-02 17:52:45','2026-11-02 17:52:45',20,'法师','4',1) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (24,'2026-11-02 17:52:45','2026-11-02 17:52:45',20,'战士','3',2) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (25,'2026-11-02 17:52:45','2026-11-02 17:52:45',20,'坦克','2',3) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (26,'2026-11-02 17:52:45','2026-11-02 17:52:45',20,'刺客','1',4) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum") VALUES (27,'2026-11-02 17:52:45','2026-11-02 17:52:45',20,'射手','0',5) ON CONFLICT DO NOTHING;
INSERT INTO dict_info ("id","createTime","updateTime","typeId","name","value","orderNum","parentId") VALUES (30,'2026-11-02 17:52:45','2026-11-02 17:52:45',20,'幻影刺客','5',1,26) ON CONFLICT DO NOTHING;
INSERT INTO dict_type ("id","createTime","updateTime","name","key") VALUES (19,'2026-11-02 17:52:45','2026-11-02 17:52:45','品牌','brand') ON CONFLICT DO NOTHING;
INSERT INTO dict_type ("id","createTime","updateTime","name","key") VALUES (20,'2026-11-02 17:52:45','2026-11-02 17:52:45','职业','occupation') ON CONFLICT DO NOTHING;

-- 序列对齐
SELECT setval('dict_info_id_seq', COALESCE((SELECT max(id) FROM dict_info), 1));
SELECT setval('dict_type_id_seq', COALESCE((SELECT max(id) FROM dict_type), 1));
