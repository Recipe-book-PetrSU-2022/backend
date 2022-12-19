# Backend

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Recipe-book-PetrSU-2022/backend/go.yml?branch=main&logo=go&logoColor=white&style=for-the-badge)
![GitHub](https://img.shields.io/github/license/Recipe-book-PetrSU-2022/backend?style=for-the-badge)

Структура приложения:
- models - модели для базы данных
- claims - структуры данных о пользователе в JWT
- user_handlers.go - обработчик запросов для пользователя
- recipe_handlers.go - обработчик запросов для рецептов
- security.go - фунцкии для обработки паролей
- responses.go - структуры для создания ответов сервера
