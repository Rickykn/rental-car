CREATE TABLE public.cars (
                             id uuid NOT NULL DEFAULT gen_random_uuid(),
                             car_name bpchar(50) NOT NULL,
                             day_rate float8 NOT NULL,
                             month_rate float8 NOT NULL,
                             image bpchar(256) NOT NULL,
                             status varchar(20) NOT NULL DEFAULT 'available'::character varying,
                             CONSTRAINT cars_car_name_key UNIQUE (car_name),
                             CONSTRAINT cars_pkey PRIMARY KEY (id)
);