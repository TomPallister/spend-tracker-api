-- Table: public."Transfers"

-- DROP TABLE public."Transfers";

CREATE TABLE public."Transfers"

(

  "ID" bigserial NOT NULL,

  "ToUserID" bigint NOT NULL,

  "FromUserID" bigint NOT NULL,

  "TrackerID" bigint NOT NULL,

  "Currency" text NOT NULL DEFAULT '£'::text,

  "Value" numeric NOT NULL,

  CONSTRAINT "PK_Transfers" PRIMARY KEY ("ID"),

  CONSTRAINT "FK_Transfers_Trackers_TrackerID" FOREIGN KEY ("TrackerID")

      REFERENCES public."Trackers" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION,

  CONSTRAINT "FK_Transfers_Users_ToUserID" FOREIGN KEY ("ToUserID")

      REFERENCES public."Users" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION,
      
  CONSTRAINT "FK_Transfers_Users_FromUserID" FOREIGN KEY ("FromUserID")

      REFERENCES public."Users" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION
)

WITH (

  OIDS=FALSE

);

ALTER TABLE public."Transfers"

  OWNER TO godutch;



-- Index: public."NonClusteredIndex-20151207-211821"



-- DROP INDEX public."NonClusteredIndex-20151207-211821";



CREATE INDEX "NonClusteredIndex-Transfers-TrackerID"

  ON public."Transfers"

  USING btree

  ("TrackerID");



