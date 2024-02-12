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

Makefile에 마이그레이션을 위한 명령어가 정의되어 있습니다.

```bash
# 새 마이그레이션 생성
make goose c=create name=마이그레이션 이름

# 기타 goose 명령어
# make goose env={local|dev|test} c=(goose 명령어)

# 예시
## 마이그레이션 적용
make goose env=local c=up
## 마이그레이션 롤백
make goose env=local c=down
```

## 서버 실행

앞 문단의 로컬 환경 구성을 끝낸 뒤 아래 흐름을 통해 서버를 실행합니다.

```bash
# 서버 및 데이터베이스를 docker로 실행
docker compose up -f docker-compose.local.yml

# 데이터베이스 마이그레이션
make goose env=local c=up
```
