CREATE TABLE carousel (
    id SERIAL PRIMARY KEY,
    "order" int,
    file_path varchar(255) NOT NULL,
    file_size int NOT NULL,
    file_type varchar(50) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL
);