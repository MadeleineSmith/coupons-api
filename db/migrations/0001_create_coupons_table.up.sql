CREATE TABLE IF NOT EXISTS coupons (
  id uuid DEFAULT uuid_generate_v1mc() PRIMARY KEY,
  name VARCHAR NOT NULL,
  brand VARCHAR NOT NULL,
  value INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
  expiry TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP + (1 * interval '1 year')
)