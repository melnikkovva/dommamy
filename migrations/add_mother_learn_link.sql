-- Перенос связи mother ↔ learn в отдельную таблицу, принадлежащую контексту learn.
-- learn_user_id удалён из mother; mother читает связь через JOIN с mother_learn_link.

CREATE TABLE IF NOT EXISTS mother_learn_link (
    mother_id INT NOT NULL,
    learn_user_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (mother_id),
    FOREIGN KEY (mother_id) REFERENCES mother(id),
    FOREIGN KEY (learn_user_id) REFERENCES lms_learners(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Миграция данных из старого столбца в новую таблицу.
INSERT INTO mother_learn_link (mother_id, learn_user_id)
    SELECT id, learn_user_id FROM mother WHERE learn_user_id IS NOT NULL;

-- Удаление столбца learn_user_id (Foreign Key удаляется автоматически вместе со столбцом).
ALTER TABLE mother DROP COLUMN learn_user_id;
