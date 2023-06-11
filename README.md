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

Теперь про REST сервис. Запускается все это дело на порту 8090. 
Как запускается
```
cd docker
docker compose build
docker compose up infoserver
```

Основные методы:
```
# изменить профиль пользователя
# доступные параметры: name, id, img (название картинки в папке content/img), email, sex
curl http://localhost:8090/user --header "Content-Type: application/json" --request "PUT" -d '{"id": "albert", "name": "bebrik", "image": "img.png"}'

# добавить пользователя
# доступные параметры: name, id, img (название картинки в папке content/img), email, sex
curl http://localhost:8090/user --header "Content-Type: application/json" --request "POST" -d '{"id": "albert", "name": "bebrik", "image": "img.png"}'

# Получить пользователя
# доступные параметры: id
curl http://localhost:8090/user --header "Content-Type: application/json" --request "GET" -d '{"id": "albert"}'

# Удалить пользователя
# доступные параметры: id
curl http://localhost:8090/user --header "Content-Type: application/json" --request "DELETE" -d '{"id": "albert"}'

# Получить всех пользователей
curl http://localhost:8090/users --header "Content-Type: application/json" --request "GET"'

# Сгенерировать pdf по пользователю
# доступные параметры: id
curl http://localhost:8090/genUserPdf --header "Content-Type: application/json" --request "GET" -d '{"id": "albert"}'

# Сгенерировать pdf по всем пользователям
curl http://localhost:8090/genPdf --header "Content-Type: application/json" --request "GET"'
```