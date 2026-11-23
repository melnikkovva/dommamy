CREATE TABLE IF NOT EXISTS consultation (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mother_id INT NOT NULL,
    specialist_type VARCHAR(50) NOT NULL,
    consultation_date DATE NOT NULL,
    notes TEXT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
