# fetchdou

> Download the Brazilian DOU - Diário Oficial da União

A web scraper written in Go which helps to download the PDF format copy of the [Diário Oficial da União][pt-wikipedia-url] - the official newspaper of the Brazilian government. Laws and governmental decrees are published here.

Note that the Brazilian government has weird rules about when this can be downloaded (only from 12:00 - 23:59 during the day in Brasília time zone, for example).

## Table of Contents
- [Install](#install)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Install

```sh
go get -u github.com/nmarley/fetchdou

go install github.com/nmarley/fetchdou/cmd/fetchdou
```

## Usage

Example to download DOU from date of 2019-11-05:

```sh
fetchdou 2019-11-05
```

## Contributing

Feel free to dive in! [Open an issue](https://github.com/nmarley/fetchdou/issues/new) or submit PRs.

## License

[ISC](LICENSE)

[pt-wikipedia-url]: https://pt.wikipedia.org/w/index.php?title=Di%C3%A1rio_Oficial_da_Uni%C3%A3o&oldid=56613715
