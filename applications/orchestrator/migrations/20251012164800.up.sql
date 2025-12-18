-- migrations/001_initial_schema.up.sql
-- Create users table
CREATE TABLE IF NOT EXISTS users
(
    id             INT auto_increment PRIMARY KEY,
    identification VARCHAR (255) NOT NULL UNIQUE,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE
    CURRENT_TIMESTAMP,
    INDEX idx_identification ( identification )
    )
    engine = innodb
    DEFAULT charset = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

-- Create conversations table
CREATE TABLE IF NOT EXISTS conversations
(
    id         INT auto_increment PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
    engine = innodb
    DEFAULT charset = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

-- Create user_conversations table (junction table)
CREATE TABLE IF NOT EXISTS user_conversations
(
    id                  INT auto_increment PRIMARY KEY,
    user_identification VARCHAR (255) NOT NULL,
    conversation_id     INT NOT NULL,
    role                VARCHAR (50) NOT NULL,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE
    CURRENT_TIMESTAMP,
    FOREIGN KEY ( user_identification ) REFERENCES users ( identification ) ON
                                                                DELETE CASCADE,
    FOREIGN KEY ( conversation_id ) REFERENCES conversations ( id ) ON DELETE
    CASCADE,
    INDEX idx_user_identification ( user_identification ),
    INDEX idx_conversation_id ( conversation_id ),
    UNIQUE KEY unique_user_conversation ( user_identification, conversation_id
                                        )
    )
    engine = innodb
    DEFAULT charset = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

-- Create messages table
CREATE TABLE IF NOT EXISTS messages
(
    id              INT auto_increment PRIMARY KEY,
    sender          VARCHAR (255) NOT NULL,
    conversation_id BIGINT NOT NULL,
    content         TEXT NOT NULL,
    type            VARCHAR (50) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE
    CURRENT_TIMESTAMP,
    INDEX idx_conversation_id ( conversation_id ),
    INDEX idx_sender ( sender ),
    INDEX idx_created_at ( created_at )
    )
    engine = innodb
    DEFAULT charset = utf8mb4
    COLLATE = utf8mb4_unicode_ci;