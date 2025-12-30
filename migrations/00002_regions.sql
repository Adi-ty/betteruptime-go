-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "region" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    CONSTRAINT "region_pkey" PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "region";
-- +goose StatementEnd
