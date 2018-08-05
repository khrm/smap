# SMAP [![Go Report Card](https://goreportcard.com/badge/github.com/khrm/smap)](https://goreportcard.com/report/github.com/khrm/smap) [![Build Status](https://travis-ci.org/khrm/smap.svg?branch=master)](https://travis-ci.org/khrm/smap)

<!-- MarkdownTOC -->

- [What is it?](#what-is-it)
- [Prerequisites](#prerequisites)
- [Building](#building)
- [Running(Go)](#running-go)
- [Running(Docker)](#running-docker)

<!-- /MarkdownTOC -->

<a name="what-is-it"></a>
# What is it?

SMAP print the details of interconnections between various link in a site. Also, it prints all links found in a particular domain.
```Json
{
      "URLs": {
          "https://goharbor.io": {},
          "https://goharbor.io/blogs": {},
          "https://goharbor.io/blogs/harbor-joins-cncf": {},
          "https://goharbor.io/blogs/hello-world": {},
          "https://goharbor.io/community": {},
          "https://goharbor.io/docs": {}
      },
      "Connections": {
          "https://goharbor.io": {
              "https://goharbor.io": {},
              "https://goharbor.io/blogs": {},
              "https://goharbor.io/community": {},
              "https://goharbor.io/docs": {}
          },
          "https://goharbor.io/blogs": {
              "https://goharbor.io": {},
              "https://goharbor.io/blogs": {},
              "https://goharbor.io/blogs/harbor-joins-cncf": {},
              "https://goharbor.io/blogs/hello-world": {},
              "https://goharbor.io/community": {},
              "https://goharbor.io/docs": {}
          },
          "https://goharbor.io/blogs/harbor-joins-cncf": {
              "https://goharbor.io": {},
              "https://goharbor.io/blogs": {},
              "https://goharbor.io/community": {},
              "https://goharbor.io/docs": {}
          },
          "https://goharbor.io/blogs/hello-world": {
              "https://goharbor.io": {},
              "https://goharbor.io/blogs": {},
              "https://goharbor.io/community": {},
              "https://goharbor.io/docs": {}
          },
          "https://goharbor.io/community": {
              "https://goharbor.io": {},
              "https://goharbor.io/blogs": {},
              "https://goharbor.io/community": {},
              "https://goharbor.io/docs": {}
          },
          "https://goharbor.io/docs": {
              "https://goharbor.io": {},
              "https://goharbor.io/blogs": {},
              "https://goharbor.io/community": {},
              "https://goharbor.io/docs": {}
          }
      }
  }
StdSiteMap:
 <urlset xmlns="https://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://goharbor.io/docs</loc></url><url><loc>https://goharbor.io/blogs</loc></url><url><loc>https://goharbor.io/blogs/harbor-joins-cncf</loc></url><url><loc>https://goharbor.io/blogs/hello-world</loc></url><url><loc>https://goharbor.io</loc></url><url><loc>https://goharbor.io/community</loc></url></urlset>

```
<a name="prerequisites"></a>
## Prerequisites
Either Go or Docker should be installed.

<a name="building"></a>
## Building
You can use following commands to build if go is installed:

```shell
   $ make build
```

You will get a binary smap in the root folder.

Or you can use following commands if docker is installed.

```shell
   $ make image
```

Which will you an image smap:latest.

<a name="running-go"></a>
## Running(Go)
You can use :-

```shell
   $ ./smap -domain=example.com -depth=3
```

<a name="running-docker"></a>
## Running(Docker)
You can use :-

```shell
   $ docker run -ti smap smap -domain=kubeless.io

```



