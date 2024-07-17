// Статичесий анализатор кода. Проверят на возможные ошибки (см. ниже)
package main

// сборка: из корневой директории:
//  go build cmd/staticlint/multichecker.go
// запуск проверки:
//  ./multichecker ./...

// Используемые анализаторы:
// Собственные:
// osexitchekc - запрещаtn использовать прямой вызов os.Exit в функции main пакета main
//
// Сторонние:
// errcheck (https://github.com/kisielk/errcheck) - проверка на наличие непроверенных ошибок в коде Go. В некоторых случаях эти непроверенные ошибки могут быть критическими.
// ineffassign (https://github.com/gordonklaus/ineffassign) - определяет, когда назначения существующим переменным не используются.
// wrapcheck (https://github.com/tomarrell/wrapcheck) - проверяет, то ошибки из внешних пакетов обернуты, чтобы определить источник ошибки.

// Из паекта "golang.org/x/tools/go/analysis/passes/
// appends         проверяет, что в append только одна переменная, т.е. не передается значение для добавления в слайс.
// asmdecl         проверяет, что файлы сборки соответствуют объявлениям Go.
// assign          проверяет бесполезные присваивания, например x = x
// atomic          проверяет распространенные ошибка использоsвания sync/atomic
// atomicalign     чтш-то с разрядностью, не понял
// bools           проверяет распространенные ошибки, связанные с использованием bool
// buildssa        ??
// buildtag        проверяет тэги сборки(buildtag)
// cgocall         проверяет нарушения правил передачи указателей cgo
// composite       ?? преаеряет наличие составных литералов без ключа
// copylock        проверяет блокировки, которые были установлены по-значению
// ctrlflow        ??
// deepequalerrors проверяет использование reflect.deepEqual со значениями ошибок.
// defers          проверяет ошибочное использованеие defer
// directive       проверяет известные директивы инструментов
// errorsas        проверяте, что второй аргумент в errors.As - это указатель на тип, реализующий интерфейс ошибки
// // fieldalignment      обнаруживает структуры, которые использовали бы меньше памяти, если бы их поля были отсортированы
// findcall            ?
// framepointer        ?
// httpmux             ?
// httpresponse        проверяет ошибки в HTTP-ответах
// ifaceassert         находит невозожное приведение интерфейса-в-интерфейс, например из за одинаковых имен, но разных сигнатур приводимых интерфейсов
// inspect             ?
// loopclosure         проверяет наличие ссылок на переменные цикла, входящие во вложенные функции.
// lostcancel          проверят выхов отмены контекста
// nilfunc             проверяет бесполезное сравнение функции с nil (fun == nil)
// nilness             проверьте на избыточные или невозможных сравнений с nil
// printf              проверяет согласованность строк и аргументов формата в printf
// reflectvaluecompare  проверяет случайное использование == или reflect.deepEqual для сравнения значений reflect.Value
// shadow              проверяет сокрытие пременных
// shift               проверяет наличие сдвигов, превышающих диапазон целого числа.
// sigchanyzer         проверяет неправильное использование небуферизованного сигнала в качестве аргумента для signal.Notify.
// sortslice           проверяет наличие вызовов sort.Slice, которые не используют тип slice в качестве первого аргумента.
// stdmethods          проверяет сигнатуры методов известных интерфейсов.
// stdversion          проверяет библиотечные символы, которые являются "слишком новыми" для действующей версии Go
// stringintconv       проверяет наличие преобразований string(int), что приведет к получению UTF-8 символа, а не преобразования числа в строку как есть
// structtag           проверяет правильность формирования тегов структурных полей.
// testinggoroutine    проверят вызовов Fatal из тестовой программы.
// tests               проверяет типичные ошибки при использовании тестов и примеров.
// timeformat          проверяет использование time.Format или time.Parse вызовы с неправильным форматом.
// unmarshal           проверят, что в unmarshal передаются значения отличные от указателя или интерфейса.
// unreachable         проверяет, что в коде есть недостижимые участки
// unsafeptr           проверяет, нет ли недопустимых преобразований uintptr в unsafe.Pointer
// unusedresult        проверяет неиспользуемые результаты вызовов некоторых функций.
// unusedwrite          проверяет наличие неиспользуемых записей
// usesgenerics        проверяет, использует ли программа generic-и

// Используемые анализаторы из staticcheck.io

// SA	staticcheck
// SA1	Неправильное использование стандартной библиотеки
// SA1000	неправильное регулярное выражение
// SA1001	неправильный шаблон
// SA1002	направильное форматирование в time.Parse
// SA1003	неподдерживаемый аргумент в encoding/binary
// SA1004	подозрительно маленькая нетипизированая константа в time.Sleep
// SA1005	недопустимый первый аргумент для exec.Command
// SA1006	Printf с динамическим первым аргументом и без дополнительных аргументов
// SA1007	недопустимый URL-адрес в net/url.Parse
// SA1008	неканонический ключ в http.Header map
// SA1010	(*regexp.Regexp).FindAll вызывается с n == 0, который всегда будет возвращать нулевые результаты
// SA1011	различные методы в пакете strings предполагают наличие допустимого значения UTF-8, но вводятся недопустимые данные
// SA1012	в функцию передается nil context.Context, возможно стоит использовать context.TODO
// SA1013	io.Seeker.Seek вызывается с whence первым аргументов, но этот аргумент должен быть вторым
// SA1014	значение, не являющееся указателем, передается в Unmarshal или Decode
// SA1015	Using time.Tick in a way that will leak. Consider using time.NewTicker, and only use time.Tick in tests, commands and endless functions
// SA1016	перехват неперехвачиваемого сигнала
// SA1017	каналы, используемые с os/signal.Notify, должны быть буферизованы
// SA1018	strings.Replace вызываемая с n == 0, ничего не делает
// SA1019	использование устаревшей функции, переменной, константы или поля
// SA1020	использование неправильной связки host:port в net.Listen-подобных функциях
// SA1021	использование bytes.Equal для сравнения двух net.IP
// SA1023	модификация буфера в io.Writer реализации
// SA1024	вырезанные строки содержат дубликаты символов
// SA1025	It is not possible to use (*time.Timer).Reset’s return value correctly
// SA1026	Cannot marshal channels or functions
// SA1027	? атомарный доступ к 64-bit переменной должен быть выронен по 64-битно разрядности
// SA1028	sort.Slice должен быть использован только со слайсами
// SA1029	неподходящий ключ при вызова context.WithValue
// SA1030	недопустимый аргумент при вызовва strconv
// SA1031	? Overlapping byte slices passed to an encoder
// SA1032	неправильный порядок аргументов в errors.Is

// SA2	Проблемы с конкурентностью
// SA2000	sync.WaitGroup.Add вызываемая в горутине, приводит к состоянию гонки
// SA2001	пустая критическая секция, возможно нужно вызвать unlock в defer?
// SA2002	вызов testing.T.FailNow или SkipNow в горутение запрещен
// SA2003	Deferred Lock right after locking, likely meant to defer Unlock instead

// SA3	Проблемы в тестах
// SA3000	TestMain doesn’t call os.Exit, hiding test failures
// SA3001	присвоение к b.N в бенчмарках искажает результаты

// SA4	Бесполезный код (который ничего не делает)
// SA4000	Binary operator has identical expressions on both sides
// SA4001	&*x gets simplified to x, it does not copy x
// SA4003	Comparing unsigned values against negative values is pointless
// SA4004	The loop exits unconditionally after one iteration
// SA4005	Field assignment that will never be observed. Did you mean to use a pointer receiver?
// SA4006	A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?
// SA4008	The variable in the loop condition never changes, are you incrementing the wrong variable?
// SA4009	A function argument is overwritten before its first use
// SA4010	The result of append will never be observed anywhere
// SA4011	Break statement with no effect. Did you mean to break out of an outer loop?
// SA4012	Comparing a value against NaN even though no value is equal to NaN
// SA4013	Negating a boolean twice (!!b) is the same as writing b. This is either redundant, or a typo.
// SA4014	An if/else if chain has repeated conditions and no side-effects; if the condition didn’t match the first time, it won’t match the second time, either
// SA4015	Calling functions like math.Ceil on floats converted from integers doesn’t do anything useful
// SA4016	Certain bitwise operations, such as x ^ 0, do not do anything useful
// SA4017	Discarding the return values of a function without side effects, making the call pointless
// SA4018	Self-assignment of variables
// SA4019	Multiple, identical build constraints in the same file
// SA4020	Unreachable case clause in a type switch
// SA4021	x = append(y) is equivalent to x = y
// SA4022	Comparing the address of a variable against nil
// SA4023	Impossible comparison of interface value with untyped nil
// SA4024	Checking for impossible return value from a builtin function
// SA4025	Integer division of literals that results in zero
// SA4026	Go constants cannot express negative zero
// SA4027	(*net/url.URL).Query returns a copy, modifying it doesn’t change the URL
// SA4028	x % 1 is always zero
// SA4029	Ineffective attempt at sorting slice
// SA4030	Ineffective attempt at generating random number
// SA4031	Checking never-nil value against nil
// SA4032	Comparing runtime.GOOS or runtime.GOARCH against impossible value

// SA5	Correctness issues
// SA5000	Assignment to nil map
// SA5001	Deferring Close before checking for a possible error
// SA5002	The empty for loop (for {}) spins and can block the scheduler
// SA5003	Defers in infinite loops will never execute
// SA5004	for { select { ... with an empty default branch spins
// SA5005	The finalizer references the finalized object, preventing garbage collection
// SA5007	Infinite recursive call
// SA5008	Invalid struct tag
// SA5009	Invalid Printf call
// SA5010	Impossible type assertion
// SA5011	Possible nil pointer dereference
// SA5012	Passing odd-sized slice to function expecting even size
//
// SA6	Проблемы с производительностью
// SA6000	Using regexp.Match or related in a loop, should use regexp.Compile
// SA6001	Missing an optimization opportunity when indexing maps by byte slices
// SA6002	Storing non-pointer values in sync.Pool allocates memory
// SA6003	Converting a string to a slice of runes before ranging over it
// SA6005	Inefficient string comparison with strings.ToLower or strings.ToUpper
// SA6006	Using io.WriteString to write []byte
//
// SA9	Сомнительные конструкции кода, которые могут быть ошибочными
// SA9001	Defers in range loops may not run when you expect them to
// SA9002	Using a non-octal os.FileMode that looks like it was meant to be in octal.
// SA9003	Empty body in an if or else branch
// SA9004	Only the first constant has an explicit type
// SA9005	Trying to marshal a struct with no public fields nor custom marshaling
// SA9006	Dubious bit shifting of a fixed size integer value
// SA9007	Deleting a directory that shouldn’t be deleted
// SA9008	else branch of a type assertion is probably not reading the right value
// SA9009	Ineffectual Go compiler directive
//
// S	simple
// S1	Код может быть проще
// S1000	Use plain channel send or receive instead of single-case select
// S1001	Replace for loop with call to copy
// S1002	Omit comparison with boolean constant
// S1003	Replace call to strings.Index with strings.Contains
// S1004	Replace call to bytes.Compare with bytes.Equal
// S1005	Drop unnecessary use of the blank identifier
// S1006	Use for { ... } for infinite loops
// S1007	Simplify regular expression by using raw string literal
// S1008	Simplify returning boolean expression
// S1009	Omit redundant nil check on slices, maps, and channels
// S1010	Omit default slice index
// S1011	Use a single append to concatenate two slices
// S1012	Replace time.Now().Sub(x) with time.Since(x)
// S1016	Use a type conversion instead of manually copying struct fields
// S1017	Replace manual trimming with strings.TrimPrefix
// S1018	Use copy for sliding elements
// S1019	Simplify make call by omitting redundant arguments
// S1020	Omit redundant nil check in type assertion
// S1021	Merge variable declaration and assignment
// S1023	Omit redundant control flow
// S1024	Replace x.Sub(time.Now()) with time.Until(x)
// S1025	Don’t use fmt.Sprintf("%s", x) unnecessarily
// S1028	Simplify error construction with fmt.Errorf
// S1029	Range over the string directly
// S1030	Use bytes.Buffer.String or bytes.Buffer.Bytes
// S1031	Omit redundant nil check around loop
// S1032	Use sort.Ints(x), sort.Float64s(x), and sort.Strings(x)
// S1033	Unnecessary guard around call to delete
// S1034	Use result of type assertion to simplify cases
// S1035	Redundant call to net/http.CanonicalHeaderKey in method call on net/http.Header
// S1036	Unnecessary guard around map access
// S1037	Elaborate way of sleeping
// S1038	Unnecessarily complex way of printing formatted string
// S1039	Unnecessary use of fmt.Sprint
// S1040	Type assertion to current type
//
// ST	стилистика
// ST1	Стилистические пробелмы
// ST1000	Incorrect or missing package comment
// ST1001	Dot imports are discouraged
// ST1003	Poorly chosen identifier
// ST1005	Incorrectly formatted error string
// ST1006	Poorly chosen receiver name
// ST1008	A function’s error value should be its last return value
// ST1011	Poorly chosen name for variable of type time.Duration
// ST1012	Poorly chosen name for error variable
// ST1013	Should use constants for HTTP error codes, not magic numbers
// ST1015	A switch’s default case should be the first or last case
// ST1016	Use consistent method receiver names
// ST1017	Don’t use Yoda conditions
// ST1018	Avoid zero-width and control characters in string literals
// ST1019	Importing the same package multiple times
// ST1020	The documentation of an exported function should start with the function’s name
// ST1021	The documentation of an exported type should start with type’s name
// ST1022	The documentation of an exported variable or constant should start with variable’s name
// ST1023	Redundant type in variable declaration
//
// QF	quickfix
// QF1	Quickfixes
// QF1001	Apply De Morgan’s law
// QF1002	Convert untagged switch to tagged switch
// QF1003	Convert if/else-if chain to tagged switch
// QF1004	Use strings.ReplaceAll instead of strings.Replace with n == -1
// QF1005	Expand call to math.Pow
// QF1006	Lift if+break into loop condition
// QF1007	Merge conditional assignment into variable declaration
// QF1008	Omit embedded fields from selector expression
// QF1009	Use time.Time.Equal instead of == operator
// QF1010	Convert slice of bytes to string when printing it
// QF1011	Omit redundant type from variable declaration
// QF1012	Use fmt.Fprintf(x, ...) instead of x.Write(fmt.Sprintf(...))
