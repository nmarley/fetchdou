# fetchdou

> A package to help automated download the Brazilian government newspaper, Diário Oficial da União (DOU)

A web scraper written in Go which downloads the PDF format copy of the [Diário Oficial da União][pt-wikipedia-url] - the official newspaper of the Brazilian government. Laws and governmental decrees are published here.

The idea is to download a full PDF copy of the DOU every day when published and store it in S3 (public for all to access) with a reliable URL schema so that anyone can download this whenever without having to use the bullshit website which sucks ass or rely on the Brazilian government to store historical copies of it (LOL).

Ideally I will also include a full-text version which can scrape from the PDF version and enable some kinda of full-text copy (.txt) or full-text search so that anyone can search, again without having to use the bullshit Brazilian govt website which sucks ass.


Notes on sls app:

// 2018
// 2019
//   01
//   11
//     21
// 2020
//   01
//     01
// sha256sums of each PDF
// index.html of the whole thing
// /sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf

// Note: All this to be done in SLS.
//
// Code for parsing the .gov.br code to get PDF links could be in another Go
// package (which also has a command-line tool for downloading for a given
// day).
//
// Needs to have CloudFront simply for the caching if nothing else.
//
// Let's do a s3 structure of YYYY/MM/DD/FILENAME.pdf
//
// W/every file laid down, do a scan of the "directory" and create an
// index.html
//   (This will be a Lambda triggered by the s3 put)
//
// One option is to use DynamoDB for a metadata store. Can keep sha256sums of
// each PDF and also assist in the index.html for each "dir".



## Table of Contents
- [Install](#install)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Install

```sh
go get -u github.com/nmarley/go-baixar-dou

go install github.com/nmarley/go-baixar-dou/cmd/baixadou
```

## Usage

Example to download DOU from date of 2019-11-22:

```sh
baixadou 2019-11-22
```

## Contributing

Feel free to dive in! [Open an issue](https://github.com/nmarley/go-baixar-dou/issues/new) or submit PRs.

## License

[ISC](LICENSE)

[pt-wikipedia-url]: https://pt.wikipedia.org/w/index.php?title=Di%C3%A1rio_Oficial_da_Uni%C3%A3o&oldid=56613715
