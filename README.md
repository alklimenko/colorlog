# colorlog

Простой логгер с цветовым оформлением для небольших проектов. Предназначен исключительно для ручного просмотра логов, 
поэтому его вывод не загромождён никакой информацией, кроме самой необходимой. В начале каждого дня в лог выводится дата,
в каждой строке с сообщением добавляется только время. Сообщения разного приоритета выделяются разными цветами,
дополнительно, цветом выделяются даты и времена сообщений.

**Пример 1:**

Простое использование
```go
colorlog.Debug("Debug message %d", 1)
colorlog.Info("Info message %s", "2")
colorlog.Warn("Warning message %s", "2")
colorlog.Error("Error message")
colorlog.Fatal("Fatal message")
```

**Пример 2:**

Вывод лога в файл
```go
errLog, _ := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
outLog, _ := os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
colorlog.WithErr(errLog)
colorlog.WithOut(outLog)

colorlog.Debug("Debug message %d", 1)
colorlog.Info("Info message %s", "2")
colorlog.Warn("Warning message %s", "2")
colorlog.Error("Error message")
colorlog.Fatal("Fatal message")
```

**Пример 3:**

Дополнительный лог со своими настройками
```go
log := colorlog.New().WithErr(errLog).WithOut(outLog).WithConfig(cfg)

log.Info("Info message")
log.Debug("Debug message")

colorlog.Info("Info message")
```
