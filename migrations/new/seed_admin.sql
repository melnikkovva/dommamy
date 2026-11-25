INSERT IGNORE INTO app_user (login, password, first_name, last_name, role, created_at)
VALUES ('admin', '$2a$10$XFeNLjx/AHPnDjZ2.aYuje3uT5iBJVm9wvs7UxWaK3XG1pa4YNUxS', 'Admin', 'Admin', 'admin', NOW());
