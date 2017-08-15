-- Table: public."SpendSummaries"

-- DROP TABLE public."SpendSummaries";

CREATE TABLE public."SpendSummaries"

(

  "ID" bigserial NOT NULL,

  "TrackerID" bigint NOT NULL,

  "UserID" bigint NOT NULL,

  "Currency" text NOT NULL DEFAULT '£'::text,

  "Value" numeric NOT NULL,

  CONSTRAINT "PK_SpendSummaries" PRIMARY KEY ("ID"),

  CONSTRAINT "FK_SpendSummaries_Trackers_TrackerID" FOREIGN KEY ("TrackerID")

      REFERENCES public."Trackers" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION,

  CONSTRAINT "FK_SpendSummaries_Users_UserID" FOREIGN KEY ("UserID")

      REFERENCES public."Users" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION
)

WITH (

  OIDS=FALSE

);

ALTER TABLE public."SpendSummaries"

  OWNER TO godutch;

-- Index: public."NonClusteredIndex-20151207-211821"

-- DROP INDEX public."NonClusteredIndex-20151207-211821";

CREATE INDEX "NonClusteredIndex-20151207-211821"

  ON public."SpendSummaries"

  USING btree

  ("TrackerID");



