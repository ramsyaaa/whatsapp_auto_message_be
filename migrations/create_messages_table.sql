-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    recipient_number VARCHAR(20) NOT NULL,
    message_content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Sent',
    sent_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
