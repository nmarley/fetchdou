import sys
import os
import boto3
import base64
import mechanize
import http.cookiejar as cookielib
import datetime
import json


class DOUGet(object):

    def __init__(self, **kwargs):
        self.agent = self._browser()
        self.init_url_re = 'http://www.in.gov.br/leiturajornal?data=DD-MM-YYYY#daypicker'
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

    def get_initial_page(self, day):
        year = day.year
        month = day.month
        day = day.day

        date = f"{year}-{month}-{day}"
        print(f"Hi, date: [{date}]")

        #self.agent.open(self.login_url)
        #self.agent.select_form('LoginForm')
        #form = self.agent.form
#
#        form['username'] = self.username
#        form['password'] = self.password
#        self.agent.submit()

        return

    def log_time_entry(self, date):
        req = self.agent.click_link(text='Time Entry')
        page = self.agent.open(req)

        self.agent.select_form('ProjectEntryForm')
        form = self.agent.form

        form.find_control('StartDate').readonly = False
        # -OR- form.set_all_readonly(False)
        # allow changing the .value of all controls

        form['hours'] = '8'
        form['StartDate'] = date.strftime('%m-%d-%Y')

        self.agent.submit()

        return


def lambda_handler(event, context):
    timesheet = DOUGet()
    # today = datetime.date.today()
    today = datetime.datetime(2019, 11, 22, 13, 0, 0, tzinfo=None)

    timesheet.get_initial_page(today)
    # today = datetime.datetime(2017, 6, 16, 13, 0, 0, tzinfo=None)
    #if today.weekday() <= 4:
    #    timesheet.log_time_entry(today)

    return None


if __name__ == "__main__":
    lambda_handler(None, None)
