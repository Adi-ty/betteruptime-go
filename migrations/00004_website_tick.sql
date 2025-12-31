-- +goose Up
-- +goose StatementBegin
CREATE TYPE "website_status" AS ENUM ('UP', 'DOWN', 'UNKNOWN');

CREATE TABLE IF NOT EXISTS "website_tick" (
    "id" TEXT NOT NULL,
    "response_time_ms" INTEGER NOT NULL,
    "status_code" "website_status" NOT NULL,
    "website_id" INTEGER NOT NULL,
    "region_id" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "website_tick_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "website_tick" ADD CONSTRAINT "website_tick_website_id_fkey" FOREIGN KEY ("website_id") REFERENCES "website"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "website_tick" ADD CONSTRAINT "website_tick_region_id_fkey" FOREIGN KEY ("region_id") REFERENCES "region"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "website_tick";
DROP TYPE "website_status";
-- +goose StatementEnd