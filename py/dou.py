import sys
import os
import boto3
import base64
import mechanize
import http.cookiejar as cookielib
import datetime
import json
import re
# debugging
from pprint import pprint


class DOUGet(object):

    def __init__(self, **kwargs):
        self.agent = self._browser()
        self.pdf_dl_url = 'http://pesquisa.in.gov.br/imprensa/core/jornalList.action'

    def _browser(self):
        br = mechanize.Browser()
        cj = cookielib.LWPCookieJar()
        br.set_cookiejar(cj)

        br.set_handle_equiv(True)
        br.set_handle_gzip(True)
        br.set_handle_redirect(True)
        br.set_handle_referer(True)
        br.set_handle_robots(False)

        br.set_handle_refresh(mechanize._http.HTTPRefreshProcessor(),
                              max_time=1)
        br.addheaders = [('User-agent', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36')]

        return br

    def get_dou_url_for_day(self, day):
        year = "%04d" % day.year
        month = "%02d" % day.month
        day = "%02d" % day.day
        url = f'http://www.in.gov.br/leiturajornal?data={day}-{month}-{year}#daypicker'
        return url

    def get_search_params(self, day):
        # format dates for this idiotic asinine piece of shit app
        year = "%04d" % day.year
        month = "%02d" % day.month
        day = "%02d" % day.day
        paramDayMonth = f"{day}/{month}"
        paramYear = year

        params = {
            'search-bar': '',
            'tipo-pesquisa': '0',
            'sistema-busca': '2',
            'termo-pesquisado': '0',
            'jornal': 'do1',
            't': 'com.liferay.journal.model.JournalArticle',
            'g': '68942',
            'edicao.jornal': '1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,600,601,602,603,612,613,614,615,701',
            '__checkbox_edicao.jornal': '1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,2,2000,529,525,3,3000,3020,1040,526,530,600,601,602,603,604,605,606,607,608,609,610,611,612,613,614,615,701,702',
            '__checkbox_edicao.jornal': '1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,600,601,602,603,612,613,614,615,701',
            '__checkbox_edicao.jornal': '2,2000,529,525,604,605,606,607,702',
            '__checkbox_edicao.jornal': '3,3000,3020,1040,526,530,608,609,610,611',
            'edicao.txtPesquisa': '',
            'edicao.jornal_hidden': '1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,2,2000,529,525,3,3000,3020,1040,526,530,600,601,602,603,604,605,606,607,608,609,610,611,612,613,614,615,701,702',
            'edicao.dtInicio': paramDayMonth,
            'edicao.dtFim': paramDayMonth,
            'edicao.ano': paramYear,
        }

        return params

    def get_initial_page(self, day):
        url = self.get_dou_url_for_day(day)
        self.agent.open(url)
        self.agent.select_form('form_busca_dou')
        form = self.agent.form

        # get_search_params
        params = self.get_search_params(day)
        for key in params:
            val = params[key]
            ctrl = {}
            try:
                ctrl = form.find_control(key)
            except:
                # who the fuck cares, g-d Python sucks ass
                form.new_control(type='hidden', name=key, attrs='', ignore_unknown=True)
                ctrl = form.find_control(key)
            ctrl.readonly = False

            # check if "listcontrol"... my fucking g-d how much this shit sucks...
            # the final Go rewrite is gonna be heaven.
            #if ctrl is a ListControl
            try:
                form[key] = val
            except:
                form[key] = [val]

        form.action = 'http://pesquisa.in.gov.br/imprensa/core/jornalList.action'
        self.agent.submit()

        import pdb;pdb.set_trace()
        resp = self.agent.response()
        data = resp.read()
        links = self.extract_pdf_download_links(data)
        print(links)

        return

    def write_page(self, page):
        with open('page.html', 'wb') as f:
            f.write(page)

    def extract_pdf_download_links(self, data):
        import html
        links = []
        re_pdf_download_link = re.compile(r'(http://download.in.gov.br/[^\'"]*)')

        text = data.decode('utf-8')
        for m in re_pdf_download_link.finditer(text):
            # they actually HTML escape the links (what idiotic assholes)
            links.append(html.unescape(m.group(0)))

        return links


def lambda_handler(event, context):
    obj = DOUGet()
    # today = datetime.date.today()
    today = datetime.datetime(2019, 11, 21, 13, 0, 0, tzinfo=None)
    # today = datetime.datetime(2019, 1, 2, 13, 0, 0, tzinfo=None)

    obj.get_initial_page(today)

    return None


if __name__ == "__main__":
    lambda_handler(None, None)
