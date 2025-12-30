-- +goose Up
-- +goose StatementBegin
CREATE TYPE "WebsiteStatus" AS ENUM ('UP', 'DOWN', 'UNKNOWN');

CREATE TABLE IF NOT EXISTS "WebsiteTick" (
    "id" TEXT NOT NULL,
    "responseTimeMs" INTEGER NOT NULL,
    "statusCode" "WebsiteStatus" NOT NULL,
    "websiteId" TEXT NOT NULL,
    "regionId" TEXT NOT NULL,

    CONSTRAINT "WebsiteTick_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "WebsiteTick" ADD CONSTRAINT "WebsiteTick_websiteId_fkey" FOREIGN KEY ("websiteId") REFERENCES "Website"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "WebsiteTick" ADD CONSTRAINT "WebsiteTick_regionId_fkey" FOREIGN KEY ("regionId") REFERENCES "Region"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "WebsiteTick";
DROP TYPE "WebsiteStatus";
-- +goose StatementEnd