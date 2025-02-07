# Ozon-MyHub Проект 
**Ozon-MyHub** — это система для управления постами и комментариями, предоставляющая возможность работы как с базой данных Postgres, так и с in-memory хранилищем. 
## Запуск проекта
 Для запуска проекта используйте команду:

```
bash
docker-compose up --build -d
```
Данный проект поддерживает два типа хранилища: *Postgres* и *in-memory* хранилище.

*Использование in-memory хранилища*
Если вы хотите использовать in-memory хранилище, в файле .env установите значение:

> env
> 
> USE_POSTGRES=false

*Использование Postgres хранилища*
Если вы хотите использовать Postgres хранилище, в файле .env установите значение:

> env
> 
> USE_POSTGRES=true

Пример конфигурации для Postgres:

> env
> 
> USE_POSTGRES=true DB_HOST=myhub-postgres
> 
> DB_PORT=5432
>  DB_USER=test
>   DB_PASSWORD=myhub_password
>    DB_NAME=myhub_db
> APP_PORT=8080

Взаимодействие с пользовательскими данными:
Если вы выбрали локальное хранилище, в логах консоли появится информация о трёх созданных пользователях: Для хранилища на локальной машине  по умолчанию используются следующие пользователи:
При желании тип ID можно поменять 

    myhub-app       | Пользователь: Alice (ID: 1)
    myhub-app       | Пользователь: Bob (ID: 2)
    myhub-app       | Пользователь: Charlie (ID: 3)
    
Если вы выбрали хранилище Postgres, система сгенерирует 3 случайных новых записи о пользователях в БД в таблице users. Для хранилища в БД  передаваемый ID пользователя должен быть типа UUID. 

    myhub-app       | Пользователь: test_user_1 (ID: 7b465018-ea3d-4bf3-be46-ec2d19100c84)
    myhub-app       | Пользователь: test_user_2 (ID: 5a0cf470-f003-417b-8276-568a493bef7e)
    myhub-app       | Пользователь: test_user_3 (ID: 3725e2a8-17ce-4267-82db-b34fdc1f4fbe)
Именно эти данные нужно использовать в полях userID при тестировании работоспособности системы.

## Основные эндпоинты:
> 
> _(limit:10, offset: 0) - настраиваемые параметры_



Если вы хотите использовать Postgres хранилище, в файле `.env` установите значение `POSTGRES=true`.

## Основные эндпоинты

_(limit:10, offset: 0) - настраиваемые параметры_

**Получение списка постов:**
```graphql
query {
  posts(limit:10, offset: 0) {
    posts {
      id
      title
      content
      commentsAllowed
    }
    totalCount
    hasNextPage
  }
}

2. Создание нового поста
graphql

mutation {
  createPost(
    userId: "1"
    title: "22222222"
    content: "новыйтест"
    commentsAllowed: true
  ) {
    id
    title
    content
    commentsAllowed
  }
}
3. Просмотр комментариев к посту с возможностью пагинации
graphql

query {
  post(id: "0ff89a9b-9f9b-403f-b4e4-7e41ecb218e5") {
    id
    title
    content
    commentsAllowed
    comments(limit: 10, offset: 0) {
      id
      content
      user {
        nickname
        id
      }
    }
  }
}
4. Получение комментариев с вложенностью (ответы на комментарии)
Если хотите получить вложенные комментарии (ответы на комментарии), добавьте параметр replies:

graphql

query {
  post(id: "cad0e69e-bb73-45fd-a4c0-2b49dacda01b") {
    id
    title
    content
    commentsAllowed
    comments(limit: 10, offset: 0) {
      id
      content
      user {
        nickname
        id
      }
      replies {
        id
        content
        user {
          nickname
          id
        }
        replies {
          id
          content
          user {
            nickname
            id
          }
        }
      }
    }
  }
}
5. Добавление нового комментария к посту
graphql

mutation {
  createCommentForPost(
    userId: "1"
    postId: "cad0e69e-bb73-45fd-a4c0-2b49dacda01b"
    content: "2 коммент"
  ) {
    id
    content
    postId
  }
}
6. Добавление ответа на комментарий (создание новой ветки комментариев)
graphql

mutation {
  createReplyForComment(
    userId: "1",
    parentCommentId: "d94e22dd-be91-40d2-a8c9-323b84630417",
    content: "Это мой 5 ответ на комментарий"
  ) {
    id
    content
  }
}
7. Поменять статус поста (разрешить или запретить комментарии)
graphql

mutation {
  toggleCommentsForPost(
    postId: "bbd7e741-72e8-4325-a66f-a8624c347500",
    userId: "843b3f66-5421-44e8-942a-54e1f0384ff7"
  ) {
    id
    title
    commentsAllowed
  }
}
```
Для тестов используйте команду 
(тест внутреннего хранилища)

go test -v ./internal/storage 

В рамках данной системы реализованы система пагинации, основные функции создания и получения постов и комментариев, а также валидация. В проекте используются in-memory хранилище и PostgreSQL в качестве альтернативного хранилища. 
Также мне бы хотелось добавить, но не хватило времени:
1. Систему аутентификации и авторизации при помощи директивы isAuthenticated и DataLoaderMiddleware.
2. В связке с авторизацией описала бы механизм Subscription.
3. Проблему n+1 запросов я бы решила через библиотеку dataloaden.
4.  Ограничила бы максимальную сложность запроса через подобъекты запроса, используя handler.ComplexityLimit(). 

Жду обратной связи!