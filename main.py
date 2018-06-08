# -*- coding: utf-8 -*-
import json
from random import sample
from instapy import InstaPy
from pprint import pprint

from selenium.common.exceptions import NoSuchElementException

try:
    with open('config.json') as f:
        data = json.load(f)
        insta_username = data['username']
        insta_password = data['password']
        clarifai_api_key = data['clarifai_api_key']
    pass
except Exception as e:
    raise

# set headless_browser=True if you want to run InstaPy on a server
# try:
    # set these if you're locating the library in the /usr/lib/pythonX.X/ directory
    # Settings.database_location = '/path/to/instapy.db'


session = InstaPy(username=insta_username,
                  password=insta_password,
                  headless_browser=True,
                  multi_logs=True)

try:
    with open('resources.json', encoding='utf-8') as f:
        data = json.load(f)
        comments = data['comments']
        hashtags = data['hashtags']
    pass
    session.login()

    # settings
    selection = sample(hashtags, 10) # Select 10 random hashtags
    session.set_comments(selection)
    session.set_do_comment(True, percentage=50)
    session.set_smart_hashtags(selection, limit=3, sort='top', log_tags=True)
    session.like_by_tags(amount=10, use_smart_hashtags=True)
    session.clarifai_check_img_for(['face'])
    session.set_use_clarifai(enabled=True, api_key=clarifai_api_key)
    session.like_by_feed(amount=100, randomize=True, unfollow=False, interact=True)

except Exception as exc:
    # if changes to IG layout, upload the file to help us locate the change
    if isinstance(exc, NoSuchElementException):
        file_path = os.path.join(gettempdir(), '{}.html'.format(time.strftime('%Y%m%d-%H%M%S')))
        with open(file_path, 'wb') as fp:
            fp.write(session.browser.page_source.encode('utf8'))
        print('{0}\nIf raising an issue, please also upload the file located at:\n{1}\n{0}'.format(
            '*' * 70, file_path))
    # full stacktrace when raising Github issue
    raise

finally:
    # end the bot session
    session.end()
