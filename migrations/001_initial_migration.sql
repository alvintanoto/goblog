CREATE TABLE IF NOT EXISTS public."users"
(
    id bigint NOT NULL DEFAULT nextval('user_id_seq'::regclass),
    username text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT unique_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS public.post
(
    id bigint NOT NULL DEFAULT nextval('post_id_seq'::regclass),
    title text COLLATE pg_catalog."default",
    content text COLLATE pg_catalog."default",
    is_public boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by bigint NOT NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_by bigint,
    is_deleted boolean NOT NULL DEFAULT false,
    CONSTRAINT post_pkey PRIMARY KEY (id),
    CONSTRAINT post_created_by_fkey FOREIGN KEY (created_by)
        REFERENCES public."users" (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT post_updated_by_fkey FOREIGN KEY (updated_by)
        REFERENCES public."users" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);
