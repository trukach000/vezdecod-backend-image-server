import requests
import sys

host = 'https://trukach000.ru'
port = '443'
prefix = 'imloader'

if len(sys.argv) > 1:
    host= sys.argv[1] 
if len(sys.argv) > 2:
    port= sys.argv[2] 
if len(sys.argv) > 3:
    prefix = sys.argv[3]

files = {
    'image': open('./testdata/JPG_example.jpg', 'rb')
}

base_url = host + ':' + port
if prefix is not None and prefix != '':
    base_url = base_url + '/' + prefix

url = base_url + '/upload'    
print('uploading image to ' + url)
r = requests.post(url, verify=False, files=files)

if r.status_code != 200:
    print('wrong status code of upload response: ', r.status_code)
    exit(1)

obj = r.json()
print('Image token is '+ obj['imageToken'])
get_url = base_url + '/get/' + obj['imageToken']

print('getting image from ' + url)
r = requests.get(get_url, verify=False)
print(r)

