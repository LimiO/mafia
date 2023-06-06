Мафия на минималках. Общение между клиентом и сервером происходит
посредством gRPC, где клиент отправляет запросы, а сервер стримит.

Предварительно необходимо запустить rabbitmq
```
cd docker
docker-compose build
docker-compose up rabbitmq
```

Запуск сервера происходит через
```
docker-compose up server
```

Для клиентов необходимо предварительно подготовить окружение
```bash
make env
```

Запуск клиентов происходит через Makefile (так как идет прямое взаимодействие с консолью) и следующие две команды:
```
cd <project_dir>
make player  # создает игрока
make bot  # создает бота, который будет выбирать случайные действия
```

Выбор действий за пользователя происходит через консоль.
Будет предложено N действий, нужно в консоли прописать соответствующее действие.
Если говорят, что "2) end day", то необходимо прописать 2.
