-- migrations/001_initial_schema.up.sql

-- Create users table
CREATE TABLE IF NOT EXISTS users
(
    id
               INT
        AUTO_INCREMENT
        PRIMARY
            KEY,
    identification
               VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_identification
        (
         identification
            )
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- Create conversations table
CREATE TABLE IF NOT EXISTS conversations
(
    id
        INT
        AUTO_INCREMENT
        PRIMARY
            KEY,
    created_at
        TIMESTAMP
        DEFAULT
                    CURRENT_TIMESTAMP,
    updated_at
        TIMESTAMP
        DEFAULT
                    CURRENT_TIMESTAMP
        ON
            UPDATE
                    CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- Create user_conversations table (junction table)
CREATE TABLE IF NOT EXISTS user_conversations
(
    id
                    INT
        AUTO_INCREMENT
        PRIMARY
            KEY,
    user_identification
                    VARCHAR(255) NOT NULL,
    conversation_id INT          NOT NULL,
    role            VARCHAR(50)  NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY
        (
         user_identification
            ) REFERENCES users
        (
         identification
            )
        ON DELETE CASCADE,
    FOREIGN KEY
        (
         conversation_id
            ) REFERENCES conversations
        (
         id
            )
        ON DELETE CASCADE,
    INDEX idx_user_identification
        (
         user_identification
            ),
    INDEX idx_conversation_id
        (
         conversation_id
            ),
    UNIQUE KEY unique_user_conversation
        (
         user_identification,
         conversation_id
            )
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- Create messages table
CREATE TABLE IF NOT EXISTS messages
(
    id
                    INT
        AUTO_INCREMENT
        PRIMARY
            KEY,
    sender
                    VARCHAR(255) NOT NULL,
    conversation_id BIGINT       NOT NULL,
    content         TEXT         NOT NULL,
    type            VARCHAR(50)  NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_conversation_id
        (
         conversation_id
            ),
    INDEX idx_sender
        (
         sender
            ),
    INDEX idx_created_at
        (
         created_at
            )
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- ============================================
-- SEED DATA
-- ============================================

-- Insert users
INSERT INTO users (identification)
VALUES ('user_alice'),
       ('user_bob'),
       ('user_charlie'),
       ('user_diana');

-- Insert conversations
INSERT INTO conversations (id)
VALUES (1),
       (2),
       (3);

-- Insert user_conversations (map users to conversations with roles)
INSERT INTO user_conversations (user_identification, conversation_id, role)
VALUES ('user_alice', 1, 'admin'),
       ('user_bob', 1, 'member'),
       ('user_alice', 2, 'member'),
       ('user_charlie', 2, 'admin'),
       ('user_bob', 3, 'admin'),
       ('user_diana', 3, 'member');

-- Insert messages
INSERT INTO messages (sender, conversation_id, content, type)
VALUES
    -- Conversation 1: Alice and Bob
    ('user_alice', 1, 'Hey Bob, how are you?', 'text'),
    ('user_bob', 1, 'Hi Alice! I am doing great, thanks!', 'text'),
    ('user_alice', 1, 'Did you finish the project?', 'text'),
    ('user_bob', 1, 'Yes, just submitted it this morning.', 'text'),

    -- Conversation 2: Alice and Charlie
    ('user_alice', 2, 'Charlie, can we schedule a meeting?', 'text'),
    ('user_charlie', 2, 'Sure! How about tomorrow at 3 PM?', 'text'),
    ('user_alice', 2, 'Perfect, see you then!', 'text'),
    ('user_charlie', 2, 'https://meet.example.com/alice-charlie', 'link'),

    -- Conversation 3: Bob and Diana
    ('user_bob', 3, 'Diana, welcome to the team!', 'text'),
    ('user_diana', 3, 'Thanks Bob! Excited to be here.', 'text'),
    ('user_bob', 3, 'Let me know if you need any help.', 'text'),
    ('user_diana', 3, 'Will do, appreciate it!', 'text'),
    ('user_bob', 3, 'Here is the onboarding doc', 'text'),
    ('user_diana', 3, 'Got it, reviewing now.', 'text');