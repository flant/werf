---
title: Работа с электронной почтой
sidebar: applications-guide
permalink: documentation/guides/applications-guide/template/060-email.html
layout: guide
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- ____________
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с почтой.

Для того чтобы использовать почту мы предлагаем лишь один вариант - использовать внешнее API. В нашем примере это [mailgun](https://www.mailgun.com/).

Для того, чтобы ____________ приложение могло работать с mailgun необходимо установить и сконфигурировать зависимость и начать её использовать. Установим через `____________` зависимость:

```bash
____________
```

И [сконфигурируем согласно документации](____________) пакета

{% snippetcut name="____________" url="#" %}
{% raw %}
```____________
____________
____________
____________
```
{% endraw %}
{% endsnippetcut %}

В коде приложения подключение к API и отправка сообщения может выглядеть так:

{% snippetcut name="____________" url="#" %}
{% raw %}
```____________
____________
____________
____________
```
{% endraw %}
{% endsnippetcut %}

Для работы с mailgun необходимо пробросить в ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных](https://ru.werf.io/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, [когда мы делали базовое приложение](020-basic.html#secret-values-yaml)*

{% snippetcut name="secret-values.yaml (расшифрованный)" url="#" %}
{% raw %}
```yaml
app:
  ____________
  ____________
```
{% endraw %}
{% endsnippetcut %}

А не секретные значения — храним в `values.yaml`

{% snippetcut name="values.yaml" url="#" %}
{% raw %}
```yaml
  ____________
  ____________
  ____________
```
{% endraw %}
{% endsnippetcut %}

После того, как значения корректно прописаны и зашифрованы — мы можем пробросить соответствующие значения в Deployment.

{% snippetcut name="deployment.yaml" url="#" %}
{% raw %}
```yaml
        - name: ____________
          value: ____________
```
{% endraw %}
{% endsnippetcut %}

TODO: надо дать отсылку на какой-то гайд, где описано, как конкретно использовать ____________. Мало же просто его установить — надо ещё как-то юзать в коде.


<div>
    <a href="070-redis.html" class="nav-btn">Далее: Подключаем redis</a>
</div>