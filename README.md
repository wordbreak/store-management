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
docker compose up -f docker-compose.local.yaml

# 데이터베이스 마이그레이션
make goose env=local c=up
```

----

## Note

- API의 추후 변경을 고려하여 버저닝을 적용했습니다. 여기서 작성된 모든 API는 `v1` 버전을 가집니다.
- 현실 세계의 케이스를 고려하여 데이터베이스 Reader, Writer를 분리했고, 인메모리 캐시를 구현하여 사용했습니다.
- goose가 스키마의 변경사항을 추적하기 쉽다고 생각하여 데이터베이스 마이그레이션 툴로 사용했습니다.
- 쿼리가 복잡해질경우 ORM의 쿼리 빌더보다 직접 작성한 쿼리가 더 나은 성능을 보여준다 생각하였고, 그 생각의 연장으로 ORM보다는 struct mapper인 sqlx를 사용했습니다.
- Repository, Service, Controller로 레이어를 분리하여 비교적 의존성을 낮고, 테스트 작성이 쉬운 구조로 설계했습니다.
- 특히 검색쪽은 구현이 미흡한 상태로 제출하였습니다.
  - 상품 목록의 경우 커서 기반 페이지네이션이 적용되어 있지만 검색 결과에서는 페이지네이션이 되어 있지 않습니다.
  - 검색 결과의 정렬은 구현되어 있지 않습니다.
  - 특히 모든 검색이 LIKE 검색으로만 이루어져있어 검색 성능이 좋지 않습니다.
  - ElasitcSearch등의 옵션을 사용하여 검색을 구현하면 더 좋은 성능을 기대할 수 있을 것 같습니다. 특히 한글 초성 검색의 경우 간단한 토크나이저를 사용하면 쉽게 구현이 가능할 것으로 보입니다.
- 시간 관계상 테스트 코드는 거의 작성하지 못했습니다.
-
