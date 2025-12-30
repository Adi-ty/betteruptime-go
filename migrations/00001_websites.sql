-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "Website" (
    "id" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Website_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "Website_url_key" ON "Website"("url");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "Website";
-- +goose StatementEnd
