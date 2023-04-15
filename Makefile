include .env
NOW := $(shell date +%s)


run:
	go build -o gigeelock cmd/server/main.go && ./gigeelock

run-sample:
	go run cmd/sample/main.go

open-test:
	open ./tmp/test.aiff

curl-whisper:
	curl --request POST \
    --url https://api.openai.com/v1/audio/transcriptions \
    --header 'Authorization: Bearer $(OPENAI_API_KEY)' \
    --header 'Content-Type: multipart/form-data' \
    --form file=@$(TARGET_PATH) \
    --form model=whisper-1 >> ./tmp/$(NOW).txt

list-model:
	curl https://api.openai.com/v1/models \
	-H "Authorization: Bearer $(OPENAI_API_KEY)" \
	| jq '.data | sort_by(.created) | map(.id)'