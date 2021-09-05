from bs4 import BeautifulSoup
from datetime import datetime
import json
import requests
import pandas as pd
import sys
import os

os.chdir(os.path.dirname(sys.argv[0]))

# coding=utf-8
page = requests.get("https://www.hasd.org/community/covid-19-daily-updates.cfm")
html = page.content

soup = BeautifulSoup(html, 'html.parser')

list = list(soup.children)
start = soup.find('span', text='HMS-FWA')
tableToParse = start.parent.parent.parent.parent

# empty list
data = []

# for getting the header from
# the HTML file
list_header = []
soup = tableToParse
header = soup.find_all("tbody")[0].find("tr")

# for element in header:
#     try:
#         text = element.get_text()
#         if text == ' ':
#             text = ''
#         elif text == '-':
#             text = '0'
#
#         list_header.append(text)
#     except:
#         continue

# for getting the data
HTML_data = soup.find_all("tbody")[0].find_all("tr")

for element in HTML_data:
    sub_data = []
    for sub_element in element:
        if sub_element != '\n':

            try:
                text = sub_element.get_text()
                text = text.rstrip('\n')
                text = text.lstrip('\n')
                if text == ' ':
                    text = ''
                elif text == '-':
                    text = '0'
                elif text.isdigit():
                    text = text
                elif text.isalnum():
                    text = text
                elif len(text) > 3:
                    text = text
                else:
                    alphanumeric_filter = filter(str.isalnum, text)
                    alphanumeric_string = "".join(alphanumeric_filter)
                    text = alphanumeric_string

                sub_data.append(text)
            except:
                raise

    data.append(sub_data)
list_header = data.pop(0)

# Storing the data into Pandas
# DataFrame
dataFrame = pd.DataFrame(data = data, columns = list_header)



now = datetime.today().strftime('%Y%m%d%H%M%S')
# Converting Pandas DataFrame
# into CSV file
dataFrame.to_csv('incoming/covid_data_' + now + '.csv')