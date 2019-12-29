CREATE TABLE IF NOT EXISTS service_scene_scenes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS service_scene_actions (
    id INT AUTO_INCREMENT PRIMARY_KEY,
    scene_id INT NOT NULL,
    stage INT NOT NULL, -- Ordering within the scene
    index INT NOT NULL, -- Ordering within the stage

    function VARCHAR(64), -- e.g. sleep()
    controller_name VARCHAR(64),
    command VARCHAR(64),
    property VARCHAR(64),
    property_value VARCHAR(64),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
    deleted_at TIMESTAMP
);
