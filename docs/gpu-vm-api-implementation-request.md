# gpu-vm-api 구현 요청서

> 작성일: 2026-08-14 · 범위: **1단계 (구현까지)**

## 필수 확인

- **public repo** 다. secret 관리가 필수적임
- 접속 정보는 **전부 환경변수**, 문서엔 `{BASE}` 자리표시자로만 표기
- 호스트·포트·계정·크리덴셜을 코드·문서·커밋 메시지에 **일체 남기지 않는다**

## 작업 목표

- go echo 서버 기반으로 Ubuntu 기반 GPU VM 생성 및 정보 제공 테스트
- Swagger 기반 API 구성
- **단계별로 구현 계획 먼저 세운다**

### 완료 판정 — 2단계로 나눈다

접속 정보가 아직 제공되지 않았으므로 이번 요청은 1단계까지다.

| 단계 | 범위 | 판정 |
|---|---|---|
| **1단계 (이번)** | 서버 구현 · Swagger · **정리(삭제) 로직 포함** | 빌드 통과 + Swagger 생성 + dry-run 경로 검증 |
| 2단계 (추후) | 접속 정보 수령 후 실기동 | 실제로 VM 이 뜨고 API 로 정보가 나옴 |

## 목표 아키텍처

- Go 서버는 **CB-Tumblebug REST 를 호출**한다
- **AWS SDK 직접 호출 아님**

## 프로젝트 경로

| 변수 | 경로 |
|---|---|
| `$PROJ_ROOT` | `~/projects/ai-mcmp` |
| `$PROJ` | `$PROJ_ROOT/gpu-vm-api` |
| `$PROJ_DOC` | `$PROJ/docs` |
| `$CBTB` | `~/projects/m-cmp/cb-tumblebug` (로컬 clone) |

## 분석 요청

작업 전 CB-Tumblebug 을 파악할 것. 다만 **읽는 순서와 정본을 지킨다.**

1. **사용법 정본** — `~/projects/_docs/m-cmp/design/cb-tumblebug-usage.md` (511줄)
   스펙 찾기 → 이미지 → 조합 검증 → 생성 → 상태/접속 → 제어/정리까지 정리돼 있다. **여기부터 읽는다**
2. **API 정본** — 서버 swagger `{BASE}/api/doc.json`
3. `$CBTB` 로컬 clone — **구조 파악용**

### ⛔ 로컬 clone 을 스펙 근거로 삼지 마라

clone HEAD 는 `v0.12.30-21-g6d9f1040` — **태그 `v0.12.30` 보다 21커밋 앞서 있고 breaking 변경이 섞여 있다**(`feat(rdbms)!:`). 서버는 0.12.30 이다.

clone 에서 스펙을 확인해야 하면 **태그 `v0.12.30` 기준**으로 보거나 서버 `doc.json` 을 쓴다.

### ⛔ 옛 용어·경로를 쓰지 마라

v0.12.6 에서 전면 개명됐다. 학습된 옛 API 를 그대로 쓰면 **전부 404** 다.

| 옛것 | 지금 |
|---|---|
| MCI | **Infra** (`/infra`) |
| VM | **Node** |
| SubGroup | **NodeGroup** |
| `commonSpec` / `commonImage` | `specId` / `imageId` |
| `subGroupSize` (문자열) | `nodeGroupSize` (정수) |

## VM 생성 시 필요한 정보

| 항목 | 값 |
|---|---|
| Tumblebug | 0.12.30 |
| CSP / 리전 | AWS **`us-west-2`** |
| 템플릿 | `infra-aws-gpu-simple` (system NS) |
| 이미지 | 템플릿이 들고 있음 — Ubuntu 22.04 DL Base OSS Nvidia Driver |

접속 가능한 서버에 m-cmp 스택 전체가 떠 있다. **접근 정보는 2단계에서 제공**한다.

### 진행 방식 — 템플릿 기반

가이드가 CB-Tumblebug 의 **Template 기능**으로 VM 을 만든다. 이 프로젝트는 **그 방식을 API 로 실행**하는 것이 목적이다.

### 🔑 생성 전 사전검증은 필수다

GPU 는 zone 별 재고 부족이 흔하다.

- `POST {BASE}/specImagePairReview` 를 먼저 치고 응답의 **`suggestedZone`** 을 `nodeGroups[].zone` 에 넣는다
- 통과 결과는 재고에 따라 변하는 **런타임 값**이다. 값을 코드에 고정하지 말고 **매번 조회**한다

### 🔑 비용 — 정리까지가 구현이다

- 템플릿에 GPU NodeGroup 이 **2개** 들어 있다 (L4 $2.0144/h + L40S $2.2421/h). 그대로 띄우면 **시간당 $4.26**
  → 안 쓸 NodeGroup 은 제거하고 만든다
- **삭제에 캐스케이드가 없다.** `Delete Infra` → 공유 리소스(`sharedResources`) → NS 순서
- `option=force` 는 **쓰지 마라** — CSP 종료 확인 없이 기록만 지워서 **과금되는 고아 인스턴스**가 남는다
- 정리 API 를 **1단계 구현 범위에 포함**한다. 실기동 전에 회수 경로가 먼저 있어야 한다

### 그 밖의 함정

- **제어는 GET** 이다 (`/control/infra/{id}?action=...`). PUT 아님
- **nsId 가 갈린다** — spec/image 계열은 `system`, infra 계열은 `default`
- `x-credential-holder` 헤더가 추천 결과를 필터링한다. GPU 가 안 보이면 여기부터 의심
- **rate limit** 글로벌 50 req/s, Infra 조회 계열은 더 조임 → **폴링 간격 주의**

## 작업 이력 작성

다음을 남긴다.

- 작업 시 어떤 내용을 진행했는지
- 어떤 문제가 있었고 어떻게 해결했는지
- 결정이 필요했던 사항과 무엇으로 정했는지
- 추가로 제시한 대안들

**기록 위치는 `_local-docs`** (= `~/projects/_docs/ai-mcmp/`). 기존 kit 체계를 그대로 따른다. 파일 하나에 몰아 쓸 필요 없음.

### LLM 2종 증적

과제 요구사항이다 — Claude + 알파(ChatGPT / Gemini 등).

**어느 부분을 어느 LLM 이 만들었는지 작업하면서 같이 기록한다.** 사후 복원은 불가능하다.

## 주의 사항

- **작업 선행 금지.** 원칙은 **계획 → 확인 → 작업**
- 위험한 내용은 반드시 사전에 확인받기
- 애매하면 반드시 확인 요청하기
