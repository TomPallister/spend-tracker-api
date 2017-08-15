-- Table: public."TrackerUsers"

-- DROP TABLE public."TrackerUsers";

CREATE TABLE public."TrackerUsers"

(

  "ID" bigserial NOT NULL,

  "TrackerID" bigint NOT NULL,

  "UserID" bigint NOT NULL,

  CONSTRAINT "PK_TrackerUsers" PRIMARY KEY ("ID", "TrackerID", "UserID"),

  CONSTRAINT "FK_TrackerUsers_Trackers_TrackerID" FOREIGN KEY ("TrackerID")

      REFERENCES public."Trackers" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION,

  CONSTRAINT "FK_TrackerUsers_Users_UserID" FOREIGN KEY ("UserID")

      REFERENCES public."Users" ("ID") MATCH SIMPLE

      ON UPDATE NO ACTION ON DELETE NO ACTION

)

WITH (

  OIDS=FALSE

);

ALTER TABLE public."TrackerUsers"

  OWNER TO godutch;

