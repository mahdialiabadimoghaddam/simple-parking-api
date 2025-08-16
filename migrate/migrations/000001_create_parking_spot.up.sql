CREATE TABLE IF NOT EXISTS parking_spot(
  id bigserial PRIMARY KEY,
  row_number int NOT NULL,
  column_number int NOT NULL,
  vehicle_type char(10) NOT NULL,
  empty boolean NOT NULL
);