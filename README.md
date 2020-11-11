# go_hasd_covid


Redis help.... 

open cli
redis-cli

get all keys
KEYS *

find all data for key 
smembers SCHOOL_HES_DATA



Data Extraction Process
- simple python program that fetches data from hortonville and processes it with the output being a csv file
- scheduled to extract morning noon and evening

Data Loading Service
- scheduled to run immediately after a run of the data extraction process
- loads data from data extraction process into a more refined and persistent data store