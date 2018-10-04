# Formatting

import re
import os
import os.path
if 0: from typing import Pattern

# Pattern ie
# event_009[tag some-tag].jpg
# event2_012.jpg
FORMATPATTERN = r"{}_(\d+)(\[[\w\-_\s]+\])?\.(\w+)"

def getFormatReg(eventname=""): # type: (str) -> Pattern[str]
    """
    Format regular expression from folder name (event).
    If not provided... use the working directory as an event.
    """
    if not eventname: # Use working directory
        eventname = os.path.basename(os.getcwd())
    return re.compile(FORMATPATTERN.format(re.escape(eventname)))
