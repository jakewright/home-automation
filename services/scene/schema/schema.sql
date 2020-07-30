USE home_automation;

CREATE TABLE IF NOT EXISTS service_scene_scenes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    owner_id INT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS service_scene_actions (
    scene_id INT NOT NULL,
    stage INT NOT NULL, -- Ordering within the scene
    sequence INT NOT NULL, -- Ordering within the stage

    func VARCHAR(64), -- e.g. sleep()
    controller_name VARCHAR(64),
    device_id VARCHAR(64),
    command VARCHAR(64),
    property VARCHAR(64),
    property_value VARCHAR(64),
    property_type VARCHAR(64),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),

    PRIMARY KEY (scene_id, stage, sequence),

    FOREIGN KEY (scene_id) REFERENCES service_scene_scenes(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);
