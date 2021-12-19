INSERT INTO dbo.t_flights
select
t.flight_id + 200 as flight_id,
flight_name,
airline_id,
ticket_num_economy_class,
ticket_num_pr_economy_class,
ticket_num_business_class,
ticket_num_first_class,
cost_economy_class_rub,
cost_pr_economy_class_rub,
cost_business_class_rub,
cost_first_class_rub,
aircraft_model_id,
t.landing_airport_id as departure_airport_id,
t.departure_airport_id as landing_airport_id,
t.departure_time + interval '3 days' as departure_time,
t.landing_time + interval '3 days' as landing_time,
max_luggage_weight_kg,
cost_luggage_weight_rub,
max_hand_luggage_weight_kg,
cost_hand_luggage_weight_rub,
wifi_flg,
food_flg,
usb_flg,
change_dttm
from dbo.t_flights as t;

select setval('dbo.t_flights_flight_id_seq', (select max(flight_id) from dbo.t_flights));