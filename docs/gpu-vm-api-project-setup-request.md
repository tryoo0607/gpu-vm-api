# gpu-vm-api 프로젝트 생성 요청서

> 작성일: 2026-08-14

## 목표

- 프로젝트 생성 후 AI 활용 가능하게 skill 배치까지 진행하는 것

## 프로젝트 경로

- `$PROJ_ROOT` : `~/projects/ai-mcmp`
- `$PROJ_ROOT_DOC` : `$PROJ_ROOT/docs`
- `$PROJ` : `$PROJ_ROOT/gpu-vm-api`
- `$PROJ_DOC` : `$PROJ/docs`

## 프로젝트 정보

- 프로젝트 생성 위치 : `$PROJ`
- 프로젝트 명 : `gpu-vm-api`
- lang : golang

## 진행 방식

1. project 생성
2. 프로젝트 내에 디렉토리 생성 docs 디렉토리 생성
3. 다음 파일 복사
   - `$PROJ_ROOT_DOC/AI-MCMP-Project-Skills.md` -> `$PROJ_DOC/AI-MCMP-Project-Skills.md`
4. git init 실행하기
5. 이 요청서를 `$PROJ_DOC`에 md로 생성해두기

## 수행 결과

| 단계 | 결과 |
| --- | --- |
| 1. 프로젝트 생성 | `~/projects/ai-mcmp/gpu-vm-api` 생성, module `github.com/tryoo0607/gpu-vm-api` (go 1.25.9) |
| 2. docs 디렉토리 생성 | `gpu-vm-api/docs` 생성 |
| 3. 스킬 문서 복사 | `docs/AI-MCMP-Project-Skills.md` 복사 완료 |
| 4. git init | `gpu-vm-api/.git` 초기화, 브랜치는 `main`으로 지정 |
| 5. 요청서 문서화 | 본 문서 (`docs/gpu-vm-api-project-setup-request.md`) |
| 추가. skill 배치 | `CLAUDE.md`에서 `@docs/AI-MCMP-Project-Skills.md` 참조 |
| 추가. `.gitignore` 생성 | 사용자 요청으로 Go 프로젝트용 `.gitignore` 추가 |

## 비고

- 원격 저장소는 개인 계정(https://github.com/tryoo0607/gpu-vm-api)에 생성했으며, Go module 경로는 `github.com/tryoo0607/gpu-vm-api`로 지정함.
