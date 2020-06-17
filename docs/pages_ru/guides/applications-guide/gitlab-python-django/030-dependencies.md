---
title: Подключение зависимостей
sidebar: applications-guide
permalink: documentation/guides/applications-guide/template/030-dependencies.html
layout: guide
---

{% filesused title="Файлы, упомянутые в главе" %}
- werf.yaml
{% endfilesused %}

Werf подразумевает, что лучшей практикой будет разделить сборочный процесс на этапы, каждый с четкими функциями и своим назначением. Каждый такой этап соответствует промежуточному образу, подобно слоям в Docker. В werf такой этап называется стадией, и конечный образ в итоге состоит из набора собранных стадий. Все стадии хранятся в хранилище стадий, которое можно рассматривать как кэш сборки приложения, хотя по сути это скорее часть контекста сборки.

Стадии — это этапы сборочного процесса, кирпичи, из которых в итоге собирается конечный образ. Стадия собирается из группы сборочных инструкций, указанных в конфигурации. Причем группировка этих инструкций не случайна, имеет определенную логику и учитывает условия и правила сборки. С каждой стадией связан конкретный Docker-образ. Подробнее о том, какие стадии для чего предполагаются можно посмотреть в [документации](https://ru.werf.io/documentation/reference/stages_and_images.html).

Werf предлагает использовать для стадий следующую стратегию:

*   использовать стадию beforeInstall для инсталляции системных пакетов;
*   использовать стадию install для инсталляции системных зависимостей и зависимостей приложения;
*   использовать стадию beforeSetup для настройки системных параметров и установки приложения;
*   использовать стадию setup для настройки приложения.

Подробно про стадии описано в [документации](https://ru.werf.io/documentation/configuration/stapel_image/assembly_instructions.html).

Одно из основных преимуществ использования стадий в том, что мы можем не перезапускать нашу сборку с нуля, а перезапускать её только с той стадии, которая зависит от изменений в определенных файлах.

В нашем случае в качестве примера мы можем взять файл `requirements.txt`.

Те кто уже сталкивался с разработкой на python приложений знают, что в файле `requirements.txt` указываются зависимости которые нужны для сборки приложения. Потому самое логичное указать данный файл в зависимости сборки, чтобы в случае изменений в нём, была перезапущена сборка только со стадии **_install_**.

Для этого в одной из первых глав мы сразу и добавляли наш файл requirements.txt в зависимости werf.
В Django в качестве менеджера зависимостей используется pip. Пропишем его использование в файле `werf.yaml` и затем оптимизируем его использование.

## Подключение менеджера зависимостей

Пропишем команды `pip install` в нужные стадии сборки в `werf.yaml`

```yaml
  install:
  - name: Install python requirements
    pip:
      requirements: /usr/src/app/requirements.txt
      executable: pip3.6

```

Однако, если оставить всё так — стадия `install` не будет запускаться при изменении файла со списком пакетов. Подобная зависимость пользовательской стадии от изменений [указывается с помощью параметра git.stageDependencies](https://ru.werf.io/documentation/configuration/stapel_image/assembly_instructions.html#%D0%B7%D0%B0%D0%B2%D0%B8%D1%81%D0%B8%D0%BC%D0%BE%D1%81%D1%82%D1%8C-%D0%BE%D1%82-%D0%B8%D0%B7%D0%BC%D0%B5%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9-%D0%B2-git-%D1%80%D0%B5%D0%BF%D0%BE%D0%B7%D0%B8%D1%82%D0%BE%D1%80%D0%B8%D0%B8):

```yaml
git:
- add: /
  to: /app
  stageDependencies:
    install:
    - requirements.txt
```

При изменении файла `requirements.txt` стадия `install` будет запущена заново.

## Оптимизация сборки

У pip есть кеш, чтобы каждый раз менеджер зависимостей не скачивал заново один и тот-же пакет.
Находится он в папке: `~/.cache/pip/`

Для того, чтобы оптимизировать работу с этим кешом при сборке, мы добавим специальную конструкцию в werf.yaml:

```yaml
mount:
- from: build_dir
  to: /app/.cache/pip
```

При каждом запуске билда, эта директория будет мантироваться с сервера, где запускается билд, и  не будет очищаться между билдами.
Так между сборками, у нас сохранится этот кеш.
