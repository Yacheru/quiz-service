CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    authorized BOOLEAN NOT NULL,
    authorized_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    quit_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

-- Варианты
CREATE TABLE IF NOT EXISTS variants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) UNIQUE
);
--

-- Использование Many-To-Many позволит нам по желанию менять количество ответов не изменяя структуру таблицы
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    variant_id INTEGER NOT NULL,
    question VARCHAR(50) NOT NULL,
    answer VARCHAR(50) NOT NULL,
    FOREIGN KEY (variant_id) REFERENCES variants(id) ON DELETE CASCADE,
    UNIQUE (variant_id, question)
);

CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    answer VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS questions_and_answers (
    questions_id INTEGER NOT NULL,
    answers_id INTEGER NOT NULL,
    FOREIGN KEY (questions_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (answers_id) REFERENCES answers(id) ON DELETE CASCADE
);
--

-- Начало тестирования
CREATE TABLE IF NOT EXISTS testing (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    variant_id INTEGER NOT NULL,
    start_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    correct_answers INTEGER DEFAULT 0,
    finish_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES auth(id),
    FOREIGN KEY (variant_id) REFERENCES variants(id),
    UNIQUE (user_id, variant_id)
);
--

-- Список ответов пользователей
CREATE TABLE IF NOT EXISTS user_answers (
    id SERIAL PRIMARY KEY,
    test_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answer VARCHAR(50) NOT NULL,
    FOREIGN KEY (test_id) REFERENCES testing(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);
--
