USE home_automation;

INSERT INTO `service_scene_scenes` (`id`, `name`, `owner_id`) VALUES
(1, 'Hue light test', 1);

INSERT INTO `service_scene_actions` (`scene_id`, `stage`, `sequence`, `func`, `controller_name`, `device_id`, `command`, `property`, `property_value`, `property_type`) VALUES
(1, 1, 1, '', 'service.controller.hue', 'jake-desk-lamp', '', 'power', 'false', 'boolean'),
(1, 1, 2, 'sleep 2s', '', '', '', '', '', ''),
(1, 2, 1, '', 'service.controller.hue', 'jake-desk-lamp', '', 'power', 'true', 'boolean');
