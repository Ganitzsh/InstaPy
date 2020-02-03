# -*- coding: utf-8 -*-
import json
import sys
from random import sample
from instapy import InstaPy
from pprint import pprint

from selenium.common.exceptions import NoSuchElementException

try:
    with open(sys.argv[1]) as f:
        data = json.load(f)
        print(data);
        insta_username = data['Account']['Username']
        insta_password = data['Account']['Password']
        #  headless = data['headless']
        headless = True
    pass
except Exception as e:
    raise

# set headless_browser=True if you want to run InstaPy on a server
# try:
    # set these if you're locating the library in the /usr/lib/pythonX.X/ directory
    # Settings.database_location = '/path/to/instapy.db'


session = InstaPy(username=insta_username,
                  password=insta_password,
                  headless_browser=headless,
                  multi_logs=True)

try:
    with open(sys.argv[1], encoding='utf-8') as f:
        root = json.load(f)
        data = root['Settings']
        comments = data['Comments']
        hashtags = data['Hashtags']
        total_likes = data['TotalLikes']
        max_follower = data['MaxFollowers']
        min_follower = data['MinFollowers']
        max_following = data['MaxFollowing']
        min_following = data['MinFollowing']
        potency = data['Potency']
        per_user = data['PerUser']
    pass
    session.login()

    # settings
    # selection = sample(hashtags, s) # Select sample random hashtags
    # session.set_comments(comments)
    # session.set_do_comment(True, percentage=50)
    # session.set_smart_hashtags(selection, limit=3, sort='top', log_tags=True)
    ratio = -1.0
    # if potency == 'positive':
    #     ratio = max_follower / max_following
    # elif potency == 'negative':
    #     ratio = -(max_following / max_follower)

    session.set_user_interact(amount=per_user, randomize=True, percentage=50, media='Photo')
    session.set_relationship_bounds(
            enabled=True,
            potency_ratio=ratio,
            delimit_by_numbers=True,
            max_followers=max_follower,
            min_followers=min_follower,
            max_following=max_following,
            min_following=min_following,
    )
    session.like_by_tags(hashtags, use_smart_hashtags=False, randomize=True)
    session.like_by_feed(amount=total_likes, randomize=True, unfollow=False, interact=True)

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
