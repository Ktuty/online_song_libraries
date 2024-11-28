<div align="center">
  <a href="https://git.io/typing-svg"><img src="https://readme-typing-svg.herokuapp.com?font=Tektur&size=40&duration=4000&color=28DAFF&center=true&vCenter=true&width=435&height=100&lines=online_song_libraries" alt="Typing SVG" /></a>
</div>

# Технологии

- Golang
- PostgreSQL
- RESTful API
- CleanCode

# Подготовка проекта

1. **Клонирование репозитория**:
   ```sh
   git clone https://github.com/Ktuty/online_song_libraries

2. **Установка зависимостей в проекте**:
   ```go
   // из корневой дирректории проетка
   go mod tidy

3. **Проверка и изменение файлов конфигурации .env:**
   ```env
   #Example:
   
    port="8080"

    #URL="external_api"

   настрока уровеня логирования проекта:
    #LOG_LEVEL="debug"
    LOG_LEVEL="info"
    
    DB_NAME="your_database"
    DB_PORT="5432"
    DB_HOST="your_host"
    DB_PASS="your_password"
    DB_SSLMODE="disable"
    DB_USER="postgres"


4. **Проект готов к запуску:**
   ```go
   // из корневой дириктории
   go run cmd/main.go

5. **Ссылка на сваггер:**
   ```sh
    http://localhost:8080/swagger/index.html

# Миграции
1. **Подготовка утилиты Make для Windows 10**:
 * Установка Scoop:
     ```sh
      irm get.scoop.sh -outfile 'install.ps1'
      .\install.ps1 -RunAsAdmin

  * Установка Migrate:
     ```sh
     scoop install migrate

  * Установка Make:
     ```sh
     @powershell -NoProfile -ExecutionPolicy unrestricted -Command "iex ((new-object net.webclient).DownloadString('https://chocolatey.org/install.ps1'))" && SET PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin
     choco install make

2. **Make методы для миграций вручную**
 * путь к бд будет создан из конфигурации .env
     ```env
     #Example:
      DB_NAME="your_database"
      DB_PORT="5432"
      DB_HOST="your_host"
      DB_PASS="your_password"
      DB_SSLMODE="disable"
      DB_USER="postgres"
     
  * Создание файлов для миграции в ./schema
    ```sh
     make create-migration MIGRATION_NAME=<name>
   
  * migration up __(миграции поднимаются из ./schema при запуске программы)__
     ```sh
     make migrate-up
     
  * migration down __(параметр степ указывает на количество миграций для отката)__
    ```sh
    make migrate-down STEPS=<number>

  * Справка по методам
    ```sh
    make help
  
