# Success example response

```
$ curl http://localhost:8080/v1/company/7736207543

{"inn":"7736207543","kpp":"770401001","name":"ОБЩЕСТВО С ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТЬЮ \"ЯНДЕКС\"","ceo":"Савиновский Артем Геннадьевич"}
```

# Error example response
```
$ curl http://localhost:8080/v1/company/123

{"error":"некорректный ИНН","code":2,"message":"некорректный ИНН"}
```
