package staticlint

import (
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
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
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
		"SA": true,
	}

	var staticcheckers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if saChecks[v.Analyzer.Name] {
			staticcheckers = append(staticcheckers, v.Analyzer)
		}
	}

	stdAnalysers := []*analysis.Analyzer{
		appends.Analyzer,             // проверяет, что в append только одна переменная, т.е. не передается значение для добавления в слайс.
		asmdecl.Analyzer,             // проверяет, что файлы сборки соответствуют объявлениям Go.
		assign.Analyzer,              // проверяет бесполезные присваивания, например x = x
		atomic.Analyzer,              // проверяет распространенные ошибка использоsвания sync/atomic
		atomicalign.Analyzer,         // чтш-то с разрядностью, не понял
		bools.Analyzer,               // проверяет распространенные ошибки, связанные с использованием bool
		buildssa.Analyzer,            // ??
		buildtag.Analyzer,            // проверяет тэги сборки(buildtag)
		cgocall.Analyzer,             // проверяет нарушения правил передачи указателей cgo
		composite.Analyzer,           // ?? преаеряет наличие составных литералов без ключа
		copylock.Analyzer,            // проверяет блокировки, которые были установлены по-значению
		ctrlflow.Analyzer,            //??
		deepequalerrors.Analyzer,     // проверяет использование reflect.deepEqual со значениями ошибок.
		defers.Analyzer,              // проверяет ошибочное использованеие defer
		directive.Analyzer,           // проверяет известные директивы инструментов
		errorsas.Analyzer,            // проверяте, что второй аргумент в errors.As - это указатель на тип, реализующий интерфейс ошибки
		fieldalignment.Analyzer,      // обнаруживает структуры, которые использовали бы меньше памяти, если бы их поля были отсортированы
		findcall.Analyzer,            // ?
		framepointer.Analyzer,        // ?
		httpmux.Analyzer,             // ?
		httpresponse.Analyzer,        // проверяет ошибки в HTTP-ответах
		ifaceassert.Analyzer,         // находит невозожное приведение интерфейса-в-интерфейс, например из за одинаковых имен, но разных сигнатур приводимых интерфейсов
		inspect.Analyzer,             //?
		loopclosure.Analyzer,         // проверяет наличие ссылок на переменные цикла, входящие во вложенные функции.
		lostcancel.Analyzer,          // проверят выхов отмены контекста
		nilfunc.Analyzer,             // проверяет бесполезное сравнение функции с nil (fun == nil)
		nilfunc.Analyzer,             // проверьте на избыточные или невозможных сравнений с nil
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

	staticcheckers = append(staticcheckers, stdAnalysers...)
	multichecker.Main(
		staticcheckers...,
	)
}
