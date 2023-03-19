run:
	go build -o gigeelock cmd/server/main.go && ./gigeelock
open-test:
	open ./tmp/test.aiff

curl-whisper:
	curl --request POST \
    --url https://api.openai.com/v1/audio/transcriptions \
    --header 'Authorization: Bearer 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' \
    --header 'Content-Type: multipart/form-data' \
    --form file=/tmp/target.mp3 \
    --form model=whisper-1