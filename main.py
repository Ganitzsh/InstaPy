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


session = InstaPy(
        username=insta_username,
        password=insta_password,
        headless_browser=headless,
        want_check_browser=False,
        multi_logs=True
)

try:
    with open(sys.argv[1], encoding='utf-8') as f:
        root = json.load(f)
        data = root['Settings']
        comments = data['Comments']
        hashtags = data['Hashtags']
        total_likes = data['TotalLikes']
        min_posts = data['MinPosts']
        max_follower = data['MaxFollowers']
        min_follower = data['MinFollowers']
        max_following = data['MaxFollowing']
        min_following = data['MinFollowing']
        per_user = data['PerUser']
    pass

    session.login()

    #  session.set_quota_supervisor(
    #          enabled=True,
    #          sleep_after=["likes", "comments_d", "follows", "unfollows", "server_calls_h"],
    #          sleepyhead=True,
    #          stochastic_flow=True,
    #          notify_me=True,
    #          peak_likes_hourly=57,
    #          peak_likes_daily=585,
    #          peak_comments_hourly=21,
    #          peak_comments_daily=182,
    #          peak_follows_hourly=48,
    #          peak_follows_daily=None,
    #          peak_unfollows_hourly=35,
    #          peak_unfollows_daily=402,
    #          peak_server_calls_hourly=None,
    #          peak_server_calls_daily=4700
    #  )

    session.set_skip_users(
            skip_private=True,
            skip_no_profile_pic=True,
            no_profile_pic_percentage=100
    )

    #  session.set_user_interact(
    #          amount=per_user,
    #          randomize=True,
    #          percentage=100,
    #          media='Photo'
    #  )

    session.set_relationship_bounds(
            enabled=True,
            delimit_by_numbers=True,
            min_posts=min_posts,
            max_followers=max_follower,
            min_followers=min_follower,
            max_following=max_following,
            min_following=min_following,
    )

    session.like_by_tags(hashtags, randomize=True)

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
