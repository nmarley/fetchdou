# go-baixar-dou

A web scraper written in Go which downloads the PDF format copy of the [Diário Oficial da União][pt-wikipedia-url] - the official newspaper of the Brazilian government. Laws and governmental decrees are published here.

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
