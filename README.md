# Asserlang 어쩔랭

Made with ♥️ in South Korea by [chul0721](https://github.com/chul0721) & [sujang958](https://github.com/sujang958)

이 프로젝트는 개발 중인 프로젝트입니다.  
[디스코드 서버에 참가하여 어쩔랭 개발에 기여하기](https://discord.gg/YBZd39qVfF)  
| 종류 | 경로 | 제작자 | 상태 |
|------|------|------|------|
| Node.JS 구현체 | [/Asserlang_nodeJS](https://github.com/assertive-lang/asserlang/tree/main/Asserlang_nodeJS) | chul0721, sujang958 | v1 |
| C# 구현체 | [/Asserlang_CSharp](https://github.com/assertive-lang/asserlang/tree/main/Asserlang_CSharp) | c3nb | v1 |
| Go 구현체 | [/Asserlang_Go](https://github.com/assertive-lang/asserlang/tree/main/Asserlang_Go) |  | on process |
| AsserFuck Rust 구현체 | [/extras/AsserFuck_Rust](https://github.com/assertive-lang/asserlang/tree/main/Asserlang_Go) | sujang958 | v1 |

유행어를 본따 만든 [엄랭](https://github.com/rycont/umjunsik-lang), [몰랭](https://github.com/ArpaAP/mollang), 그리고 [슈숙 언어](https://github.com/yf-dev/syusuk)와 같은 언어들에 영감을 받아 만들게 되었습니다.

코드가 다소 이상하다고 느껴지신다면 언제든 PR로 리펙토링 해주시면 감사하겠습니다. 💻  
그 외에도 PR은 환영입니다! 🙋

> 이 문서는 표준 구현체인 Node.js 기준으로 작성되었습니다.

## 도움을 주신 분들 ✨
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/c3nb"><img src="https://avatars.githubusercontent.com/u/73321185?v=4?s=100" width="100px;" alt=""/><br /><sub><b>C#Newbie</b></sub></a><br /><a href="https://github.com/assertive-lang/asserlang/commits?author=c3nb" title="Platform">📦</a></td>
    <td align="center"><a href="https://github.com/wjdqhry"><img src="https://avatars.githubusercontent.com/u/30039641?v=4?s=100" width="100px;" alt=""/><br /><sub><b>정보교</b></sub></a><br /><a href="https://github.com/assertive-lang/asserlang/commits?author=wjdqhry" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

# 문법

무조건 코드의 시작과 끝에는 각각 '쿠쿠루삥뽕'와 '슉슈슉슉'가 포함되어야 합니다.  
키워드는 변수의 이름이 될 수 없습니다.  
파일 확장자는 .astv를 사용합니다.  
띄어쓰기 대신 `~`를 사용합니다.  
줄바꿈을 통해 코드를 인식합니다.

## 연산자

```
ㅋ: +1를 의미합니다.
ㅎ: -1를 의미합니다.
ㅌ: 곱하기를 의미합니다.
```

> 예) ㅋㅋㅋㅋㅋㅌㅋㅋㅌㅋㅋㅋㅋ = 4 x 2 x 4
>
> 예) ㅋㅋㅎㅌㅋ = 1 x 1

## 변수

## 수를 담는 변수

#### 선언

```
쿠쿠루삥뽕
어쩔냉장고~ㅋㅋ
슉슈슉슉
```

> 변수 "냉장고"을 선언과 동시에 2라는 값으로 초기화 합니다.
>
> - 키워드는 변수 이름이 될 수 없습니다. (연산자 또한 키워드에 포함됩니다)
>   - 잘못된 예) 어쩔어쩔~ㅋㅋ
>   - 잘못된 예) 어쩔ㅋㅋ~ㅋㅋ
> - 변수 선언 시 초기화를 하지 않을 경우 0이 할당됩니다.
>   - 예) 어쩔초고속진공블랜딩믹서기

#### 할당

```
쿠쿠루삥뽕
어쩔냉장고~ㅋㅋㅋ
저쩔냉장고~ㅋㅋ
슉슈슉슉
```

> 변수 "냉장고"을 선언하며 동시에 3이라는 값으로 초기화 합니다.
>
> 변수 "냉장고"에 2라는 값을 할당 해 줍니다.

## 유니코드를 담는 변수

#### 선언

```
쿠쿠루삥뽕
우짤래미냉장고~ㅋㅋㅋㅋㅋㅋㅋㅌㅋㅋㅋㅋㅋㅋㅋㅋ
ㅇㅉ냉장고
슉슈슉슉
```

> 변수 "냉장고"을 선언과 동시에 "C"라는 값으로 초기화 합니다.
>
> - 변수 선언 시 초기화를 하지 않을 경우 0에 해당하는 유니코드 값이 할당됩니다.
>   - 예) 어쩔초고속진공블랜딩믹서기

#### 할당

```
쿠쿠루삥뽕
우짤래미냉장고~ㅋㅋㅋㅋㅋㅋㅋㅌㅋㅋㅋㅋㅋㅋㅋㅋ
저짤래미냉장고~ㅋㅋㅋㅋㅋㅋㅋㅌㅋㅋㅋㅋㅋ
ㅇㅉ냉장고
슉슈슉슉
```

> 변수 "냉장고"을 선언하며 동시에 "C"라는 값으로 초기화 합니다.
>
> 변수 "냉장고"에 "A"라는 값을 할당 해 줍니다.

## 입출력

#### 입력

```
쿠쿠루삥뽕
ㅌㅂ
슉슈슉슉
```

> 사용자에게 입력을 받습니다.

```
쿠쿠루삥뽕
어쩔다이슨v15디렉트앱솔루트엑스트라청소기~ㅌㅂ
슉슈슉슉
```

> 사용자에게 입력을 받은 후 변수 "다이슨v15디렉트앱솔루트엑스트라청소기"에 해당 값을 저장합니다.

#### 출력

```
쿠쿠루삥뽕
어쩔냉장고~ㅋㅋㅋㅋㅋ
ㅇㅉ냉장고
슉슈슉슉
```

> 결과: 5
>
> 냉장고 변수를 출력합니다.

```
쿠쿠루삥뽕
ㅇㅉㅌㅂ
슉슈슉슉
```

> 사용자에게 입력 받은 후 해당 값을 출력합니다.

## 함수

```
쿠쿠루삥뽕
안물수고염~킹받죠~빡쳤죠
어쩔냉장고~ㅋ
ㅇㅉ냉장고
ㅇㅉ킹받죠
ㅇㅉ빡쳤죠
안물
안궁수고염~ㅋㅋㅋ~ㅋㅋㅋㅋㅋ
슉슈슉슉
```

> 결과: 1 3 5
>
> 안물 키워드로 함수를 선언하고 안궁 키워드로 함수를 사용합니다.
>
> 안물{함수명}~{인자1}~{인자2}~... ... 안물
>
> 함수 선언시의 블록 구분은 안물 키워드를 시작과 끝에 둠으로써 구분합니다.

### Return

```
쿠쿠루삥뽕

안물반환~와샌즈
무지개반사와샌즈
안물

ㅇㅉ안궁반환~ㅋㅋㅌㅋㅋ
슉슈슉슉
```

> 결과: 22
>
> 무지개반사{Return할 값}

## 조건문

```
쿠쿠루삥뽕
어쩔개~ㅋㅎ
어쩔냉장고~ㅋㅋ
ㅇㅉ냉장고
화났쥬?개킹받쥬?저쩔냉장고~ㅋ
ㅇㅉ냉장고
슉슈슉슉
```

> 결과: 2 1
>
> 화났쥬?(조건)킹받쥬?(조건이 0일 때 실행할 코드)

## 실행

#### Node.JS 구현체를 이용하여 실행

터미널 및 콘솔에 아래 코드를 순서대로 입력하세요.
최신 버전의 git과 node.js가 설치되어 있어야 합니다.

```
$ git clone http://github.com/assertive-lang/asserlang
$ cd asserlang
$ npm i
$ tsc
$ node dist/index.js 파일명
```

## 에러

- 아무것도 모르죠?
  - 시작과 끝에 "쿠쿠루삥뽕"와 "슉슈슉슉"를 포함하지 않은 경우 발생하는 에러
- 어쩔변수
  - 변수에 관련된 구문에서 발생한 에러
- 어쩔파일
  - 파일을 제대로 불러오지 못한 경우 발생하는 에러
- 안물
  - 함수의 선언 과정에서 생긴 에러
- 어쩔조건
  - 조건문 코드에서 생긴 에러

🥕

~~슉슈슉슉~~
