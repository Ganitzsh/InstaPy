 # -*- coding: utf-8 -*-
from instapy import InstaPy

insta_username = ''
insta_password = ''

# set headless_browser=True if you want to run InstaPy on a server
# try:
    # set these if you're locating the library in the /usr/lib/pythonX.X/ directory
    # Settings.database_location = '/path/to/instapy.db'
    # Settings.browser_location = '/path/to/chromedriver'

session = InstaPy(username=insta_username,
                  password=insta_password,
                  headless_browser=False,
                  multi_logs=True)
try:
    session.login()

    # settings
    session.set_comments([
        'Cool shot! 📷',
        'Nice one!',
        'Really nice 🤩! Come visit our page sometimes 🙌',
        'Great content! 😃',
        'We like it!',
        'Nailed it! 🤙',
        'So cool! 🙌',
        'Neat! 📸',
        'Stunning 💫',
        '👍'
    ])
    session.set_upper_follower_count(limit=750)
    session.set_lower_follower_count(limit=400)
    session.set_do_comment(True, percentage=25)
    session.set_do_follow(enabled=True, percentage=5, times=2)
    # session.set_user_interact(amount=4, randomize=True, percentage=100, media='Photo')

    # actions
    session.like_by_tags([
        'toplisbonphoto',
        'tilestyle',
        'straightfacades',
        'visitlisboa',
        'wonderlustportugal',
        'archilovers',
        'cityshots',
        'lisbonbuildings',
        'colorfulhouses',
        'cityviews',
        'classicarchitecture',
        'visiteurope',
        'discovereurope',
        'archdesign',
        'buildingstyles',
        'streetview',
        'urbanscape',
        'citylandscape',
        'cityarchitecture'
    ],
    amount=30
    # interact=True
    )
    # session.interact_user_followers(['topluxembourgphoto'], amount=50, randomize=True)

except Exception as exc:
    # full stacktrace when raising Github issue
    traceback.print_exc(exc)

finally:
    # end the bot session
    session.end()
