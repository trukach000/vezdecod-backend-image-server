Для запуска используейте start.sh и start-docker.sh 
(версии with_prefix нужны моего хостинга когда есть прокси вебсервер для множества проектов на одном домене)
для запуска на винде достаточно команды на проброс порта и 
docker run -dp 30001:80 -e DB_HOST=host.docker.internal imloader 

У проекта также при запуску хостится swagger 2.0 по пути /swagger/index.html


Основные руты
POST /upload - загружает JPG картинки в соотвествии с требованиями заданий
GET /get/{id}?scale=x - возвращает картинку по ID (и масштабирует при необходимости)


Тестовый скрипт написан на питоне 3: test-script.py
он делает загрузку изображения из /testdata/ и потом запрашивает его же по полученной идентификатору
и проверяет что ответ 200


python3 ./test-script.py http://127.0.0.1 30001 ""