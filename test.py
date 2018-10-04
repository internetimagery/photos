
from __future__ import print_function
from mediaformat import getFormatReg, MediaFormat
import unittest

class TestFormat(unittest.TestCase):
    def test_format(s): # type: () -> None
        names = [
            ("34-234-434 event_0045[hello_to you].jpg", "34-234-434 event", True),
            ("3234 some event_0434.gfg", "3234 some event", True),
            ("_434234.jpg", "", False),
            ("3434gfgd_her4[dfs sf].jpg", "3434gfgd", False)]
        for test in names:
            matcher = getFormatReg(test[1])
            match = matcher.match(test[0])
            if test[2]:
                s.assertIsNotNone(match)
            else:
                s.assertIsNone(match)

    def test_mediaformat(s): # type: () -> None
        event = "18-02-23 some event"
        names = [
            ("{}_034[one two three].jpg".format(event), {"event": event, "index": 34, "tags": ["one", "two", "three"]}),
            ("{}_001.jpg".format(event), {"event": event, "index": 1, "tags": []})]
        matcher = getFormatReg(event)
        for test in names:
            match = matcher.match(test[0])
            s.assertIsNotNone(match)
            assert match is not None # make mypy happy
            media = MediaFormat.from_match(match)
            for attr, value in test[1].items():
                s.assertEqual(getattr(media, attr), value)



if __name__ == '__main__':
    unittest.main()
