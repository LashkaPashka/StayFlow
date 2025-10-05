CREATE TABLE hotels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  city VARCHAR(100) NOT NULL,
  address VARCHAR(255)
);

CREATE TABLE room_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hotel_id UUID NOT NULL REFERENCES hotels(id),
  type VARCHAR(50) NOT NULL,
  price NUMERIC(10,2) NOT NULL,
  total_count INT NOT NULL
);

CREATE TABLE room_inventory (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hotel_id UUID NOT NULL REFERENCES hotels(id),
  room_type_id UUID NOT NULL REFERENCES room_types(id),
  date DATE NOT NULL,
  available INT NOT NULL,
  reserved INT NOT NULL DEFAULT 0,
  UNIQUE (room_type_id, date)
);
