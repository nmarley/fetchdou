package main

import (
	"fmt"
	"time"
)

func requestHeaders() map[string]string {
	headers := make(map[string]string)

	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	headers["accept-encoding"] = "gzip, deflate"
	headers["accept-language"] = "pt-BR,pt;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6"
	headers["cookie"] = "GUEST_LANGUAGE_ID=pt_BR"
	headers["dnt"] = "1"
	headers["proxy-connection"] = "keep-alive"
	headers["upgrade-insecure-requests"] = "1"
	headers["user-agent"] = userAgent

	// TODO: address this later
	// headers["host"] = "download.in.gov.br"
	// headers["referer"] = "http://pesquisa.in.gov.br/imprensa/core/jornalList.action"

	return headers
}

func searchParams(date time.Time) map[string]string {
	strDDMM := fmt.Sprintf("%02d/%02d", date.Day(), date.Month())

	params := make(map[string]string)
	params["search-bar"] = ""
	params["tipo-pesquisa"] = "0"
	params["sistema-busca"] = "2"
	params["termo-pesquisado"] = "0"
	params["jornal"] = "do1"
	params["t"] = "com.liferay.journal.model.JournalArticle"
	params["g"] = "68942"
	params["edicao.jornal"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,600,601,602,603,612,613,614,615,701"
	params["checkbox_edicao.jornal"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,2,2000,529,525,3,3000,3020,1040,526,530,600,601,602,603,604,605,606,607,608,609,610,611,612,613,614,615,701,702"
	params["__checkbox_edicao.jornal"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,600,601,602,603,612,613,614,615,701"
	params["__checkbox_edicao.jornal"] = "2,2000,529,525,604,605,606,607,702"
	params["__checkbox_edicao.jornal"] = "3,3000,3020,1040,526,530,608,609,610,611"
	params["edicao.txtPesquisa"] = ""
	params["edicao.jornal_hidden"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,2,2000,529,525,3,3000,3020,1040,526,530,600,601,602,603,604,605,606,607,608,609,610,611,612,613,614,615,701,702"
	params["edicao.dtInicio"] = strDDMM
	params["edicao.dtFim"] = strDDMM
	params["edicao.ano"] = fmt.Sprintf("%04d", date.Year())

	return params
}
