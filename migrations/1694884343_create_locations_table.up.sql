CREATE TABLE IF NOT EXISTS locations (
    id INT NOT NULL AUTO_INCREMENT,
    ip_address VARCHAR(255) NOT NULL,
    country_code VARCHAR(255) NOT NULL,
    country  VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    mystery_value INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uc_location UNIQUE (ip_address,country_code,country,city,latitude,longitude,mystery_value),
    PRIMARY KEY(id)
)
CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;