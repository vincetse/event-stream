SHELL := /bin/bash
uri = amqp://guest:guest@localhost:5672
exchange_name = xyz
exchange_type = fanout
routing_key = fanning-out
queue_name_prefix = q
event_count = 5

mq:
	docker run --rm \
		-p 5672:5672 \
		--detach \
		--name mq \
		rabbitmq:3-alpine

consumer:
	$(eval n = 3)
	for ((i=0; i<$(n); i++)); do \
		echo ./event-consumer \
			--uri $(uri) \
			--exchange-name $(exchange_name) \
			--exchange-type $(exchange_type) \
			--queue-name $(queue_name_prefix)-$$i \
		; \
	done \
	| parallel --keep-order --ungroup --jobs $(n)

producer:
	./event-producer \
		--uri $(uri) \
		--exchange-name $(exchange_name) \
		--exchange-type $(exchange_type) \
		--routing-key $(routing_key) \
		--event-count $(event_count)
