KAFKA_BIN = docker exec -it kafka-1 /opt/kafka/bin
BS        = --bootstrap-server kafka-1:19092
TOPIC    ?= orders
GROUP    ?= orders-consumers

.PHONY: up down reset ps logs deps topic topics describe groups group produce consume

up:        ## поднять кластер
	docker compose up -d

down:      ## остановить (данные сохраняются)
	docker compose down

reset:     ## остановить и удалить все данные
	docker compose down -v

ps:
	docker compose ps

logs:
	docker compose logs -f kafka-1

deps:      ## скачать Go-зависимости
	go mod tidy

topic:     ## создать топик: make topic TOPIC=orders
	$(KAFKA_BIN)/kafka-topics.sh $(BS) --create --if-not-exists --topic $(TOPIC) --partitions 3 --replication-factor 3

topics:
	$(KAFKA_BIN)/kafka-topics.sh $(BS) --list

describe:  ## партиции, лидеры, реплики, ISR
	$(KAFKA_BIN)/kafka-topics.sh $(BS) --describe --topic $(TOPIC)

groups:
	$(KAFKA_BIN)/kafka-consumer-groups.sh $(BS) --list

group:     ## offset'ы и lag группы
	$(KAFKA_BIN)/kafka-consumer-groups.sh $(BS) --describe --group $(GROUP)

produce:
	go run ./cmd/producer -topic $(TOPIC)

consume:
	go run ./cmd/consumer -topic $(TOPIC) -group $(GROUP)
