CREATE TABLE media (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    model_type        VARCHAR(100) NOT NULL,
    model_id          CHAR(36)     NOT NULL,
    collection_name   VARCHAR(100) NOT NULL,
    disk              VARCHAR(50)  NOT NULL,
    directory         VARCHAR(255) NOT NULL,
    file_name         VARCHAR(255) NOT NULL,
    name              VARCHAR(255) NOT NULL,
    mime_type         VARCHAR(150) NOT NULL,
    size              BIGINT       NOT NULL DEFAULT 0,
    conversions       JSON         NULL,
    custom_properties JSON         NULL,
    order_column      INT          NOT NULL DEFAULT 0,
    created_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_media_model (model_type, model_id),
    INDEX idx_media_collection (model_type, model_id, collection_name, order_column)
);