CREATE TABLE IF NOT EXISTS vehicle (
  id bigserial PRIMARY KEY,
  type varchar(10) NOT NULL,
  plate_number char(7) NOT NULL
);