-- Table: "Users"

-- DROP TABLE "Users";

CREATE TABLE "Users"
(
  "ID" bigserial NOT NULL,
  "EmailAddress" text NOT NULL,
  "Name" text NOT NULL,
  "AuthenticationID" text NOT NULL,
  "DateCreated" timestamp without time zone NOT NULL DEFAULT now(),
  CONSTRAINT "PK_Users" PRIMARY KEY ("ID")
)
WITH (
  OIDS=FALSE
);
ALTER TABLE "Users"
  OWNER TO godutch;

ALTER TABLE "Users" ADD CONSTRAINT EmailAddressIsUnique UNIQUE ("EmailAddress");

ALTER TABLE "Users" ADD CONSTRAINT AuthenticationIDIsUnique UNIQUE ("AuthenticationID");

