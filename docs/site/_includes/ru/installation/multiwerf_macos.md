```shell
# добавим директорию ~/bin в PATH
export PATH=$PATH:$HOME/bin
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc

# включим стартовый скрипт bash ~/.bashrc в логин скрипте bash ~/.bash_profile
echo 'if [ -f "$HOME/.bashrc" ]; then' >> ~/.bash_profile
echo '  . "$HOME/.bashrc"' >> ~/.bash_profile
echo 'fi' >> ~/.bash_profile

# установим multiwerf в директорию ~/bin
mkdir -p ~/bin
cd ~/bin
curl -L https://raw.githubusercontent.com/werf/multiwerf/master/get.sh | bash
```

##### Использование werf в текущей сессии shell

Следующий вызов создаст shell-функцию `werf`, которая вызывает бинарный файл той версии werf, которую multiwerf скачал и активировал:

```shell
source $(multiwerf use {{ include.version }} {{ include.channel }} --as-file)
werf version
...
```

##### Опционально: автоматически активировать werf при запуске терминала

```shell
echo '. $(multiwerf use {{ include.version }} {{ include.channel }} --as-file)' >> ~/.bashrc
```
