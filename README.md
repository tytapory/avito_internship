### Возникшие проблемы
* В доке апи не указано, что нужно делать если запрос не подходит по типу, возвращаю стандартную ошибку 405
* Не указаны ограничения на длину логина и пароля, поставил самые щадящие (на случай если у вас есть автотесты которые могут упасть на моих ограничениях). Логин и пароль не пустые и логин меньше 32