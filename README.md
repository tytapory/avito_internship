# Avito Shop
Бекэнд внутреннего магазина авито для сотрудников. Поддерживает регистрацию, авторизацию, покупку мерча, переводы коинов между пользователями и вывод общей информации о пользователе.

## Запуск
Запуск осуществляется в докере. Для запуска можно использовать команду
```make prod```
Эта команда запустит контейнер со стандарнтыми настройками на http://localhost:8080/

## Запуск тестов
Запуск юнит тестов и E2E тестов можно выполнив используя
```make test```
Эта команда поднимает сервер с тестовой базой данных и контейнер с тестами, которые будут выполнены при запуске.

## Нагрузочное тестирование
При желании можно запустить нагрузочное тестирование через K6. Для этого нужно запустить сервер на http://localhost:8080/ и в директории tests/load_test выполнить команду
```k6 run k6_load_test.js```
Этот сценарий выполняется от имени 30 виртуальных пользователей одновременно, каждая итерация состоит из 9 разных запросов на сервер.
На моем ноутбуке этот сценарий выдает примерно 2000 RPS при 0% ошибок и со средним временем ожидания 15 мс.

## Возникшие проблемы
* В доке апи не указано, что нужно делать если запрос не подходит по типу, поэтому ко всем ручкам добавил стандартную ошибку 405 (method not allowed).
* Не указаны ограничения на длину логина и пароля, поставил самые щадящие условия (на случай если у вас есть автотесты). Логин и пароль не пустые и логин меньше 32
