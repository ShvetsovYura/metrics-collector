package main

import (
	osexitcheck "github.com/ShvetsovYura/metrics-collector/cmd/staticlint/osexitcheck"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/tomarrell/wrapcheck/wrapcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	saChecks := map[string]bool{
		"SA1000": true, //	Invalid regular expression
		"SA1001": true, //	Invalid template
		"SA1002": true, //	Invalid format in time.Parse
		"SA1003": true, //	Unsupported argument to functions in encoding/binary
		"SA1004": true, //	Suspiciously small untyped constant in time.Sleep
		"SA1005": true, //	Invalid first argument to exec.Command
		"SA1006": true, //	Printf with dynamic first argument and no further arguments
		"SA1007": true, //	Invalid URL in net/url.Parse
		"SA1008": true, //	Non-canonical key in http.Header map
		"SA1010": true, //	(*regexp.Regexp).FindAll called with n == 0, which will always return zero results
		"SA1011": true, //	Various methods in the strings package expect valid UTF-8, but invalid input is provided
		"SA1012": true, //	A nil context.Context is being passed to a function, consider using context.TODO instead
		"SA1013": true, //	io.Seeker.Seek is being called with the whence constant as the first argument, but it should be the second
		"SA1014": true, //	Non-pointer value passed to Unmarshal or Decode
		"SA1015": true, //	Using time.Tick in a way that will leak. Consider using time.NewTicker, and only use time.Tick in tests, commands and endless functions
		"SA1016": true, //	Trapping a signal that cannot be trapped
		"SA1017": true, //	Channels used with os/signal.Notify should be buffered
		"SA1018": true, //	strings.Replace called with n == 0, which does nothing
		"SA1019": true, //	Using a deprecated function, variable, constant or field
		"SA1020": true, //	Using an invalid host:port pair with a net.Listen-related function
		"SA1021": true, //	Using bytes.Equal to compare two net.IP
		"SA1023": true, //	Modifying the buffer in an io.Writer implementation
		"SA1024": true, //	A string cutset contains duplicate characters
		"SA1025": true, //	It is not possible to use (*time.Timer).Reset’s return value correctly
		"SA1026": true, //	Cannot marshal channels or functions
		"SA1027": true, //	Atomic access to 64-bit variable must be 64-bit aligned
		"SA1028": true, //	sort.Slice can only be used on slices
		"SA1029": true, //	Inappropriate key in call to context.WithValue
		"SA1030": true, //	Invalid argument in call to a strconv function
		"SA1031": true, //	Overlapping byte slices passed to an encoder
		"SA1032": true, //	Wrong order of arguments to errors.Is

		"SA2000": true, //	sync.WaitGroup.Add called inside the goroutine, leading to a race condition
		"SA2001": true, //	Empty critical section, did you mean to defer the unlock?
		"SA2002": true, //	Called testing.T.FailNow or SkipNow in a goroutine, which isn’t allowed
		"SA2003": true, //	Deferred Lock right after locking, likely meant to defer Unlock instead
		"SA3000": true, //	TestMain doesn’t call os.Exit, hiding test failures
		"SA3001": true, //	Assigning to b.N in benchmarks distorts the results

		"SA4000": true, //	Binary operator has identical expressions on both sides
		"SA4001": true, //	&*x gets simplified to x, it does not copy x
		"SA4003": true, //	Comparing unsigned values against negative values is pointless
		"SA4004": true, //	The loop exits unconditionally after one iteration
		"SA4005": true, //	Field assignment that will never be observed. Did you mean to use a pointer receiver?
		"SA4006": true, //	A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?
		"SA4008": true, //	The variable in the loop condition never changes, are you incrementing the wrong variable?
		"SA4009": true, //	A function argument is overwritten before its first use
		"SA4010": true, //	The result of append will never be observed anywhere
		"SA4011": true, //	Break statement with no effect. Did you mean to break out of an outer loop?
		"SA4012": true, //	Comparing a value against NaN even though no value is equal to NaN
		"SA4013": true, //	Negating a boolean twice (!!b) is the same as writing b. This is either redundant, or a typo.
		"SA4014": true, //	An if/else if chain has repeated conditions and no side-effects; if the condition didn’t match the first time, it won’t match the second time, either
		"SA4015": true, //	Calling functions like math.Ceil on floats converted from integers doesn’t do anything useful
		"SA4016": true, //	Certain bitwise operations, such as x ^ 0, do not do anything useful
		"SA4017": true, //	Discarding the return values of a function without side effects, making the call pointless
		"SA4018": true, //	Self-assignment of variables
		"SA4019": true, //	Multiple, identical build constraints in the same file
		"SA4020": true, //	Unreachable case clause in a type switch
		"SA4021": true, //	x = append(y) is equivalent to x = y
		"SA4022": true, //	Comparing the address of a variable against nil
		"SA4023": true, //	Impossible comparison of interface value with untyped nil
		"SA4024": true, //	Checking for impossible return value from a builtin function
		"SA4025": true, //	Integer division of literals that results in zero
		"SA4026": true, //	Go constants cannot express negative zero
		"SA4027": true, //	(*net/url.URL).Query returns a copy, modifying it doesn’t change the URL
		"SA4028": true, //	x % 1 is always zero
		"SA4029": true, //	Ineffective attempt at sorting slice
		"SA4030": true, //	Ineffective attempt at generating random number
		"SA4031": true, //	Checking never-nil value against nil
		"SA4032": true, //	Comparing runtime.GOOS or runtime.GOARCH against impossible value

		"SA5000": true, //	Assignment to nil map
		"SA5001": true, //	Deferring Close before checking for a possible error
		"SA5002": true, //	The empty for loop (for {}) spins and can block the scheduler
		"SA5003": true, //	Defers in infinite loops will never execute
		"SA5004": true, //	for { select { ... with an empty default branch spins
		"SA5005": true, //	The finalizer references the finalized object, preventing garbage collection
		"SA5007": true, //	Infinite recursive call
		"SA5008": true, //	Invalid struct tag
		"SA5009": true, //	Invalid Printf call
		"SA5010": true, //	Impossible type assertion
		"SA5011": true, //	Possible nil pointer dereference
		"SA5012": true, //	Passing odd-sized slice to function expecting even size

		"SA6000": true, //	Using regexp.Match or related in a loop, should use regexp.Compile
		"SA6001": true, //	Missing an optimization opportunity when indexing maps by byte slices
		"SA6002": true, //	Storing non-pointer values in sync.Pool allocates memory
		"SA6003": true, //	Converting a string to a slice of runes before ranging over it
		"SA6005": true, //	Inefficient string comparison with strings.ToLower or strings.ToUpper
		"SA6006": true, //	Using io.WriteString to write []byte

		"SA9001": true, //	Defers in range loops may not run when you expect them to
		"SA9002": true, //	Using a non-octal os.FileMode that looks like it was meant to be in octal.
		"SA9003": true, //	Empty body in an if or else branch
		"SA9004": true, //	Only the first constant has an explicit type
		"SA9005": true, //	Trying to marshal a struct with no public fields nor custom marshaling
		"SA9006": true, //	Dubious bit shifting of a fixed size integer value
		"SA9007": true, //	Deleting a directory that shouldn’t be deleted
		"SA9008": true, //	else branch of a type assertion is probably not reading the right value
		"SA9009": true, //	Ineffectual Go compiler directive
	}
	sChecks := map[string]bool{"S1000": true, //	Use plain channel send or receive instead of single-case select
		"S1001": true, //	Replace for loop with call to copy
		"S1002": true, //	Omit comparison with boolean constant
		"S1003": true, //	Replace call to strings.Index with strings.Contains
		"S1004": true, //	Replace call to bytes.Compare with bytes.Equal
		"S1005": true, //	Drop unnecessary use of the blank identifier
		"S1006": true, //	Use for { ... } for infinite loops
		"S1007": true, //	Simplify regular expression by using raw string literal
		"S1008": true, //	Simplify returning boolean expression
		"S1009": true, //	Omit redundant nil check on slices, maps, and channels
		"S1010": true, //	Omit default slice index
		"S1011": true, //	Use a single append to concatenate two slices
		"S1012": true, //	Replace time.Now().Sub(x) with time.Since(x)
		"S1016": true, //	Use a type conversion instead of manually copying struct fields
		"S1017": true, //	Replace manual trimming with strings.TrimPrefix
		"S1018": true, //	Use copy for sliding elements
		"S1019": true, //	Simplify make call by omitting redundant arguments
		"S1020": true, //	Omit redundant nil check in type assertion
		"S1021": true, //	Merge variable declaration and assignment
		"S1023": true, //	Omit redundant control flow
		"S1024": true, //	Replace x.Sub(time.Now()) with time.Until(x)
		"S1025": true, //	Don’t use fmt.Sprintf("%s", x) unnecessarily
		"S1028": true, //	Simplify error construction with fmt.Errorf
		"S1029": true, //	Range over the string directly
		"S1030": true, //	Use bytes.Buffer.String or bytes.Buffer.Bytes
		"S1031": true, //	Omit redundant nil check around loop
		"S1032": true, //	Use sort.Ints(x), sort.Float64s(x), and sort.Strings(x)
		"S1033": true, //	Unnecessary guard around call to delete
		"S1034": true, //	Use result of type assertion to simplify cases
		"S1035": true, //	Redundant call to net/http.CanonicalHeaderKey in method call on net/http.Header
		"S1036": true, //	Unnecessary guard around map access
		"S1037": true, //	Elaborate way of sleeping
		"S1038": true, //	Unnecessarily complex way of printing formatted string
		"S1039": true, //	Unnecessary use of fmt.Sprint
		"S1040": true, //	Type assertion to current type
	}

	stChecks := map[string]bool{"ST1000": true, //	Incorrect or missing package comment
		"ST1001": true, //	Dot imports are discouraged
		"ST1003": true, //	Poorly chosen identifier
		"ST1005": true, //	Incorrectly formatted error string
		"ST1006": true, //	Poorly chosen receiver name
		"ST1008": true, //	A function’s error value should be its last return value
		"ST1011": true, //	Poorly chosen name for variable of type time.Duration
		"ST1012": true, //	Poorly chosen name for error variable
		"ST1013": true, //	Should use constants for HTTP error codes, not magic numbers
		"ST1015": true, //	A switch’s default case should be the first or last case
		"ST1016": true, //	Use consistent method receiver names
		"ST1017": true, //	Don’t use Yoda conditions
		"ST1018": true, //	Avoid zero-width and control characters in string literals
		"ST1019": true, //	Importing the same package multiple times
		"ST1020": true, //	The documentation of an exported function should start with the function’s name
		"ST1021": true, //	The documentation of an exported type should start with type’s name
		"ST1022": true, //	The documentation of an exported variable or constant should start with variable’s name
		"ST1023": true, //	Redundant type in variable declaration
	}

	qfChecks := map[string]bool{
		"QF1001": true, //	Apply De Morgan’s law
		"QF1002": true, //	Convert untagged switch to tagged switch
		"QF1003": true, //	Convert if/else-if chain to tagged switch
		"QF1004": true, //	Use strings.ReplaceAll instead of strings.Replace with n == -1
		"QF1005": true, //	Expand call to math.Pow
		"QF1006": true, //	Lift if+break into loop condition
		"QF1007": true, //	Merge conditional assignment into variable declaration
		"QF1008": true, //	Omit embedded fields from selector expression
		"QF1009": true, //	Use time.Time.Equal instead of == operator
		"QF1010": true, //	Convert slice of bytes to string when printing it
		"QF1011": true, //	Omit redundant type from variable declaration
		"QF1012": true, //	Use fmt.Fprintf(x, ...) instead of x.Write(fmt.Sprintf(...))
	}

	stdAnalysers := []*analysis.Analyzer{
		appends.Analyzer,         // проверяет, что в append только одна переменная, т.е. не передается значение для добавления в слайс.
		asmdecl.Analyzer,         // проверяет, что файлы сборки соответствуют объявлениям Go.
		assign.Analyzer,          // проверяет бесполезные присваивания, например x = x
		atomic.Analyzer,          // проверяет распространенные ошибка использоsвания sync/atomic
		atomicalign.Analyzer,     // чтш-то с разрядностью, не понял
		bools.Analyzer,           // проверяет распространенные ошибки, связанные с использованием bool
		buildssa.Analyzer,        // ??
		buildtag.Analyzer,        // проверяет тэги сборки(buildtag)
		cgocall.Analyzer,         // проверяет нарушения правил передачи указателей cgo
		composite.Analyzer,       // ?? преаеряет наличие составных литералов без ключа
		copylock.Analyzer,        // проверяет блокировки, которые были установлены по-значению
		ctrlflow.Analyzer,        //??
		deepequalerrors.Analyzer, // проверяет использование reflect.deepEqual со значениями ошибок.
		defers.Analyzer,          // проверяет ошибочное использованеие defer
		directive.Analyzer,       // проверяет известные директивы инструментов
		errorsas.Analyzer,        // проверяте, что второй аргумент в errors.As - это указатель на тип, реализующий интерфейс ошибки
		// fieldalignment.Analyzer,      // обнаруживает структуры, которые использовали бы меньше памяти, если бы их поля были отсортированы
		findcall.Analyzer,            // ?
		framepointer.Analyzer,        // ?
		httpmux.Analyzer,             // ?
		httpresponse.Analyzer,        // проверяет ошибки в HTTP-ответах
		ifaceassert.Analyzer,         // находит невозожное приведение интерфейса-в-интерфейс, например из за одинаковых имен, но разных сигнатур приводимых интерфейсов
		inspect.Analyzer,             //?
		loopclosure.Analyzer,         // проверяет наличие ссылок на переменные цикла, входящие во вложенные функции.
		lostcancel.Analyzer,          // проверят выхов отмены контекста
		nilfunc.Analyzer,             // проверяет бесполезное сравнение функции с nil (fun == nil)
		nilness.Analyzer,             // проверьте на избыточные или невозможных сравнений с nil
		printf.Analyzer,              //проверяет согласованность строк и аргументов формата в printf
		reflectvaluecompare.Analyzer, //  проверяет случайное использование == или reflect.deepEqual для сравнения значений reflect.Value
		shadow.Analyzer,              // проверяет сокрытие пременных
		shift.Analyzer,               // проверяет наличие сдвигов, превышающих диапазон целого числа.
		sigchanyzer.Analyzer,         // проверяет неправильное использование небуферизованного сигнала в качестве аргумента для signal.Notify.
		sortslice.Analyzer,           // проверяет наличие вызовов sort.Slice, которые не используют тип slice в качестве первого аргумента.
		stdmethods.Analyzer,          // проверяет сигнатуры методов известных интерфейсов.
		stdversion.Analyzer,          // проверяет библиотечные символы, которые являются "слишком новыми" для действующей версии Go
		stringintconv.Analyzer,       // проверяет наличие преобразований string(int), что приведет к получению UTF-8 символа, а не преобразования числа в строку как есть
		structtag.Analyzer,           // проверяет правильность формирования тегов структурных полей.
		testinggoroutine.Analyzer,    // проверят вызовов Fatal из тестовой программы.
		tests.Analyzer,               // проверяет типичные ошибки при использовании тестов и примеров.
		timeformat.Analyzer,          // проверяет использование time.Format или time.Parse вызовы с неправильным форматом.
		unmarshal.Analyzer,           // проверят, что в unmarshal передаются значения отличные от указателя или интерфейса.
		unreachable.Analyzer,         // проверяет, что в коде есть недостижимые участки
		unsafeptr.Analyzer,           // проверяет, нет ли недопустимых преобразований uintptr в unsafe.Pointer
		unusedresult.Analyzer,        // проверяет неиспользуемые результаты вызовов некоторых функций.
		unusedwrite.Analyzer,         //  проверяет наличие неиспользуемых записей
		usesgenerics.Analyzer,        // проверяет, использует ли программа generic-и
	}

	externalAnalyzers := []*analysis.Analyzer{
		errcheck.Analyzer,
		ineffassign.Analyzer,
		wrapcheck.Analyzer,
	}

	var staticcheckers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if saChecks[v.Analyzer.Name] {
			staticcheckers = append(staticcheckers, v.Analyzer)
		}
		if sChecks[v.Analyzer.Name] {
			staticcheckers = append(staticcheckers, v.Analyzer)
		}
		if stChecks[v.Analyzer.Name] {
			staticcheckers = append(staticcheckers, v.Analyzer)
		}
		if qfChecks[v.Analyzer.Name] {
			staticcheckers = append(staticcheckers, v.Analyzer)
		}
	}
	staticcheckers = append(staticcheckers, stdAnalysers...)
	staticcheckers = append(staticcheckers, osexitcheck.Analyzer)
	staticcheckers = append(staticcheckers, externalAnalyzers...)

	multichecker.Main(
		staticcheckers...,
	)
}
