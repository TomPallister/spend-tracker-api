-- Table: public."Spends"

-- DROP TABLE public."Spends";

CREATE TABLE public."Spends"

(

  "ID" bigserial NOT NULL,

  "TrackerID" bigint NOT NULL,

  "UserID" bigint NOT NULL,

  "Name" text NOT NULL,
  
  "DateCreated" timestamp without time zone NOT NULL DEFAULT now(),

  "Value" numeric NOT NULL,

  "Currency" text NOT NULL DEFAULT '£'::text,

  CONSTRAINT "TrackerSpends_pkey" PRIMARY KEY ("ID"),

  CONSTRAINT "FK_Tracks_Trackers_TrackerID" FOREIGN KEY ("TrackerID")

      REFERENCES public."Trackers" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION,

  CONSTRAINT "FK_Tracks_Users_UserID" FOREIGN KEY ("UserID")

      REFERENCES public."Users" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION

)

WITH (

  OIDS=FALSE

);

ALTER TABLE public."Spends"

  OWNER TO godutch;



  

  