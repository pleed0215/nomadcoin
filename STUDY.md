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
