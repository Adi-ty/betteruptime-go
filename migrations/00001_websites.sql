-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "website" (
    "id" SERIAL PRIMARY KEY,
    "url" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "website_url_key" ON "website"("url");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "website";
-- +goose StatementEnd