# phone_book_json
phone book / store json format / golang
### Этот пакет предоставляет возможности телефонного справочника в формате json

#### Консольные комманды:
- **insert** добавляет запись, пример ```go run phone_book_json insert <first_name string> <last_name string> <phone int>```

- **list** выводит список записей, использование ```go run phone_book_json list``` вывод  ```[
        {
                "first_name": "lesha",
                "last_name": "shirnin",
                "phone": 89015849542,
                "updated_at": "1743505396"
        },
        {
                "first_name": "lesha",
                "last_name": "shirnin",
                "phone": 890158495100,
                "updated_at": "1743506754"
        }
]```
- **remove** удаляет запись с указаным номером, использование ```go run phone_book_json remove <phone_number:int>```
- **search** поиск по номеру телефона использование ```go run phone_book_json search <phone_number:int>``` если запись существует, вернет ```{
                "first_name": "lesha",
                "last_name": "shirnin",
                "phone": 89015849542,
                "updated_at": "1743505396"
        }```,если использовать флаг **--startWith**, тогда будет искать записи, чьи номера начинаються с заданного номера, например ```go run phone_book_json search 8901 --startWith``` вывод ```[
        {
                "first_name": "lesha",
                "last_name": "shirnin",
                "phone": 89015849542,
                "updated_at": "1743505396"
        }, ...
]```
