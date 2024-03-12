CREATE TABLE carousel (
    id SERIAL PRIMARY KEY,
    file_name varchar(255) NOT NULL,
    "order" int,
    file_path varchar(255) NOT NULL,
    file_size int NOT NULL,
    file_type varchar(50) NOT NULL,
    uploaded_at timestamp(0) NOT NULL
);