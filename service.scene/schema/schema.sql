CREATE TABLE IF NOT EXISTS service_scene_scenes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS service_scene_actions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    scene_id INT NOT NULL,
    stage INT NOT NULL, -- Ordering within the scene
    sequence INT NOT NULL, -- Ordering within the stage

    func VARCHAR(64), -- e.g. sleep()
    controller_name VARCHAR(64),
    command VARCHAR(64),
    property VARCHAR(64),
    property_value VARCHAR(64),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
    deleted_at TIMESTAMP,

    FOREIGN KEY (scene_id) REFERENCES service_scene_scenes(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);
