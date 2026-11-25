SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mother (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    first_name VARCHAR(50) NOT NULL, 
    last_name VARCHAR(50) NOT NULL, 
    phone VARCHAR(20) NOT NULL DEFAULT '', 
    email VARCHAR(255) NOT NULL DEFAULT '', 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS help_items (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    mother_id INT NOT NULL, 
    kind VARCHAR(50) NOT NULL, 
    sum INT NOT NULL,
    check_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS child (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    mother_id INT NOT NULL, 
    first_name VARCHAR(50) NOT NULL,
    birth_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mother_secret_info (
    mother_id INT PRIMARY KEY,
    criminal_record VARCHAR(100) NOT NULL DEFAULT '',
    nd_pnd VARCHAR(255) NOT NULL DEFAULT '',
    drugs VARCHAR(255) NOT NULL DEFAULT '',
    chronic_diseases VARCHAR(255) NOT NULL DEFAULT '',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS watching_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    action VARCHAR(50) NOT NULL,
    mother_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS event_type (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO event_type (name) VALUES
    ('День матери'), ('Обучение кулинарии'),
    ('Встреча с батюшкой'), ('Встреча с психологом');

CREATE TABLE IF NOT EXISTS event (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    starts_at TIMESTAMP NOT NULL,
    place VARCHAR(255) NOT NULL DEFAULT '',
    max_people INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (type_id) REFERENCES event_type(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS participation (
    event_id INT NOT NULL,
    mother_id INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'registered',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, mother_id),
    FOREIGN KEY (event_id) REFERENCES event(id),
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mother_profile (
    mother_id INT PRIMARY KEY,
    birth_date DATE NULL,
    birth_place VARCHAR(150) NOT NULL DEFAULT '',
    citizenship VARCHAR(100) NOT NULL DEFAULT '',
    registration VARCHAR(255) NOT NULL DEFAULT '',
    living VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '',
    religion VARCHAR(100) NOT NULL DEFAULT '',
    education VARCHAR(100) NOT NULL DEFAULT '',
    profession VARCHAR(150) NOT NULL DEFAULT '',
    income_source VARCHAR(100) NOT NULL DEFAULT '',
    income_amount INT NOT NULL DEFAULT 0,
    marital_status VARCHAR(100) NOT NULL DEFAULT '',
    father_info VARCHAR(255) NOT NULL DEFAULT '',
    pregnancy_weeks INT NOT NULL DEFAULT 0,
    situation TEXT NOT NULL,
    expected_help TEXT NOT NULL,
    source VARCHAR(150) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mother_status_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mother_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    reason VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS help_request (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mother_id INT NOT NULL,
    kind VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'new',
    staff_comment VARCHAR(500) NOT NULL DEFAULT '',
    source VARCHAR(20) NOT NULL DEFAULT 'staff',
    help_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mother_id) REFERENCES mother(id),
    FOREIGN KEY (help_id) REFERENCES help_items(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS document (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mother_id INT NOT NULL,
    doc_type VARCHAR(30) NOT NULL,
    help_id INT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'not_signed',
    signed_at DATE NULL,
    file_name VARCHAR(100) NOT NULL DEFAULT '',
    original_name VARCHAR(255) NOT NULL DEFAULT '',
    content_type VARCHAR(100) NOT NULL DEFAULT '',
    file_size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY one_statement_per_help (help_id),
    FOREIGN KEY (mother_id) REFERENCES mother(id),
    FOREIGN KEY (help_id) REFERENCES help_items(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS donor (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    kind VARCHAR(20) NOT NULL DEFAULT 'person',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '',
    comment VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS donation (
    id INT AUTO_INCREMENT PRIMARY KEY,
    donor_id INT NOT NULL,
    amount INT NOT NULL,
    donated_at DATE NOT NULL,
    comment VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donor(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mother_learn_link (
    mother_id INT NOT NULL,
    learn_user_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (mother_id),
    FOREIGN KEY (mother_id) REFERENCES mother(id),
    FOREIGN KEY (learn_user_id) REFERENCES lms_learners(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;