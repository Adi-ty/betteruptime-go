-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "website" (
    "id" SERIAL PRIMARY KEY,
    "url" TEXT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "time_added" TIMESTAMP(3) NOT NULL
);

ALTER TABLE "website" ADD CONSTRAINT "website_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE UNIQUE INDEX "website_url_key" ON "website"("url");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "website";
-- +goose StatementEnd