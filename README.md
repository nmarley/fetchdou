# go-baixar-dou

Download the DOU - Diário Oficial da União

A web scraper written in Go which downloads the PDF format copy of the [Diário Oficial da União][pt-wikipedia-url] - the official newspaper of the Brazilian government. Laws and governmental decrees are published here.

The idea is to download a full PDF copy of the DOU every day when published and store it in S3 (public for all to access) with a reliable URL schema so that anyone can download this whenever without having to use the bullshit website which sucks ass or rely on the Brazilian government to store historical copies of it (LOL).

Ideally I will also include a full-text version which can scrape from the PDF version and enable some kinda of full-text copy (.txt) or full-text search so that anyone can search, again without having to use the bullshit Brazilian govt website which sucks ass.

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
