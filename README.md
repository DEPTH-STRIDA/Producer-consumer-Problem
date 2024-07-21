# Проект "Производитель и Потребитель" на Golang

Данный пример решает проблему производителя и потребителя на языке Golang.

Критической областью является срез `persons`. Работать с ним, добавлять или удалять элементы одновременно может только одна функция. Доступ к срезу ограничивается за счет мьютексов. Максимальный размер слайса ограничивается семафорами писателя и читателя в купе с `cond` из пакета `sync`. Таким образом, в срез может записаться только 5 человек. Код построен таким образом, чтобы соответствовать диаграмме. 💻🔄

## English

This example solves the problem of producer and consumer in the Golang language.

The critical area is the `persons` slice. Only one function can add or remove items at a time. Access to the slice is limited by mutexes. The maximum slice size is limited by writer and reader semaphores coupled with `cond` from the `sync` package. Thus, only 5 people can enroll in the slice. 💻🔄

### Useful Links

- [Go Semaphores](https://aarol.dev/posts/go-semaphores/)
- [Sync Package Documentation](https://pkg.go.dev/sync)
- [Synchronization Primitives in Go](https://medium.com/german-gorelkin/synchronization-primitives-go-8857747d9660)
