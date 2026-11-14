ALTER TABLE mother ADD COLUMN IF NOT EXISTS learn_user_id VARCHAR(100) NULL;

CREATE TABLE IF NOT EXISTS lms_learners (
    user_id VARCHAR(50) PRIMARY KEY,
    department_id VARCHAR(50),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    login VARCHAR(100),
    phone VARCHAR(100),
    status INT,
    added_date DATETIME,
    last_login_date DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE lms_learners ADD COLUMN IF NOT EXISTS phone VARCHAR(100) NULL;

CREATE TABLE IF NOT EXISTS lms_course_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    course_id VARCHAR(50) NOT NULL,
    course_title VARCHAR(255),
    enrollment_id VARCHAR(50),
    access_date DATETIME,
    due_date DATETIME NULL,
    completion_date DATETIME NULL,
    time_spent INT DEFAULT 0,
    duration INT DEFAULT 0,
    completion_status VARCHAR(50),
    progress DECIMAL(5,2) DEFAULT 0,
    awarded_score DECIMAL(5,2) NULL,
    is_overdue BOOLEAN DEFAULT FALSE,
    views_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_course (user_id, course_id),
    FOREIGN KEY (user_id) REFERENCES lms_learners(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
