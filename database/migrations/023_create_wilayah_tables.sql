-- ============================================================
-- Migration 023: Tabel Wilayah Indonesia
-- Provinsi, Kota/Kabupaten, Kecamatan, Kelurahan/Desa, Kode Pos
-- ============================================================

-- 1. Tabel Provinsi
CREATE TABLE IF NOT EXISTS provinces (
    id      INT PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,

    INDEX idx_provinces_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Tabel Kota/Kabupaten
CREATE TABLE IF NOT EXISTS cities (
    id          INT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    province_id INT NOT NULL,

    INDEX idx_cities_name (name),
    INDEX idx_cities_province_id (province_id),

    CONSTRAINT fk_cities_province FOREIGN KEY (province_id)
        REFERENCES provinces(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Tabel Kecamatan
CREATE TABLE IF NOT EXISTS districts (
    id      INT PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    city_id INT NOT NULL,

    INDEX idx_districts_name (name),
    INDEX idx_districts_city_id (city_id),

    CONSTRAINT fk_districts_city FOREIGN KEY (city_id)
        REFERENCES cities(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Tabel Kelurahan/Desa
CREATE TABLE IF NOT EXISTS subdistricts (
    id          INT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    district_id INT NOT NULL,

    INDEX idx_subdistricts_name (name),
    INDEX idx_subdistricts_district_id (district_id),

    CONSTRAINT fk_subdistricts_district FOREIGN KEY (district_id)
        REFERENCES districts(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. Tabel Kode Pos
CREATE TABLE IF NOT EXISTS postal_codes (
    id              INT PRIMARY KEY,
    subdistrict_id  INT NOT NULL,
    district_id     INT NOT NULL,
    city_id         INT NOT NULL,
    province_id     INT NOT NULL,
    postal_code     VARCHAR(10) NOT NULL,

    INDEX idx_postal_codes_code (postal_code),
    INDEX idx_postal_codes_subdistrict_id (subdistrict_id),
    INDEX idx_postal_codes_district_id (district_id),
    INDEX idx_postal_codes_city_id (city_id),
    INDEX idx_postal_codes_province_id (province_id),

    CONSTRAINT fk_postal_codes_subdistrict FOREIGN KEY (subdistrict_id)
        REFERENCES subdistricts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_postal_codes_district FOREIGN KEY (district_id)
        REFERENCES districts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_postal_codes_city FOREIGN KEY (city_id)
        REFERENCES cities(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_postal_codes_province FOREIGN KEY (province_id)
        REFERENCES provinces(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
