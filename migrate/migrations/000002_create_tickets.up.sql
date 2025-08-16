CREATE TABLE IF NOT EXISTS tickets (
  id bigserial PRIMARY KEY,
  vehicle_id int NOT NULL,
  parking_spot_id int NOT NULL,
  content varchar(255),
  entery_time bigint
);