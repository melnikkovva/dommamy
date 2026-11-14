INSERT IGNORE INTO app_user (login, password, first_name, last_name, role, created_at)
VALUES ('admin', '$2a$10$TC2ciqTNgu47Sm8t5EHCIeyPYJtCcQ26lGeSUPtJc7TGK/PZDEkCm', 'Admin', 'Admin', 'admin', NOW());
