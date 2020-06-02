---
title: Гайд по использованию Ruby On Rails + GitLab + Werf
sidebar: applications-guide
permalink: documentation/guides/applications-guide/gitlab-rails/040-assets.html
author: alexey.chazov <alexey.chazov@flant.com>
layout: guide
toc: false
author_team: "bravo"
author_name: "alexey.chazov"
ci: "gitlab"
language: "ruby"
framework: "rails"
is_compiled: 0
package_managers_possible:
 - bundler
package_managers_chosen: "bundler"
unit_tests_possible:
 - Rspec
unit_tests_chosen: "Rspec"
assets_generator_possible:
 - webpack
 - gulp
assets_generator_chosen: "webpack"
---


# Генерируем и раздаем ассеты

В какой-то момент в процессе разработки вам понадобятся ассеты (т.е. картинки, css, js).
Asset Pipeline представляет фреймворк для соединения и минимизации или сжатия ассетов JavaScript и CSS. Он также добавляет возможность писать эти ассеты на других языках и препроцессорах, таких как CoffeeScript, Sass и ERB. Это позволяет автоматически комбинировать ассеты приложения с ассетами других гемов.

Для генерации ассетов мы будем использовать команду `bundle exec rake assets:precompile`.

Интуитивно понятно, что на стадии сборки нам надо будет вызвать скрипт, который генерирует файлы, т.е. что-то надо будет дописать в `werf.yaml`. Однако, не только там — ведь какое-то приложение в production должно непосредственно отдавать статические файлы. Мы не будем отдавать файлики с помощью Rails. Хочется, чтобы статику раздавал nginx. А значит надо будет внести какие-то изменения и в helm чарт.

## Сценарий сборки ассетов

Команда `assets:precompile` по умолчанию Rails для `production` прекомпилирует файлы в директорию `public/assets`

Тут есть один нюанс - при сборке приложения мы не рекомендуем использовать какие-либо изменяемые переменные. Потому что собранный бинарный образ должен быть независимым от конкретного окружения. А значит во время сборки у нас не может быть, например, базы данных, user-generated контента и подобных вещей.

По непонятной причине - для генерации assets rails ходит в базу данных, хотя не понятно для каких целей и для этого - нужен SECRET_KEY_BASE​. При текущей сборке - мы использовали workaround, передав фейковое значение. По этому поводу есть issue созданное более 2х лет назад, но в версии rails 2.7 - до сих пор так. Если вы знаете, зачем авторы Rails так сделали - просьба написать в комментариях.

## Какие изменения необходимо внести

Генерация ассетов происходит в артефакте на стадии `setup`, так как данная стадия рекомендуется для настройки приложения

Для уменьшения нагрузки на процесс основного приложения которое обрабатыаем логику работы rails приложения мы будем отдавать статические файлы через `nginx`
Мы запустим оба контейнера одним деплойментом и все запросы будет приходить вначале на nginx и если в запросе не будет отдача статических файлов - запрос будет отправлен прмложению.

### Изменения в сборке

Добавим стадию сборки ассетов для приложения в файл `werf.yaml`

[werf.yaml](gitlab-rails-files/examples/example_2/werf.yaml#L45)
```yaml
  setup:
  - name: build assets
    shell: RAILS_ENV=production SECRET_KEY_BASE=fake bundle exec rake assets:precompile
    args:
      chdir: /app
```

Окей, а в каком контейнере в конечном итоге должны оказаться собранные файлы? Есть минимум два варианта:

*   Делать один образ в котором: рельсы, сгенерированные ассеты, нгинкс. Запускать этот один и тот же образ двумя разными способами (с разным исполняемым файлом)
*   Делать два образа: рельсы отдельно, nginx + сгенерированные ассеты отдельно.

В первом варианте при каждом изменении будут перекатываться оба контейнера. Такое себе в большинстве случаев.

Пойдём вторым путём.

Дальше сложности, ибо для сборки нужны рельсы, нода, но в финальном образе не хочется иметь вот этого всего дерьма: нам в финальном образе нужна только статика и нгинкс.

И ВОТ ТУТ нам на помощь приходят артефакты. Мы ВОТ ТУТ объясняем что такое артефакты и поясняем, что мы сможем сгенерить в одном а пихнуть в другое и финальный образ будет вжух шустрый быстрый маленький охуеннный.


И рассказываем как конкретно будем собирать

В образе с нашим приложением мы не хотим чтобы у нас была установлена среда для сборки приложения и nginx а также для того чтобы уменьшить размеры образов - мы воспользуемся сборкой с помощью артефактов.

[Артефакт](https://ru.werf.io/documentation/configuration/stapel_artifact.html) — это специальный образ, используемый в других артефактах или отдельных образах, описанных в конфигурации. Артефакт предназначен преимущественно для отделения ресурсов инструментов сборки от процесса сборки образа приложения. Примерами таких ресурсов могут быть — программное обеспечение или данные, которые необходимы для сборки, но не нужны для запуска приложения, и т.п.

С помощью такого подхода мы сможем собрать и подготовить все файлы и зависимости в одном образе и импортировать нужные нам файлы по двум разным docker контейнерам, где в одном - будет среда для выполнения приложения ruby on rails а во втором - только ngin со статическими файлами.

Разница `artifact` и `image` так же состоит в том - что `artifact` нельза запустиль локально для дебага как `image` командой `werf run`

Импорт указывается отдельной директивой следующим в следующем синтаксисе:

[werf.yaml](gitlab-rails-files/examples/example_2/werf.yaml#L19)
```yaml
artifact: build
from: ruby:2.7.1
// build stages
---
image: rails
from: ruby:2.7.1-slim
import:
- artifact: build
  add: /usr/local/bundle
  after: install
- artifact: build
  add: /app
  after: install
// build stages
---
---
image: assets
from: nginx:alpine
ansible:
  beforeInstall:
  - name: Add nginx config
    copy:
      content: |
{{ .Files.Get ".werf/nginx.conf" | indent 8 }}
      dest: /etc/nginx/nginx.conf
import:
- artifact: build
  add: /app/public
  to: /www
  after: setup
```

Подготовленные ассеты мы будет отдавать через отдельный nginx контейнер в поде чтобы не загружать основное приложение лишними подключениями. Для этого так-же производится импорт подготовленных файлов в отдельный образ.
При работе мы планируем, что все запросы будут проксироваться через nginx, поэтому заменяем файл `/etc/nginx/nginx.conf` на необходимый нам, который находится также в репозитории с приложением. Такой подход позволит нам управлять лимитом подключений который может принять приложение.

### Изменения в деплое

При таком подходе изменим деплой нашего приложения добавив еще один контейнер в наш деплоймент с приложением.  Укажем livenessProbe и readinessProbe, которые будут проверять корректную работу контейнера в поде. preStop команда необходима для корректного завершение процесса nginx. В таком случае при новом выкате новой версии приложения будет корректное завершение всех активных сессий.

[.helm/templates/deployment.yaml](gitlab-rails-files/examples/example_2/.helm/templates/deployment.yaml#L33)
```yaml
      - name: assets
{{ tuple "assets" . | include "werf_container_image" | indent 8 }}
        lifecycle:
          preStop:
            exec:
              command: ["/usr/sbin/nginx", "-s", "quit"]
        livenessProbe:
          httpGet:
            path: /healthz
            port: 80
            scheme: HTTP
        readinessProbe:
          httpGet:
            path: /healthz
            port: 80
            scheme: HTTP
        ports:
        - containerPort: 80
          name: http
          protocol: TCP
```

В описании сервиса - так же должен быть указан правильный порт

[.helm/templates/service.yaml](gitlab-rails-files/examples/example_2/.helm/templates/service.yaml#L9)
```yaml
  ports:
  - name: http
    port: 80
    protocol: TCP
```


### Изменения в роутинге

Поскольку у нас маршрутизация запросов происходит черех nginx контейнер а не на основе ingress ресурсов - нам необходимо только указать коректный порт для сервиса

```yaml
      paths:
      - path: /
        backend:
          serviceName: {{ .Chart.Name }}
          servicePort: 80
```

Если мы хотим разделять трафик на уровне ingress - нужно разделить запросы по path и портам

[.helm/templates/ingress.yaml](gitlab-rails-files/examples/example_2/.helm/templates/ingress.yaml#L9)
```yaml
      paths:
      - path: /
        backend:
          serviceName: {{ .Chart.Name }}
          servicePort: 3000
      - path: /assets
        backend:
          serviceName: {{ .Chart.Name }}
          servicePort: 80
```