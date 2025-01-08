# LOGJAM

**logjam** is a P2P video & audio technology for decentralized social-media (DeSo).

<img src="doc/img/greatape-logo.png" style="width:150px" />

## Running The Project

---

first you need to get the binary by one of this ways :

1. downloading it from **releases** page
2. **building** it your self
3. leaving the build and run to **docker**!

- ### getting the release
1. download the latest binary from here https://github.com/reiver/logjam/releases
2. make it executable like `chmod +x ./logjam`
3. run it! `./logjam`

**note:** default listen host(`0.0.0.0`) or listen port(`8090`) can be changed using the `--listen-host` and `--listen-port` arguments.

- ### building

first make sure you have golang installed *(>=v1.20)*
1. clone the repository and `cd` into it.
2. build it using go like `go build .`
3. run it! `./logjam`

- ### using Docker

make sure you have docker installed first.
1. clone the repository and `cd` into it.
2. if you have docker-compose installed then run `docker-compose up`
3. if you choose to build and run it with docker itself then :
   1. build the image: `docker build -t logjam:latest .`
   2. create and run the container: `docker run --rm logjam:latest`

## HTTP TCP Port

**logjam** includes a web-server.

By default this web-server runs on TCP port `8080`.

This can be changed to a different TCP port using either of 2 different methods.

### HTTP TCP Port — environment-variable

Method №1 is to use the `PORT` **environment-variable**.

For example, if you wanted to change the web-server's TCP port to `9000` then you could do something similar to the following:

```
PORT=9000 ./logjam
```

### HTTP TCP Port — command-line switch / flag

Method №2 is to use the `--src` **command-line** switch/flag.

For example, if you wanted to change the web-server's TCP port to `9000` then you could do something similar to the following:

```
./logjam --src=9000
```

## Hacking

Any wishing to contribute to the **logjam** source-code, please see [HACKING.md](HACKING.md) file.

## Credits

Many people worked on **logjam**.

Some continue to work on it.
Some are alumni.

In alphabetical order —

### Credits — Back-End Software Engineers

* Benyamin Azarkhazin
* [Charles Iliya Krempeaux](http://reiver.link/)
* Mehrdad Mirsamie

### Credits — Front-End Software Engineers

* Massoud Seifi
* Muhammad Zaid Ali
* Sal Rahman

### Credits — Logo

* Chet Earl Woodside

### Credits — Product Researchers

* [Charles Iliya Krempeaux](http://reiver.link/)

### Credits — QA Specialists

* Atanaz Ostovar
* Sepideh Farsi

### Credits — UX/UI Specialists

* Farzaneh Amini

## Cookie

```
  ░▒▓█▓▒░      ░▒▓██████▓▒░ ░▒▓██████▓▒░       ░▒▓█▓▒░░▒▓██████▓▒░░▒▓██████████████▓▒░
  ░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒▒▓███▓▒░      ░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓████████▓▒░▒▓██████▓▒░ ░▒▓██████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░
```
