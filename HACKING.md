# HACKING logjam

This document helps programmers get started with software-development of **logjam**.

## Technology Stack

**logjam** is primarily written in the **Go programming-language** (**Golang**) on the back-end.

With the web-based front-end written in CSS, HTML, JavaScript, and TypeScript.

## Top-Level Files

* `main.go` — contains the `main()` procedure.
* `HACKING.md` — _this file_.
* `LICENSE` — information about ownership and licensing.
* `README.md` — more basic & more general useful information about this project.

## Sub-Directories

* `cfg/` — this is where the rest of the source-code gets configuration information from.
* `doc/` — files needs for docs (including some docs themselves).
* `env/` — handles configuration done via environment variables.
  * this (`env/`) should NOT be `import`ed by anything other than `flg/`; `import` from `cfg/` to get configuration information.
* `flg/` — handles configuration done via command-line flags/switch, while defauling to and potential configuration done via environment variables.
  * this (`flg/`) should NOT be `import`ed by anything other than `cfg/`; `import` from `cfg/` to get configuration information.
* `lib/` — libraries that are decoupled from the rest of the source-code base.
  * see also: [lib/README.md](lib/README.md)
* `srv/` — contains services that can be `import`ed and usedby the rest of te source-code.
* `web-app/` — the front-end code.
* `www/` — contains the HTTP handlers.
  * see also: [www/README.md](www/README.md)
