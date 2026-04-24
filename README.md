# colorlog

Простой логгер с цветовым оформлением для небольших проектов. Предназначен исключительно для ручного просмотра логов, 
поэтому его вывод не загромождён никакой информацией, кроме самой необходимой. В начале каждого дня в лог выводится дата,
в каждой строке с сообщением добавляется только время. Сообщения разного приоритета выделяются разными цветами,
дополнительно, цветом выделяются даты и времена сообщений.

**Пример 1:**

Простое использование
```go
colorlog.Debugf("Debug message %d", 1)
colorlog.Infof("Info message %s", "2")
colorlog.Warnf("Warning message %s", "2")
colorlog.Error("Error message")
colorlog.Fatal("Fatal message")
```

**Пример 2:**

Вывод лога в файл
```go
outLog, _ := os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
colorlog.WithOut(outLog)

colorlog.Debugf("Debug message %d", 1)
colorlog.Infof("Info message %s", "2")
colorlog.Warnf("Warning message %s", "2")
colorlog.Error("Error message")
colorlog.Fatal("Fatal message")
```

**Пример 3:**

Дополнительный лог со своими настройками
```go
log := colorlog.New().WithOut(outLog).WithConfig(cfg)

log.Info("Info message")
log.Debug("Debug message")

colorlog.Info("Info message")
```

**Пример 4:**

Ротация логов:
```go
rot := rotator.NewBuilder().
	WithStrategy(rotator.StrategySize).
	WithSize(1*1024*1024).
	WithCount(5).
	Build()
log := colorlog.New().
	WithOut(rot)
```