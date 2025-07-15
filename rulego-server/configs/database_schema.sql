-- RuleGo IoT数据处理系统数据库表结构
-- 创建数据库
CREATE DATABASE IF NOT EXISTS rulego_test;

-- 使用数据库
\c rulego_test;

-- 创建传感器数据表
CREATE TABLE IF NOT EXISTS sensor_data (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    temperature_celsius DECIMAL(5,2),
    temperature_fahrenheit DECIMAL(5,2),
    humidity_percentage DECIMAL(5,2),
    pressure_hpa DECIMAL(7,2),
    alert_level VARCHAR(20) NOT NULL DEFAULT 'normal',
    location VARCHAR(200),
    processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    raw_data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_sensor_data_device_id ON sensor_data(device_id);
CREATE INDEX IF NOT EXISTS idx_sensor_data_timestamp ON sensor_data(timestamp);
CREATE INDEX IF NOT EXISTS idx_sensor_data_alert_level ON sensor_data(alert_level);
CREATE INDEX IF NOT EXISTS idx_sensor_data_processed_at ON sensor_data(processed_at);

-- 创建设备管理表
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(100) UNIQUE NOT NULL,
    device_name VARCHAR(200),
    device_type VARCHAR(50) NOT NULL DEFAULT 'sensor',
    location VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_seen TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建告警日志表
CREATE TABLE IF NOT EXISTS alert_logs (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    alert_level VARCHAR(20) NOT NULL,
    alert_message TEXT,
    sensor_data_id INTEGER REFERENCES sensor_data(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_alert_logs_device_id ON alert_logs(device_id);
CREATE INDEX IF NOT EXISTS idx_alert_logs_alert_level ON alert_logs(alert_level);
CREATE INDEX IF NOT EXISTS idx_alert_logs_created_at ON alert_logs(created_at);

-- 创建数据处理统计表
CREATE TABLE IF NOT EXISTS processing_stats (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    total_messages INTEGER NOT NULL DEFAULT 0,
    successful_messages INTEGER NOT NULL DEFAULT 0,
    failed_messages INTEGER NOT NULL DEFAULT 0,
    warning_alerts INTEGER NOT NULL DEFAULT 0,
    critical_alerts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date)
);

-- 插入一些示例设备数据
INSERT INTO devices (device_id, device_name, device_type, location, status) VALUES
('sensor_001', '温湿度传感器01', 'sensor', '机房A-机柜1', 'active'),
('sensor_002', '温湿度传感器02', 'sensor', '机房A-机柜2', 'active'),
('sensor_003', '气压传感器01', 'sensor', '机房B-机柜1', 'active')
ON CONFLICT (device_id) DO NOTHING;

-- 创建视图：设备最新数据
CREATE OR REPLACE VIEW device_latest_data AS
SELECT 
    d.device_id,
    d.device_name,
    d.device_type,
    d.location,
    d.status,
    sd.temperature_celsius,
    sd.temperature_fahrenheit,
    sd.humidity_percentage,
    sd.pressure_hpa,
    sd.alert_level,
    sd.timestamp,
    sd.processed_at
FROM devices d
LEFT JOIN LATERAL (
    SELECT * FROM sensor_data 
    WHERE device_id = d.device_id 
    ORDER BY timestamp DESC 
    LIMIT 1
) sd ON true;

-- 创建函数：更新设备最后见到时间
CREATE OR REPLACE FUNCTION update_device_last_seen()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE devices 
    SET last_seen = NEW.timestamp, updated_at = CURRENT_TIMESTAMP
    WHERE device_id = NEW.device_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
DROP TRIGGER IF EXISTS trigger_update_device_last_seen ON sensor_data;
CREATE TRIGGER trigger_update_device_last_seen
    AFTER INSERT ON sensor_data
    FOR EACH ROW
    EXECUTE FUNCTION update_device_last_seen();

-- 创建函数：自动记录告警日志
CREATE OR REPLACE FUNCTION log_alerts()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.alert_level IN ('warning', 'critical') THEN
        INSERT INTO alert_logs (device_id, alert_level, alert_message, sensor_data_id)
        VALUES (
            NEW.device_id, 
            NEW.alert_level, 
            CASE 
                WHEN NEW.alert_level = 'warning' THEN '设备检测到警告级别异常'
                WHEN NEW.alert_level = 'critical' THEN '设备检测到严重级别异常'
            END,
            NEW.id
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
DROP TRIGGER IF EXISTS trigger_log_alerts ON sensor_data;
CREATE TRIGGER trigger_log_alerts
    AFTER INSERT ON sensor_data
    FOR EACH ROW
    EXECUTE FUNCTION log_alerts();

-- 创建函数：更新处理统计
CREATE OR REPLACE FUNCTION update_processing_stats()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO processing_stats (date, total_messages, successful_messages, warning_alerts, critical_alerts)
    VALUES (
        CURRENT_DATE,
        1,
        1,
        CASE WHEN NEW.alert_level = 'warning' THEN 1 ELSE 0 END,
        CASE WHEN NEW.alert_level = 'critical' THEN 1 ELSE 0 END
    )
    ON CONFLICT (date) DO UPDATE SET
        total_messages = processing_stats.total_messages + 1,
        successful_messages = processing_stats.successful_messages + 1,
        warning_alerts = processing_stats.warning_alerts + CASE WHEN NEW.alert_level = 'warning' THEN 1 ELSE 0 END,
        critical_alerts = processing_stats.critical_alerts + CASE WHEN NEW.alert_level = 'critical' THEN 1 ELSE 0 END,
        updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
DROP TRIGGER IF EXISTS trigger_update_processing_stats ON sensor_data;
CREATE TRIGGER trigger_update_processing_stats
    AFTER INSERT ON sensor_data
    FOR EACH ROW
    EXECUTE FUNCTION update_processing_stats();

-- 创建查询今日统计的视图
CREATE OR REPLACE VIEW today_stats AS
SELECT 
    date,
    total_messages,
    successful_messages,
    failed_messages,
    warning_alerts,
    critical_alerts,
    ROUND((successful_messages::DECIMAL / NULLIF(total_messages, 0)) * 100, 2) as success_rate
FROM processing_stats 
WHERE date = CURRENT_DATE;

-- 权限设置
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO postgres; 