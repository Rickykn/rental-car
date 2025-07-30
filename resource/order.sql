CREATE TABLE public.orders (
                               id uuid NOT NULL DEFAULT gen_random_uuid(),
                               car_id uuid NOT NULL,
                               order_date date NOT NULL,
                               pickup_date date NOT NULL,
                               dropoff_date date NOT NULL,
                               pickup_location bpchar(50) NOT NULL,
                               dropoff_location bpchar(50) NOT NULL,
                               CONSTRAINT orders_pkey PRIMARY KEY (id),
                               CONSTRAINT orders_car_id_fkey FOREIGN KEY (car_id) REFERENCES public.cars(id) ON DELETE CASCADE
);