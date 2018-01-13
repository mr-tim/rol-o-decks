import base64
from flask import Flask, jsonify, request
from sqlalchemy import text
import subprocess

import database

app = Flask(__name__)

@app.route('/api/search')
def search():
    query = request.args.get('q', '')

    session = database.Session()

    query_results = session.query(database.SlideContent)\
        .join(database.Slide)\
        .join(database.Document)\
        .filter(text('''slide_content_content match '"''' + query + '''*"' '''))\
        .all()

    results = [search_result(x, query) for x in query_results]

    response = {'results': results}

    return jsonify(response)

@app.route('/api/open/<int:slideId>')
def open(slideId):
    session = database.Session()
    slide = session.query(database.Slide)\
        .get(slideId)

    cmd = ['open', slide.document.path]
    subprocess.run(cmd)
    return 'ok'

def search_result(x, search_term):
    t = x.content
    start = t.lower().find(search_term.lower())
    sub_start = t.rfind('\n', 0, start)+1
    sub_end = t.find('\n', start+len(search_term))
    if sub_end == -1:
        sub = t[sub_start:]
    else:
        sub = t[sub_start:sub_end]
    start = sub.lower().find(search_term.lower())
    match = {
        'text': sub,
        'start': start,
        'length': len(search_term)
    }
    return {
        'slideId': x.slide.id,
        'slide': x.slide.slide,
        'path': x.slide.document.path,
        'thumbnail': str(base64.b64encode(x.slide.thumnail_png))[2:-1],
        'match': match
    }
