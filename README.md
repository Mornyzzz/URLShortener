### Как это работает:
    
1. **Docker**:
   - Сборка образа: `docker build -t myapp .`
   - Запуск с Postgres: `docker run myapp` или `docker run -e STORAGE_TYPE=postgres myapp`
   - Запуск с Redis (in-memory): `docker run -e STORAGE_TYPE=redis myapp`

3. **Docker Compose**:
   - Сборка: `docker-compose build`
   - Запуск с Postgres: `docker-compose up` или `STORAGE_TYPE=postgres docker-compose up`
   - Запуск с Redis: `STORAGE_TYPE=redis docker-compose up`
