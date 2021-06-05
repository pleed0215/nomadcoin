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
