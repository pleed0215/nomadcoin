```go
func plus(a ...int) int {
	var sum int = 0
	for _, item := range(a) {
		sum += item
	}
	return sum
}
```

## Package Sync

https://golang.org/pkg/sync/

- synchronous programming을 위한 패키지인듯하다.
- Once, WaitGroup,

## web site 만드는 방법

강의에서 다룬 내용은 모두 go의 네이티브 패키지를 이용한다.

1. net/http 패키지 사용
2. http.ListenAndServe를 이용해서 서버열기. 끗. ?

### HandleFunc 사용 방법

HandleFunc은 일종의 router를 다루는 것이라고 보면 될 것 같다.
대충 구조가 http.HandleFunc([Route], handlerFunc) 구조이다.

### HandlerFunc

위 HandleFunc의 두번째 인자로 들어가는 함수.
본격적인 request, response를 다룰 수 있게 해주는 핸들러.
템플릿을 다룰수도 있고, 데이터를 리턴해줄 수 있는 것으로 보인다.

`func (rw http.ResponseWriter, r *http.Request)`
위의 구조로 되어 있다.

### Fprint를 이용해서 화면에 표시할 글자 리턴해주기.

rw가 Writer 타입이기 때문에 Fprint를 이용하면 바로 응답을 해줄 수 있다.

### template을 이용하기

template 패키지를 이용해야 한다.
template.ParseFiles를 이용하면 템플릿을 렌더링할 수 있으며, 데이터를 템플릿에 넘길 수도 있다.
template의 Execute 함수를 이용하여 템플릿을 실행하도록 한다.

```go
	tmpl.Execute(rw, data)
```

data를 넘겨줄 때 조심해야 하는 것은, 저것역시 외부 파일로 나가야 하므로 data의 멤버들 역시 대문자로 써줘야 한다.

템플릿에서 넘겨준 데이터를 사용할 때에는 앞에 이중 중괄호를 사용하도록 한다.

```html
<title>{{.Title}}</title>
```

### template의 Must함수

솔직히 go에서는 에러처리가 좀.. 고급지지 못해서 그런지 짜증나는 부분이 있는데..
Must함수는 함수 안에서 미리 error를 핸들링 해준다.
Must 함수의 코드를 살펴 보면..

```go
	func Must(t *Template, err error) *Template {
		if err != nil {
			panic(err)
		}
		return t
	}
```

라구 되어 있다. 그래서 사용을 할 때에는..

```go
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
```

이렇게 사용해주면 된다.

### template 안에서 for 문 사용하려면?

참고 gotemplate 관련하여 vscode extension 설치를 추천한다.

```html
{{range .Blocks}}
<p class="text-lg bold">{{.Data}}</p>
<p class="text-md">{{.Hash}}</p>
<p class="text-md">{{.PrevHash}}</p>
{{end}}
```

### template 쪼개기.

django나 pug처럼 template을 쪼개서 다른 template에서 불러 올 수 있다.
{{ define }} 을 이용한다.
define으로 명명한 template은
{{template "이름"}} 이렇게 사용한다.

```html
{{ define "head" }}
<head>
  <meta charset="UTF-8" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <!--
    <link
      href="https://cdn.jsdelivr.net/npm/reset-css@5.0.1/reset.min.css"
      rel="stylesheet"
    />
    -->
  <link rel="stylesheet" href="https://unpkg.com/mvp.css" />
  <!--
    <link
      href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css"
      rel="stylesheet"
    />
    -->
  <title>{{.Title}} | Nomad Coin</title>
</head>
{{ end }}
```

header.html이 위와 같이 정의가 되어있다면...
실제로 사용할 때에는..

```html
<html lang="en">
  {{template "head"}} ...
</html>
```

위와 같이 사용한다.

template 이름들은 전부 소문자를 사용해야하는 모양이다. 대문자로 했더니 작동을 안했다.
template을 사용하기 전에 준비 과정이 필요하다.
template을 사전에 다 load를 해야 한다.
앞에서는 한 개의 단일 template은 ParseFiles를 이용해서 사용했지만.. template을 이용하기 위해서는
ParseGlob을 사용해야 한다.

전역으로 template의 포인터인 templates를 먼저 정의해서..

```go
	var templates *template.Template
```

templates를 다음과 같이 사용한다.

```go
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir+ "partials/*.gohtml"))
```

안타깝게도 go에서는 javascript처럼 `**/*.gohtml`과 같이 사용해 줄 수 없다...
그래서 위와 같이 정의해서 사용한다.
그리고 template을 rendering하는 것도 바뀐다.

```go
func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"GoGOGOOO", bc.GetBlockchain().AllBlock()}
	//tmpl.Execute(rw, data)
	templates.ExecuteTemplate(rw, "home", data)
}
```

이전과 비교를 위해서 Execute는 주석처리 했다.
Execute가 ExcuteTemplate으로 바뀌는데, 차이점은 template의 name(위의 gohtml에서 정의한..)이 추가 된 것이다.

!! 여기서 끝이 아니다.
data를 넘겨줬지만.. 실제로는 잘 안된다.
왜냐하면 각각 template에도 data를 넘겨주는 작업을 다시 해야 한다.

...
좀 짜증나는 부분이.... template에 이를테면.. .PageTitle을 넘겨줬다 치면.. 받은 template에서는 .PageTitle이라 쓰면.. 안된다.
그냥.. **.**만 써야 한다.

```go
{{template "head" .PageTitle}}
```

이를테면 위와같이 head template에 .PageTitle을 넘겨줬다치면..

```go
<title>{{.}} | Nomad Coin</title>
```

head 템플릿에서는 .으로 사용해야 한다...
이게 뭔 개소리야..
..
갑자기 go template이 sucks해 보인다...

### form in gohtml

form의 data를 읽어 오고 싶은가..?
설명은 코드로 대체하겠다.

```go
func add(rw http.ResponseWriter, r *http.Request) {
	data := addData{"Add"}
	switch(r.Method) {
		case "GET":
			templates.ExecuteTemplate(rw, "add", data)
		case "POST":
			r.ParseForm()
			data:=r.Form.Get("blockData")
			bc.GetBlockchain().AddBlock(data)
			http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
	//tmpl.Execute(rw, data)
}
```

ParseForm을 호출해줘야 form의 data를 읽어 올 수 있다. go의 코드가 공개되어 있으므로
코드를 참고해서 이해하도록 하자.

# REST API

## Marshaling?

Process of transforming the memory representation of an object into a data format subitable for storage or transmission
: object를 대표하는 메모리를 저장소나 전송에 맞는 포맷에 맞는 데이터 포맷으로 변경해주는 작업.

말이 어려움.

### go의 Marshal 함수

`func Marshal(v interface{}) ([]byte, error)`
v를 encoding한 JSON을 리턴해준다.

```go
	b, err := json.Marshal(data)
	utils.HandleError(err)
	fmt.Fprintf(rw, "%s", b)
```

javascript도 body에서 json으로 변경을 해줘서 그렇지.. 이런 작업들이 필요한 것 같다.
inspect의 network에 가서 위의 결과를 확인하면 안타깝게도 plain/text로 되어 있다. json이 아니라.
크롬에 json 뷰어 익스텐션을 설치를 했더니, 몰랐는디..

`rw.Header().Add("Content-Type", "application/json")`
위의 코드를 추가하여 Header json이라는 것을 알려주면 된다...

위 코드 블럭을 한줄로 줄일 수 있는 방법이 있다.

```go
	json.NewEncoder(rw).Encode(data)
```

## Struct field tags

go에서는 첫글자가 대문자가 아니면 export를 할 수가 없는데 사실 대부분의 json은 대부분 camel case기반이다. 그래서 사용하는 것이 Struct field tag.

```go
type URLDescription struct {
	URL string `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
}
```

뒤에 `\`json:"{name}"\`` 식으로 추가를 해주면 된다.

추가적으로, omitempty를 뒤에 ','로 덧붙여 주면 빈값일 경우에는 출력을 하지 않는다.

```go
type URLDescription struct {
	URL string `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}
```

Payload가 없다면 payload 필드는 json에 포함되지 않는다.

추가적으로 `json:"-"`하면 이 필드는 ignored가 된다.

## Marshal Text

위의 코드에서 URLDescription의 URL은 path를 포함하고 있긴하지만.. full path를 주고 싶다면 어떻게 해야할까?

```go
type Stringer interface {
	String() string
}
```

interface는 일종의 blueprint같은 역할을 한다고 한다.
Stringer interface의 String method를 이용하면 된다.
URLDescription struct의 String method를 만들어 주면된다.
사용은 fmt package가 Stringer interface를 사용하기 때문에, Println같은 fmt의 메소드를 사용하면 Stringer에서 정의한 String method를 사용할 수 있다.

```go
	type TextMarshaler interface {
		MarshalText() (text []byte, err error)
	}
```

Marshal 함수에서는 encode 패키지의 MarshalText를 이용하는 모양이다.
그래서 우리가 원하는 타입에 MarshalText 메소드를 주면 된다.

```go
func (u URL) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s", BASE_URL, u)), nil
}
```

그러면 marshal encode를 할 때의 data를 우리가 컨트롤 할 수 있게 된다.
항상 그렇듯 공식 문서를 참고해도 될 것 같다.
[TextMarshaler](https://golang.org/pkg/encoding/#TextMarshaler)

interface는 js 같이 꼭 implements를 명시할 필요는 없다.
왜그럴까..

### REST Client(VS code extension)

확장자 http인 것으로..

### POST

Post method를 받았을 때 data를 읽어 오려면...
json 패키지를 이용해서 decode를 해야 한다.

```go
var addBlockBody AddBlockBody
json.NewDecoder(r.Body).Decode(&addBlockBody)
rw.WriteHeader(http.StatusCreated)
```

status code를 주는 방법도 눈여겨 볼 것.

### NewServeMux

6.4 강의에서는 여러번 그랬듯이 패키지를 만드는 작업을 했는데,
문제는 rest api 서버와 explorer 웹서버를 동시에 돌리려고 Start를 하면...
문제가 생기다. 동기적인 작업들이다보니.. 앞에 작업에서 멈춰있게 된다.

동시에 돌리기 위해 `go 루틴`을 돌리면.. panic에러가 발생한다.
http가 같은 홈을 사용하기 때문이다. `HandleFunc("/")`
`panic: http: multiple registrations for /`

ListenAndServe 함수에 보면 두번째 변수인 handler를 여태까지 nil
로 했는데.. 이 부분을 이제 바꾸면 된다.

```go
func Start(aPort int) {
	newServer := http.NewServeMux()
	port=fmt.Sprintf(":%d", aPort)
	newServer.HandleFunc("/", documentation)
	newServer.HandleFunc("/blocks", blocks)
	fmt.Println("listening on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, newServer))
}
```

그래서 http대신 NewServerMux로 만들어줘서 newServer를 넘겨주면 오케이.

## Gorilla Mux

standard 라이브러리에서는 못하는?이 아니라 어려운 작업들이 있기 때문에 gorilla mux를 이용한다.
https://github.com/gorilla/mux

사용방법은 유사하나...

```go
r := mux.NewRouter()
r.HandleFunc("/products/{key}", ProductHandler)
r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
```

이런식으로 파라미터나 패턴들을 줄 수 있다.

`go get -u github.com/gorilla/mux`
명령어로 인스톨한다.

### NewRouter

`NewRouter()` 메소드로 새로운 mux 인스턴스를 만든다.
기존 사용방법이 똑같아서 굳이 뭘 바꿀 것이 없다.

```go
func block(rw http.ResponseWriter, r *http.Request) {
	if(r.Method == "GET") {
		vars := mux.Vars(r)
		fmt.Println(vars["id"])
	}
}

func Start(aPort int) {
	newServer := mux.NewRouter()
	port=fmt.Sprintf(":%d", aPort)
	newServer.HandleFunc("/", documentation).Methods("GET")
	newServer.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	newServer.HandleFunc("/block/{id:[0-9]+}", block).Methods("GET")
	fmt.Println("listening on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, newServer))
}
```

기존의 start 코드는 위와 같이 바뀔 수 있다.
`Vars`를 사용하여 파라미터도 가져온 것을 볼 수 있다.

## Height

블록채인 패키지로 다시가서.. height라는 것을 설명을 하시는데..
이건 id와 비슷한 것 같다. 별 것 없음..

## Error handling

GetBlock에도 에러 핸들링을 위해 에러 코드를 만들고..
rest api 패키지에도 에러 핸들링을 위한 struct를 만들자.

## Middleware

gorilla mux feature 중 하나인 것 같다.
https://github.com/gorilla/mux#middleware

니코는 adaptor를 설명하면서.. http.HandlerFunc 의 타입에 ServeHTTP를 정의해 놨기 때문에.. Handler를 이렇게 사용한것이 어썸하다 했다.
아직 이해는 잘 안된다.

# CLI

## 패키지 소개

### flag

flag parsing 패키지.
https://golang.org/pkg/flag/

### cobra

https://github.com/spf13/cobra

git, aduino 등 유명한 cli는 대부분 이 패키지를 이용했다.

command & application 만드는 걸 쉽게 해주고, intelligent suggestion도 해준다.

cobra에 설명을 많이 하고 싶지 않아서 여기서는 flag를 이용할 것이라 함.

## Parsing command

`go run main.go rest`
로 command를 입력하면 rest 같은 argument는 `os.Args`로 들어간다.

Println으로 출력해보면.. 이렇게 기괴하다..
`[/var/folders/qr/vgp27ww103xgsgh27gy0fr240000gn/T/go-build2826414399/b001/exe/main rest]`

앞에는 프로그램 정보, 뒤에부터가 argument인 것 같다.

그래서 기본적으로는 아래 코드와 같이 parsing을 진행할 수 있다.

```go
func main () {
	if len(os.Args) < 2  {
		usage()
	}

	switch os.Args[1] {
	case "explorer":
		fmt.Println("Start explorer")
	case "rest":
		fmt.Println("Start REST api")
	default:
		usage()
	}
}
```

## FlagSet

이제 flag를 진행할 차례.
`go run main.go explorer -port=4000`

Args flag를 parsing하기 위해서는 Args를 slice한다.
`os.Args[2:]`

새로운 flagset을 만드려면 flag패키지의 NewFlagSet함수를 이용하면 된다.
flagset이름과, errorHandling을 받는데, errorHandling은 세 종류의 정슈이다.
`ContinueOnError, ExitOnError, PanicOnError`

그래서, NewFlagSet은
`rest:=flag.NewFlagSet("rest", flag.ExitOnError)`

이런식으로 사용할 수가 있다.

위의 코드에서 rest는 flag패키지의 FlagSet struct인데

`flag:=rest.Int("port", 4000, "Defalt port: 4000")`

식으로, flag 이름, 기본값, 사용설명 순으로 인자를 받는다.

```go
case "rest":
		rest.Parse(flags)
		fmt.Println(*flag)
```

위의 코드처럼 Parsing이 진행되면 flag는 포인터 값이므로 \*flag로 해야 값을 읽어 올 수 있다.

FlagSet의 멤버 함수 중에는 `Parsed`가 있는데, 파싱이 잘되었는지 확인할 수가 있다.

## Flag

Flag파트는 따로 정리할 내용이 없다.
위의 flagset을 따로 떼서 flag로 .. 코드 보면 될 것 같다.

# Persistence

## Bolt

Bolt는 key/value 방식의 db. bolt는 완성 단계라, 더이상 수정도 없고, 안정적이고, 이해하기 쉽다구 한다.
헤로쿠도 쓴다고.

### Open

```go
func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	...
}
```

bolt 홈페이지에서 가져온 코드.
Open, Close 메소드를 이용하면 된다.
주의할점은 file lock을 걸기 때문에 동시에 db 작업을 여러프로세스에서는 할 수 없다. 작업이 다 끝나기 전에는 db가 열린 상태.
bolt.Open의 세 번째 인자가 옵션인데, 그래서 옵션에서는 db 교착상태에서 무한정 기다리는 것보다, 타입아웃 옵션을 줄 수 있는 것 같다.

```go
db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
```

### Transaction

한 번에 한개의 read-write 트랜잭션만 가능하지만, read-only 트랜잭션은 한 번에 많은 트랜잭션을 처리할 수 있다.
고루틴을 사용할 때 thread safe하지 않기 때문에 locking을 사용하여 고루틴으로 트랜잭션을 처리할 때에도 한 번만 가능하도록 해야 한다.

Read/Write transation은 `db.Update` 메소드, Read-only는 `db.View` 메소드를 이용한다.
각각의 db.Update는 디스크가 쓰기 작업을 commit하는 것을 기다리는데, `db.Batch`를 사용하여 여러 udpate를 조합하면 오버헤드를 줄일 수 있다.

여러개의 고루틴이 Batch를 호출하는 경우에만 유용하다.
