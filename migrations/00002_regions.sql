-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Region" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Region_pkey" PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "Region";
-- +goose StatementEnd
