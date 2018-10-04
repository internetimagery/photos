
from __future__ import print_function
from format import getFormatReg
import unittest

class TestFormat(unittest.TestCase):
    def test_format(s):
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

if __name__ == '__main__':
    unittest.main()
