
uri = amqp://guest:guest@localhost:5672
routing_key = hello
event_count = 5

mq:
	docker run --rm \
		-p 5672:5672 \
		--detach \
		--name mq \
		rabbitmq:3-alpine

consumer:
	./event-consumer \
		--uri $(uri) \
		--queue-name $(routing_key)

producer:
	./event-producer \
		--uri $(uri) \
		--routing-key $(routing_key) \
		--event-count $(event_count)
