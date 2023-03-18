# curl sample
### whisper
see: https://openai.com/blog/introducing-chatgpt-and-whisper-apis
```
curl https://api.openai.com/v1/audio/transcriptions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F model="whisper-1" \
  -F file="@/path/to/file/openai.mp3"
```

### chat-gpt
see: https://dev.classmethod.jp/articles/chatgpt-api-first-step-for-beginners/
```
curl https://api.openai.com/v1/chat/completions -H "Content-Type: application/json" \
-H "Authorization: Bearer sk-xxxxxxxxxxxxx" \
-d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "自己紹介をして下さい"}]}'
```