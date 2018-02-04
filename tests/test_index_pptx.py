import unittest

import index_pptx

class TestIndex(unittest.TestCase):
    def test_empty_text_elems(self):
        xml = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
        <p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" showMasterSp="1" showMasterPhAnim="1">
        <a:t/>
        </p:sld>
        """

        content = index_pptx.scrape_text_content(xml)
        self.assertEqual(content, '', 'Empty text elem should be allow on slides')
        