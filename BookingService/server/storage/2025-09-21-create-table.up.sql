CREATE TABLE bookings (
  booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  hotel_id UUID NOT NULL,
  room_type_id UUID NOT NULL,
  check_in DATE NOT NULL,
  check_out DATE NOT NULL,
  nights INT NOT NULL,
  rooms_count INT NOT NULL
  total_amount NUMERIC(12,2) NOT NULL,
  currency VARCHAR(3) DEFAULT 'USD',
  status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING','CONFIRMED','CANCELLED','FAILED')),
  idempotency_key TEXT, 
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE bookings
ADD COLUMN rooms_count INT NOT NULL DEFAULT 1;