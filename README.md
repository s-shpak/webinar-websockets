# ws

Код для демо-приложения находится в папке `ws-chat`. Команды ниже должны выполняться из нее.

Сгенерируйте сертификат для сервера:

```bash
make generate-certs
```

Скомпилируйте сервер и клиент:

```bash
make build
```

Запустите сервер:

```bash
./server -p 4443
```

Запустите двух клиентов:

```bash
./client -m "Привет от Васи!"

# в другом терминале
./client -m "Привет от Пети!"
```

# Дополнительная литература

- Removing HTTP2 Server push from Chrome: https://developer.chrome.com/blog/removing-push/
- How HTML5 Web Sockets Interact With Proxy Servers: https://www.infoq.com/articles/Web-Sockets-Proxy-Servers/
- RFC 8441: https://www.rfc-editor.org/rfc/rfc8441.html
- Миллион WebSocket и Go: https://habr.com/en/companies/vk/articles/331784/
- Writing WebSocket servers: https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers