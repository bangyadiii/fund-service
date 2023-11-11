-- buatkan kode saya table medias yang dapat menyimpan data dengan relational polimorphic, dan ada field untuk menyimpan data path, dan collection name dan juga conversion name

CREATE TABLE medias (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    mediable_id BIGINT NOT NULL,
    mediable_type VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    collection_name VARCHAR(255) NOT NULL,
    conversion JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);

CREATE INDEX medias_mediable_index ON medias (mediable_type, mediable_id);
CREATE INDEX medias_mediable_collections_index ON medias (mediable_type, mediable_id, collection_name);
