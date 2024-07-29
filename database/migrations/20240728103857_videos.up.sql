CREATE TABLE
    IF NOT EXISTS videos (
        id uuid NOT NULL DEFAULT gen_random_uuid(),
        name character varying NOT NULL,
        description character varying,
        video_url character varying NOT NULL,
        createdAt TIMESTAMP NOT NULL DEFAULT now (),
        updatedAt TIMESTAMP NOT NULL DEFAULT now (),
        CONSTRAINT "PK_1c73655f9cfc26a01df74d4e5e9" PRIMARY KEY ("id")
    );
