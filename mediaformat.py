# Formatting

import re
import os
import os.path
if 0: from typing import Pattern, Match, List

# Pattern ie
# event_009[tag some-tag].jpg
# event2_012.jpg
FORMATPATTERN = r"({})_(\d+)(\[[\w\-_\s]+\])?\.(\w+)"

def getFormatReg(eventname=""): # type: (str) -> Pattern[str]
    """
    Format regular expression from folder name (event).
    If not provided... use the working directory as an event.
    """
    if not eventname: # Use working directory
        eventname = os.path.basename(os.getcwd())
    return re.compile(FORMATPATTERN.format(re.escape(eventname)))

class MediaFormat(object):
    """
    Get formatting information for media.
    """
    __slots__ = ("name", "event", "index", "tags", "ext")
    def __init__(s): # type: () -> None
        s.name = s.event = s.ext = ""
        s.index = 0
        s.tags = [] # type: List[str]

    @classmethod
    def from_match(cls, match): # type: (Match[str]) -> MediaFormat
        media = MediaFormat()
        media.name = match.group(0)
        media.event = match.group(1)
        media.index = int(match.group(2))
        media.tags = (match.group(3) or "")[1:-1].split()
        media.ext = match.group(4)
        return media
