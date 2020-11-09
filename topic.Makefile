SHELL := /bin/bash
uri = amqp://guest:guest@localhost:5672
exchange_name = ace
exchange_type = topic
queue_name_prefix = q
event_count = 1

mq:
	docker run --rm \
		-p 5672:5672 \
		--detach \
		--name mq \
		rabbitmq:3-alpine

consumer:
	$(eval n = 4)
	{ \
		echo ./event-consumer \
			--uri $(uri) \
			--exchange-name $(exchange_name) \
			--exchange-type $(exchange_type) \
			--routing-key "americas.south.*" \
			--queue-name "americas_south" \
		&& \
		echo ./event-consumer \
			--uri $(uri) \
			--exchange-name $(exchange_name) \
			--exchange-type $(exchange_type) \
			--routing-key "americas.north.*" \
			--queue-name "americas_north" \
		&& \
		echo ./event-consumer \
			--uri $(uri) \
			--exchange-name $(exchange_name) \
			--exchange-type $(exchange_type) \
			--routing-key "americas.central.*" \
			--queue-name "americas_central" \
		&& \
		echo ./event-consumer \
			--uri $(uri) \
			--exchange-name $(exchange_name) \
			--exchange-type $(exchange_type) \
			--routing-key "americas.#" \
			--queue-name "americas" \
	; } \
	| parallel --keep-order --ungroup --jobs $(n)

producer:
	$(eval n = 3)
	$(eval routing_keys = americas.south.brazil americas.north.mexico americas.north americas.north.canada asia.south.india americas)
	for routing_key in $(routing_keys); do \
		echo ./event-producer \
			--uri $(uri) \
			--exchange-name $(exchange_name) \
			--exchange-type $(exchange_type) \
			--routing-key $$routing_key \
			--event-count $(event_count); \
	done \
	| parallel --keep-order --ungroup --jobs $(n)
