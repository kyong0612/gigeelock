# curl sample
### whisper
see: https://openai.com/blog/introducing-chatgpt-and-whisper-apis 
see: https://platform.openai.com/docs/guides/speech-to-text
```
curl --request POST \
    --url https://api.openai.com/v1/audio/transcriptions \
    --header 'Authorization: Bearer TOKEN' \
    --header 'Content-Type: multipart/form-data' \
    --form file=@/tmp/target.mp3 \
    --form model=whisper-1
```


### chat-gpt
see: https://dev.classmethod.jp/articles/chatgpt-api-first-step-for-beginners/
```
curl https://api.openai.com/v1/chat/completions -H "Content-Type: application/json" \
-H "Authorization: Bearer sk-xxxxxxxxxxxxx" \
-d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "自己紹介をして下さい"}]}'
```