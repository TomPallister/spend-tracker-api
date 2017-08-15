-- Table: public."Trackers"

-- DROP TABLE public."Trackers";

CREATE TABLE public."Trackers"

(

  "ID" bigserial NOT NULL,

  "AdminUserID" bigint NOT NULL,

  "Name" text NOT NULL,

  "Currency" text NOT NULL,

  "DateCreated" timestamp without time zone NOT NULL DEFAULT now(),

  CONSTRAINT "PK_Trackers" PRIMARY KEY ("ID"),

  CONSTRAINT "FK_Trackers_Users_AdminUserId" FOREIGN KEY ("AdminUserID")

      REFERENCES public."Users" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION

)

WITH (

  OIDS=FALSE

);

ALTER TABLE public."Trackers"

  OWNER TO godutch;

