USE home_automation;

CREATE TABLE IF NOT EXISTS service_user_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW()
);
