import json
import os
import sys

sys.path.append('../../../libbeat/tests/system')

from beat.beat import TestCase


class BaseTest(TestCase):

    @classmethod
    def setUpClass(self):
        self.beat_name = "filebeat"
        super(BaseTest, self).setUpClass()

    def get_registry(self):
        # Returns content of the registry file
        dotYapdnsBeat = self.working_dir + '/registry'
        assert os.path.isfile(dotYapdnsBeat) is True

        with open(dotYapdnsBeat) as file:
            return json.load(file)
