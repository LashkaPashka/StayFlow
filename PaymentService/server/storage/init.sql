CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    method VARCHAR(50) NOT NULL,
    payment_id VARCHAR(255),
    client_secret VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    response TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);