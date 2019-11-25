import re
# debugging
from pprint import pprint


# examples:
# http://download.in.gov.br/sgpub/do/secao1/extra/2019/2019_11_21/2019_11_21_ASSINADO_do1_extra_A.pdf?arg1=v5uFuVtjRoHyXTCBrS-ILA&amp;arg2=1574739296
# http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&amp;arg2=1574739296
def extract_pdf_download_links(data):
    import html
    links = []
    re_pdf_download_link = re.compile(r'(http://download.in.gov.br/[^\'"]*)')

    text = data.decode('utf-8')
    for m in re_pdf_download_link.finditer(text):
        # they actually HTML escape the links (what idiotic assholes)
        links.append(html.unescape(m.group(0)))

    return links


def lambda_handler(event, context):
    html = ''
    with open('page2.html', 'rb') as f:
        html = f.read()
    links = extract_pdf_download_links(html)
    print("links: ", links)

    return None


if __name__ == "__main__":
    lambda_handler(None, None)
