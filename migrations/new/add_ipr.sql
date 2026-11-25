CREATE TABLE IF NOT EXISTS lms_ipr_plans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    plan_id VARCHAR(100) NOT NULL,
    title VARCHAR(500),
    employee_id VARCHAR(100) NOT NULL,
    mentor_id VARCHAR(100),
    mentor_name VARCHAR(255),
    start_date DATETIME,
    end_date DATETIME NULL,
    status VARCHAR(50),
    progress DECIMAL(5,2) DEFAULT 0,
    completed_tasks INT DEFAULT 0,
    total_tasks INT DEFAULT 0,
    goals JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_plan (plan_id),
    INDEX idx_employee (employee_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lms_ipr_stages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    plan_id VARCHAR(100) NOT NULL,
    stage_id VARCHAR(100) NOT NULL,
    title VARCHAR(500),
    description TEXT,
    start_date DATETIME NULL,
    end_date DATETIME NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_plan (plan_id),
    FOREIGN KEY (plan_id) REFERENCES lms_ipr_plans(plan_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lms_ipr_tasks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    stage_id BIGINT NOT NULL,
    task_id VARCHAR(100) NOT NULL,
    task_type VARCHAR(50),
    title VARCHAR(500),
    description TEXT,
    status VARCHAR(50),
    is_overdue BOOLEAN DEFAULT FALSE,
    absolute_due_date DATETIME NULL,
    completion_date DATETIME NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_stage (stage_id),
    FOREIGN KEY (stage_id) REFERENCES lms_ipr_stages(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
