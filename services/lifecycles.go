package services

type lifecycle int

/*
Я собираюсь реализовать жизненный цикл Scoped, используя пакет
context в стандартной библиотеке.
Context будет автоматически создаваться для каждого HTTP-запроса,
полученного сервером, а это означает, что весь код обработки запросов,
который обрабатывает этот request может совместно использовать один
и тот же набор служб, так что, например, одна структура,
предоставляющая информацию о сеансе, может использоваться во
время обработки данного запроса.
*/
const (
	Transient lifecycle = iota
	Singleton
	Scoped
)
