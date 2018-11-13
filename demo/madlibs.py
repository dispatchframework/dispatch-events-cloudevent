import datetime
import random
import time
import uuid

import requests

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
    callback = event["extensions"]["callback-url"]
    resp = requests.post(callback, headers={"Content-Type": "application/json"}, json=found_event)
    if not resp.ok:
        return {
            "error": resp.content
        }
    return found_event