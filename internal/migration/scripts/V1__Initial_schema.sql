-- Check and create schema if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT schema_name
        FROM information_schema.schemata
        WHERE schema_name = 'cve'
    ) THEN
        EXECUTE 'CREATE SCHEMA cve';
    END IF;
END $$;

-- Check and create table if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_schema = 'cve'
          AND table_name = 'cve_details'
    ) THEN
        EXECUTE 'CREATE TABLE cve.cve_details (
            id SERIAL PRIMARY KEY,
            cve_id VARCHAR(255) NOT NULL,
            cve_data JSONB NOT NULL,
            version INTEGER NOT NULL,
            cve_data_hash VARCHAR(64) NOT NULL
        )';
    END IF;
END $$;

-- Check and create indexes if they do not exist - using default btree index method
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'cve' AND tablename = 'cve_details' AND indexname = 'idx_cve_id'
    ) THEN
        EXECUTE 'CREATE INDEX idx_cve_id ON cve.cve_details USING btree (cve_id)';
    END IF;
    IF NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'cve' AND tablename = 'cve_details' AND indexname = 'idx_version'
    ) THEN
        EXECUTE 'CREATE INDEX idx_version ON cve.cve_details USING btree (version)';
    END IF;
    IF NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'cve' AND tablename = 'cve_details' AND indexname = 'cve_data_hash'
    ) THEN
        EXECUTE 'CREATE INDEX idx_cve_data_hash ON cve.cve_details USING btree (cve_data_hash)';
    END IF;
END $$;

