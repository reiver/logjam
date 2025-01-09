# LOGJAM

**logjam** is a P2P video & audio technology for decentralized social-media (DeSo).

<img src="doc/img/greatape-logo.png" style="width:150px" />

The documentation for **logjam** is divided up into several files.

* [RUNNING.md](doc/RUNNING.md) — Any wishing to run **logjam**, please see [RUNNING.md](doc/RUNNING.md) file.
* [HACKING.md](doc/HACKING.md) — Any wishing to contribute to the **logjam** source-code, please see [HACKING.md](doc/HACKING.md) file.

## PocketBase Base URL

To run **logjam** you need to tell it where your **PocketBase** server is, using a **PocketBase base-URL**.

You can set this using 2 different methods.

### PocketBase Base URL — environment-variable

Method №1 is to use the `POCKETBASE_URL` **environment-variable**.

For example, if you wanted to change the **PocketBase base-URL** TCP port to `http://example.com/api` then you could do something similar to the following:

```
POCKETBASE_URL="http://example.com/api" ./logjam
```

### PocketBase Base URL — command-line switch / flag

Method №2 is to use the `--pburl` **command-line** switch/flag.

For example, if you wanted to change the **PocketBase base-URL** TCP port to `http://example.com/api` then you could do something similar to the following:

```
./logjam --pburl='http://example.com/api'
```

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

## HTTP TCP Port Configuration (optional)

**logjam** includes a web-server.

By default this web-server runs on TCP port `8080`.

This can be changed to a different TCP port using either of 2 different methods.

### HTTP TCP Port Configuration (optional) — environment-variable

Method №1 is to use the `PORT` **environment-variable**.

For example, if you wanted to change the web-server's TCP port to `9000` then you could do something similar to the following:

```
PORT=9000 ./logjam
```

### HTTP TCP Port Configuration (optional) — command-line switch / flag

Method №2 is to use the `--src` **command-line** switch/flag.

For example, if you wanted to change the web-server's TCP port to `9000` then you could do something similar to the following:

```
./logjam --src=9000
```

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
