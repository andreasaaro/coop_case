# Requirements go, docker desktop, docker cli, kafkacat

APP=coop_case
TOPIC_OUT=mastodon_topic

build: test
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.AppName=$(APP) -X main.Version=vDEV" -o bin/$(APP)

modules:
	go mod tidy

test: modules
	gosec ./...
	go fmt ./...
	go test ./... -timeout 5s --cover
	go vet ./...

docker-kafka: 
	docker compose up zookeeper kafka -d
	testdata/wait_for_kafka.sh kafka 0.0.0.0 9092

docker-mastodon: 
	docker compose up mastodon -d 

docker-kill: 
	docker compose down

clean: docker-kafka
	docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --zookeeper zookeeper:2181 --delete --topic $(TOPIC_OUT) || :
	sleep 2

kafka-topics clean:
	docker exec kafka /opt/kafka/bin/kafka-topics.sh --zookeeper zookeeper:2181 --create --topic $(TOPIC_OUT) --partitions 1 --replication-factor 1

get_timeline: 
	curl https://mastodon.social/api/v1/timelines/public\?limit\=2 | jq

get_timeline2: 
	curl https://mastodon.social/api/v1/timelines/public?min_id=5&max_id=10 > t.json


consume-kafka: 
	$(info echo consuming $(TOPIC_OUT))
	kcat -G $(APP) -C  -b localhost:9094 -t $(TOPIC_OUT) -e -u -q -f "%R%s"

inspect-mastodon-to-kafka: 
	docker logs -f mastodon-to-kafka 

cycle: build docker-kill docker-kafka kafka-topics docker-mastodon inspect-mastodon-to-kafka
