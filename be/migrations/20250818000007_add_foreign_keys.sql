-- +goose Up
-- +goose StatementBegin

-- コアテーブルに外部キー制約を追加
DO $$
BEGIN IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_cores_maker_id' 
        AND table_name = 'cores'
        AND constraint_type = 'FOREIGN KEY'
    ) THEN
ALTER TABLE "cores"
    ADD CONSTRAINT "fk_cores_maker_id"
        FOREIGN KEY ("maker_id")
            REFERENCES "makers" ("id")
            ON UPDATE CASCADE
            ON DELETE RESTRICT;
END IF;
END$$;

-- カバーテーブルに外部キー制約を追加
DO $$
BEGIN IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_covers_maker_id' 
        AND table_name = 'covers'
        AND constraint_type = 'FOREIGN KEY'
    ) THEN
ALTER TABLE "covers"
    ADD CONSTRAINT "fk_covers_maker_id"
        FOREIGN KEY ("maker_id")
            REFERENCES "makers" ("id")
            ON UPDATE CASCADE
            ON DELETE RESTRICT;
END IF;
END$$;

-- インデックスを追加（パフォーマンス向上のため）
CREATE INDEX IF NOT EXISTS "idx_cores_maker_id" ON "cores" ("maker_id");
CREATE INDEX IF NOT EXISTS "idx_covers_maker_id" ON "covers" ("maker_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- インデックスを削除
DROP INDEX IF EXISTS "idx_covers_maker_id";
DROP INDEX IF EXISTS "idx_cores_maker_id";

-- 外部キー制約を削除
ALTER TABLE "covers"
    DROP CONSTRAINT IF EXISTS "fk_covers_maker_id";
ALTER TABLE "cores"
    DROP CONSTRAINT IF EXISTS "fk_cores_maker_id";

-- +goose StatementEnd