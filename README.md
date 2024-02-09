## 개발 환경 구성


### pre-commit

로컬 환경에서 [pre-commit](https://pre-commit.com/)을 사용하여 린트를 수행합니다. `pre-commit`을 먼저 설치해주세요.

이후 `pre-commit install-hooks` 명령어를 실행하여 린트에 필요한 패키지를 설치합니다.

추가로 `goimports`를 설치해줍니다.

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

### Database migration

데이터베이스 마이그레이션 툴로 [goose](https://github.com/pressly/goose)를 사용합니다. `goose`를 먼저 설치해주세요.
