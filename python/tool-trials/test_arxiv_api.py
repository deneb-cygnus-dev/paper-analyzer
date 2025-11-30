import urllib.request

url = "http://export.arxiv.org/api/query?search_query=cat:cs.SE&sortBy=lastUpdatedDate&sortOrder=descending&max_results=1"
data = urllib.request.urlopen(url)
print(data.read().decode("utf-8"))
