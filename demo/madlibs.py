import datetime
import random
import time
import uuid

import requests

# event:
#     eventtype: word.found.noun
#     eventtypeversion: ""
#     cloudeventsversion: "0.1"
#     source: http://srcdog.com/madlibs
#     eventid: 96fb5f0b-001e-0108-6dfe-da6e2806f124
#     eventtime: "0001-01-01T00:00:00.000Z"
#     schemaurl: ""
#     contenttype: ""
#     extensions:
#       callback-url: https://srcdog.com/madlibs/event
#     data:
#     - 110
#     - 117
#     - 108
#     - 108
#     executedtime: 1542068285
def handle(ctx, payload):
    event = ctx["event"]

    resp = requests.get("https://srcdog.com/madlibs/words.txt")
    words = resp.json()

    # no error checking... bad developer
    print(event)
    _, _, word_type = event["eventType"].split(".")

    picked = random.choice(words[word_type])

    found_event = {
        "specversion": event["cloudEventsVersion"],
        "type": "word.picked.%s" % word_type,
        "source": "http://demo.dispatchframework.io/dispatch/madlibs",
        "id": str(uuid.uuid4()),
        "time": datetime.datetime.utcnow().isoformat("T") + "Z",
        "relatedid": event["eventID"],
        "contentType": "application/json",
        "data": {
            "word": picked
        }
    }
    print("sending event to %s" % event["extensions"]["callback-url"])
    callback = "https://srcdog.com/madlibs/event"
    # callback = event["extensions"]["callback-url"]
    resp = requests.post(callback, headers={"Content-Type": "application/json"}, json=found_event)
    if not resp.ok:
        return {
            "error": resp.content
        }
    return found_event