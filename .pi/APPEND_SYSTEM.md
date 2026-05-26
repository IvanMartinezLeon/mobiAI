# Detección automática del framework del proyecto
Al iniciar la conversación o trabajar en el código, ejecuta los siguientes comandos para identificar el framework/tecnología principal:

## Frameworks comunes por archivo de configuración
ls -la package.json pubspec.yaml pom.xml build.gradle build.gradle.kts Cargo.toml mix.exs project.clj go.mod composer.json pyproject.toml requirements.txt setup.py setup.cfg CMakeLists.txt Makefile yarn.lock package-lock.json pnpm-lock.yaml Deno.json deno.jsonc 2>/dev/null

## Flutter/Dart
test -f pubspec.yaml && head -30 pubspec.yaml

## Node.js
test -f package.json && head -30 package.json

## Python
test -f pyproject.toml && head -30 pyproject.toml
test -f requirements.txt && head -30 requirements.txt

## Rust
test -f Cargo.toml && head -30 Cargo.toml

## Go
test -f go.mod && head -30 go.mod

## Java/Kotlin (Maven)
test -f pom.xml && head -50 pom.xml

## Java/Kotlin (Gradle)
test -f build.gradle && head -30 build.gradle
test -f build.gradle.kts && head -30 build.gradle.kts

## Elixir
test -f mix.exs && head -30 mix.exs

## PHP
test -f composer.json && head -30 composer.json

## .NET
test -f *.csproj && ls *.csproj 2>/dev/null

## Docker
test -f Dockerfile && head -20 Dockerfile
test -f docker-compose.yml && head -20 docker-compose.yml

## README
test -f README.md && head -20 README.md

Con base en los resultados, determina el framework/tecnología principal y menciónalo explícitamente al inicio de la conversación.